// Package wait implements --wait: poll a resource until it is Ready at the mutation's generation,
// fails, or (for deletes) answers 404 (CLI-12, D-C9).
package wait

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/spec"
)

// Target says how one mutation is waited on.
type Target struct {
	PollOp string // operation polled; its body carries status.phase
	GetOp  string // operation printed at the end (and polled for 404 on delete)
	Delete bool   // wait for GetOp to answer 404
}

// Table maps mutations that accept --wait to their targets (data-dictionary section 4).
var Table = map[string]Target{}

func init() {
	add := func(poll, get string, ops ...string) {
		for _, o := range ops {
			Table[o] = Target{PollOp: poll, GetOp: get, Delete: strings.HasPrefix(o, "Delete")}
		}
	}
	add("GetTeamStatus", "GetTeam", "CreateTeam", "UpdateTeam", "SuspendTeam", "ResumeTeam", "PutTeamModels", "DeleteTeam")
	add("GetAgentStatus", "GetAgent", "CreateAgent", "UpdateAgent", "PauseAgent", "ResumeAgent", "RestartAgent", "DeleteAgent")
	add("GetPoolStatus", "GetPool", "CreatePool", "UpdatePool", "PausePool", "ResumePool", "ScalePool", "DeletePool")
	add("GetLogDestination", "GetLogDestination", "CreateLogDestination", "UpdateLogDestination", "DeleteLogDestination")
	add("GetTeamCollaboration", "GetTeamCollaboration", "PutTeamCollaboration")
	add("GetNotificationChannel", "GetNotificationChannel", "CreateNotificationChannel", "UpdateNotificationChannel", "DeleteNotificationChannel")
}

// Doer sends one request.
type Doer interface {
	Do(ctx context.Context, req client.Request) (*client.Response, error)
}

// Options control the polling.
type Options struct {
	Timeout  time.Duration // default 10m
	Interval time.Duration // default 5s
	Now      func() time.Time
	Sleep    func(ctx context.Context, d time.Duration) error
	// ReadTimeout lowers the cap on a poll response's Retry-After (20 s) when it is smaller; 0 = none.
	ReadTimeout time.Duration
}

// ErrUndecodable is returned (exit 255) when a status response is not the JSON object the contract
// describes, so a contract change cannot make --wait run silently until its timeout.
var ErrUndecodable = errors.New("the status response could not be decoded")

// ErrTimeout is returned (exit 255) when the resource is not ready in time.
var ErrTimeout = errors.New("timed out waiting for the resource")

// FailedError means the resource reached Failed or Degraded (exit 254).
type FailedError struct {
	Phase      string
	Conditions json.RawMessage
}

func (e *FailedError) Error() string {
	s := "the resource is " + e.Phase
	if len(e.Conditions) > 0 && string(e.Conditions) != "null" {
		s += "; conditions: " + string(e.Conditions)
	}
	return s
}

type statusDoc struct {
	Generation json.Number `json:"generation"`
	Status     *struct {
		Phase              string          `json:"phase"`
		ObservedGeneration json.Number     `json:"observedGeneration"`
		Conditions         json.RawMessage `json:"conditions"`
	} `json:"status"`
	Phase              string          `json:"phase"`
	ObservedGeneration json.Number     `json:"observedGeneration"`
	Conditions         json.RawMessage `json:"conditions"`
	// Workers is present in a pool's status: the pool is done when every desired worker reports (CLI-12).
	Workers *struct {
		Desired   json.Number `json:"desired"`
		Reporting json.Number `json:"reporting"`
	} `json:"workers"`
}

func num(n json.Number) (int64, bool) {
	if n == "" {
		return 0, false
	}
	v, err := strconv.ParseInt(string(n), 10, 64)
	return v, err == nil
}

// Evaluate reads a status body against the generation of the mutation (0 when unknown).
func Evaluate(body []byte, generation int64) (done bool, err error) {
	var d statusDoc
	dec := json.NewDecoder(strings.NewReader(string(body)))
	dec.UseNumber()
	if dec.Decode(&d) != nil {
		return false, ErrUndecodable
	}
	phase, obs, cond := d.Phase, d.ObservedGeneration, d.Conditions
	if d.Status != nil {
		phase, obs, cond = d.Status.Phase, d.Status.ObservedGeneration, d.Status.Conditions
	}
	if o, ok := num(obs); ok && generation > 0 && o < generation {
		return false, nil // stale: the controller has not seen this generation yet
	}
	if phase == "Failed" || phase == "Degraded" {
		return false, &FailedError{Phase: phase, Conditions: cond}
	}
	if d.Workers != nil {
		want, ok1 := num(d.Workers.Desired)
		have, ok2 := num(d.Workers.Reporting)
		if ok1 && ok2 {
			return have == want, nil
		}
	}
	return phase == "Ready", nil
}

