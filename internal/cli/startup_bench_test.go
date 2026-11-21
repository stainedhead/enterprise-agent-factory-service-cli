package cli_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli/services"
)

func startup(t testing.TB, url string) time.Duration {
	start := time.Now()
	reg, err := services.NewRegistry()
	if err != nil {
		t.Fatal(err)
	}
	var out, eb bytes.Buffer
	env := cli.NewEnv(strings.NewReader(""), &out, &eb)
	env.Vars = map[string]string{"AFS_ENDPOINT_URL": url, "AFS_API_KEY": "k"}
	env.Home = t.TempDir()
	if code := cli.Run(context.Background(), []string{"sso", "get-caller-identity"}, env, reg); code != 0 {
		t.Fatalf("exit %d: %s", code, eb.String())
	}
	return time.Since(start)
}

func fakeMe(first *time.Time) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if first.IsZero() {
			*first = time.Now()
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"principal":{"id":"u","type":"user"},"serviceAdmin":false,"securityOperator":false,"teams":[]}`))
	}))
}

// NFR startup: under 100 ms from process start of the command table to the first request.
func TestStartupUnder100ms(t *testing.T) {
	var first time.Time
	srv := fakeMe(&first)
	defer srv.Close()
	var ds []time.Duration
	for i := 0; i < 7; i++ {
		first = time.Time{}
		begin := time.Now()
		startup(t, srv.URL)
		ds = append(ds, first.Sub(begin))
	}
	sort.Slice(ds, func(i, j int) bool { return ds[i] < ds[j] })
	median := ds[len(ds)/2]
	t.Logf("startup to first request: median %v (min %v, max %v)", median, ds[0], ds[len(ds)-1])
	if median > 100*time.Millisecond {
		t.Errorf("median startup %v exceeds 100ms", median)
	}
}

func BenchmarkStartup(b *testing.B) {
	var first time.Time
	srv := fakeMe(&first)
	defer srv.Close()
	for i := 0; i < b.N; i++ {
		startup(b, srv.URL)
	}
}
