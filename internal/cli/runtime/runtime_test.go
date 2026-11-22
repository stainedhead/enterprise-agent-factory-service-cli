package runtime

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/client"
)

// FR-006: the identity-provider client shares the API client's transport and read timeout.
func TestIdPClientUsesTransportAndReadTimeout(t *testing.T) {
	env := cli.NewEnv(strings.NewReader(""), &strings.Builder{}, &strings.Builder{})
	hc := IdPClient(env, cli.Globals{}).(*http.Client)
	c := hc
	if c.Timeout != client.DefaultReadTimeout || c.Transport == nil || c.Transport == http.DefaultTransport {
		t.Errorf("default: %+v", c)
	}
	rt := http.DefaultTransport
	env.RoundTripper = rt
	c = IdPClient(env, cli.Globals{ReadTimeout: 7, Set: map[string]bool{"cli-read-timeout": true}}).(*http.Client)
	if c.Timeout != 7*time.Second || c.Transport != rt {
		t.Errorf("custom: %+v", c)
	}
}

// FR-R2: the credential helper used by invoke-api resolves the chain without needing an endpoint.
func TestCredentials(t *testing.T) {
	env := cli.NewEnv(strings.NewReader(""), &strings.Builder{}, &strings.Builder{})
	env.Home = t.TempDir()
	env.Vars = map[string]string{"AFS_API_KEY": "afsk_test"}
	st, fn, err := Credentials(env, cli.Globals{})
	if err != nil || st.Profile != "default" || st.EndpointURL != "" {
		t.Fatalf("%+v %v", st, err)
	}
	c, err := fn(context.Background())
	if err != nil || c.Header != "Bearer afsk_test" || c.Source != "env AFS_API_KEY" {
		t.Errorf("%+v %v", c, err)
	}
	env.Vars = map[string]string{}
	if c, err = fn(context.Background()); err != nil || c.Header != "" {
		t.Errorf("no credential must give an empty one: %+v %v", c, err)
	}
	if _, _, err = Credentials(env, cli.Globals{Profile: "nope", Set: map[string]bool{"profile": true}}); err == nil {
		t.Error("an unknown profile must fail")
	}
}

func TestTransportUsesInjected(t *testing.T) {
	env := cli.NewEnv(strings.NewReader(""), &strings.Builder{}, &strings.Builder{})
	if Transport(env, cli.Globals{}) == nil {
		t.Fatal("default transport")
	}
	env.RoundTripper = http.DefaultTransport
	if Transport(env, cli.Globals{}) != http.DefaultTransport {
		t.Error("injected transport ignored")
	}
}
