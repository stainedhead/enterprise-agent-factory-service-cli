package runtime

import (
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
