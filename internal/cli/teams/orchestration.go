package teams

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
)

func init() { registrations = append(registrations, registerOrchestration) }

// orchestrationDesc is the help text of create-team --orchestration (ORC-26, ORC-29).
const orchestrationDesc = "The team's orchestration agent, all optional: Tool=paperclip|buzz,WebUiExposure=off|private|public,WebApiExposure=off|private|public,RegistrationMode=manual|managed,ModelId=<model id>, or JSON (file://PATH). Defaults come from the orchestration.* settings."

// shorthandKeys maps the keys of the --orchestration shorthand, lower case, to what they set.
var shorthandKeys = map[string]func(o *cli.Object, v string) error{
	"tool":             setEnum("tool", "tool", []string{"paperclip", "buzz"}),
	"webuiexposure":    setSurface("webUi", "WebUiExposure"),
	"webapiexposure":   setSurface("webApi", "WebApiExposure"),
	"registrationmode": setEnum("registrationMode", "RegistrationMode", []string{"manual", "managed"}),
	"modelid":          setModel,
	"defaultmodelid":   setModel,
}

var levels = []string{"off", "private", "public"}

func setEnum(member, key string, allowed []string) func(*cli.Object, string) error {
	return func(o *cli.Object, v string) error {
		for _, a := range allowed {
			if a == v {
				b, _ := json.Marshal(v)
				o.Set(member, b)
				return nil
			}
		}
		return fmt.Errorf("--orchestration: %s must be one of %s, got %q", capitalize(key), strings.Join(allowed, ", "), v)
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func setSurface(member, key string) func(*cli.Object, string) error {
	return func(o *cli.Object, v string) error {
		inner := cli.NewObject()
		if err := setEnum("exposure", key, levels)(inner, v); err != nil {
			return err
		}
		b, _ := json.Marshal(inner)
		o.Set(member, b)
		return nil
	}
}

func setModel(o *cli.Object, v string) error {
	if v == "" {
		return errors.New("--orchestration: ModelId cannot be empty")
	}
	m := cli.NewObject()
	b, _ := json.Marshal(v)
	m.Set("defaultModelId", b)
	mb, _ := json.Marshal(m)
	o.Set("models", mb)
	return nil
}

// buildOrchestration turns the shorthand Key=Value,Key=Value of --orchestration into the object the API takes.
// JSON and file://PATH reach the body unchanged (the generic reader gives a JSON text its own value), so only a
// value that arrives as a JSON string is read as shorthand (ModelId is shorthand for models.defaultModelId, CD-36).
func buildOrchestration(_ *cli.Parsed, parts *cli.RequestParts) error {
	if parts.Body == nil {
		return nil
	}
	raw, ok := parts.Body.Get("orchestration")
	if !ok {
		return nil
	}
	var text string
	if json.Unmarshal(raw, &text) != nil {
		return nil
	}
	obj := cli.NewObject()
	seen := map[string]bool{}
	for _, pair := range strings.Split(text, ",") {
		k, v, found := strings.Cut(pair, "=")
		k = strings.ToLower(strings.TrimSpace(k))
		if !found || k == "" {
			return fmt.Errorf("--orchestration takes Key=Value pairs or JSON (keys: %s), got %q", strings.Join(shorthandNames(), ", "), pair)
		}
		set, known := shorthandKeys[k]
		if !known {
			return fmt.Errorf("--orchestration: unknown key %q (keys: %s)", strings.TrimSpace(pair[:strings.Index(pair, "=")]), strings.Join(shorthandNames(), ", "))
		}
		if seen[k] {
			return fmt.Errorf("--orchestration: %q is given twice", k)
		}
		seen[k] = true
		if err := set(obj, strings.TrimSpace(v)); err != nil {
			return err
		}
	}
	b, _ := json.Marshal(obj)
	parts.Body.Set("orchestration", b)
	return nil
}

func shorthandNames() []string {
	return []string{"Tool", "WebUiExposure", "WebApiExposure", "RegistrationMode", "ModelId"}
}

func registerOrchestration(r *cli.Registry) {
	r.Override(cli.Override{Op: "CreateTeam",
		Kinds:             map[string]cli.Kind{"orchestration": cli.Any},
		Desc:              map[string]string{"orchestration": orchestrationDesc},
		Build:             buildOrchestration,
		DeprecatedOptions: map[string]string{"mayExposeCollaboration": "use --max-collaboration-exposure (off, private or public) instead."},
		Examples: []cli.Example{
			{Title: "Create a team whose orchestration agent runs Buzz", Command: `afs teams create-team --slug payments --display-name "Payments" --kind persistent --max-budget-usd 500.00 --first-admin '{"email":"ada@example.com"}' --orchestration Tool=buzz --wait`},
			{Title: "Create a team that may expose its Web UI privately, and serve it that way", Command: `afs teams create-team --slug billing --display-name "Billing" --kind persistent --max-budget-usd 300.00 --first-admin '{"email":"ada@example.com"}' --max-collaboration-exposure private --orchestration Tool=paperclip,WebUiExposure=private,WebApiExposure=off,ModelId=anthropic.claude-sonnet-4-5`},
		}})
	r.Override(cli.Override{Op: "UpdateTeam",
		DeprecatedOptions: map[string]string{"mayExposeCollaboration": "use --max-collaboration-exposure (off, private or public) instead."},
		Examples: []cli.Example{
			{Title: "Cap how far the team's orchestration agent may be exposed (service admins)", Command: "afs teams update-team --team-id team_01J9ZQ3ABCDEFGHJKMNPQRSTVW --max-collaboration-exposure private"},
		}})
}