// Generation reads the generation of a mutation response (0 when absent).
func Generation(body []byte) int64 {
	var d struct {
		Generation json.Number `json:"generation"`
	}
	dec := json.NewDecoder(strings.NewReader(string(body)))
	dec.UseNumber()
	if dec.Decode(&d) != nil {
		return 0
	}
	g, _ := num(d.Generation)
	return g
}

// Until polls until done. mutation is the mutation's response body; pollPath and getPath are
// the paths of the poll and final read. It returns the final resource body (nil after a delete).
func Until(ctx context.Context, d Doer, t Target, mutation []byte, pollPath, getPath string, o Options) ([]byte, error) {
	if o.Timeout <= 0 {
		o.Timeout = 10 * time.Minute
	}
	if o.Interval <= 0 {
		o.Interval = 5 * time.Second
	}
	if o.Now == nil {
		o.Now = time.Now
	}
	sleep := o.Sleep
	if sleep == nil {
		sleep = func(ctx context.Context, d time.Duration) error {
			tm := time.NewTimer(d)
			defer tm.Stop()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-tm.C:
				return nil
			}
		}
	}
	gen := Generation(mutation)
	deadline := o.Now().Add(o.Timeout)
	var interval time.Duration // set by every poll: its Retry-After or o.Interval
	for {
		if t.Delete {
			resp, err := d.Do(ctx, client.Request{Method: "GET", Path: getPath, Operation: t.GetOp})
			var ae *client.APIError
			if errors.As(err, &ae) && ae.Status == 404 {
				return nil, nil
			}
			if err != nil {
				return nil, err
			}
			interval = pollInterval(resp, o)
		} else {
			resp, err := d.Do(ctx, client.Request{Method: "GET", Path: pollPath, Operation: t.PollOp})
			if err != nil {
				return nil, err
			}
			done, ferr := Evaluate(resp.Body, gen)
			if errors.Is(ferr, ErrUndecodable) {
				return nil, fmt.Errorf("%s: %w", t.PollOp, ferr)
			}
			if ferr != nil {
				return nil, ferr
			}
			interval = pollInterval(resp, o)
			if done {
				if getPath == pollPath {
					return resp.Body, nil
				}
				final, err := d.Do(ctx, client.Request{Method: "GET", Path: getPath, Operation: t.GetOp})
				if err != nil {
					return nil, err
				}
				return final.Body, nil
			}
		}
		if !o.Now().Add(interval).Before(deadline) {
			return nil, fmt.Errorf("%w (%s)", ErrTimeout, o.Timeout)
		}
		if err := sleep(ctx, interval); err != nil {
			return nil, err
		}
	}
}

// pollInterval is the wait before the next poll: the response's Retry-After (capped) or the interval (D-C9).
func pollInterval(resp *client.Response, o Options) time.Duration {
	if ra := client.RetryAfter(resp.Header, o.Now(), o.ReadTimeout); ra > 0 {
		return ra
	}
	return o.Interval
}

// Paths returns the poll and get paths of a wait target for a mutation: path parameters come
// from the request, then from the "id" member of the mutation response (creates).
func Paths(t Target, reqPath map[string]string, mutation []byte) (poll, get string) {
	var m map[string]json.RawMessage
	_ = json.Unmarshal(mutation, &m)
	build := func(opID string) string {
		op := spec.ByID(opID)
		if op == nil {
			return ""
		}
		pp := map[string]string{}
		for _, p := range op.Params {
			if p.In != "path" {
				continue
			}
			if v := reqPath[p.Name]; v != "" {
				pp[p.Name] = v
				continue
			}
			var s string
			_ = json.Unmarshal(m["id"], &s)
			pp[p.Name] = s
		}
		return client.BuildPath(op.Path, pp)
	}
	return build(t.PollOp), build(t.GetOp)
}
