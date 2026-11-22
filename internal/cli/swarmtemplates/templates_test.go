package swarmtemplates_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

const tmpl = "swt_01J9ZQ3ABCDEFGHJKMNPQRSTVW"

func TestCommandsDoneSwarmTemplates(t *testing.T) {
	clitest.CheckOperations(t, "swarmtemplates", "ListSwarmTemplates", "CreateSwarmTemplate", "GetSwarmTemplate",
		"UpdateSwarmTemplate", "DeleteSwarmTemplate", "EnableSwarmTemplate", "DisableSwarmTemplate")
}

// CLI-17, CLI-US-22.1: disable sends DisableSwarmTemplate; enable reverses it.
func TestDisableEnableSwarmTemplate(t *testing.T) {
	for _, tc := range []struct {
		cmd, path string
		enabled   bool
	}{{"disable-swarm-template", "/v1/admin/swarm-templates/" + tmpl + ":disable", false}, {"enable-swarm-template", "/v1/admin/swarm-templates/" + tmpl + ":enable", true}} {
		t.Run(tc.cmd, func(t *testing.T) {
			body := `{"id":"` + tmpl + `","name":"cve","kind":"cve_remediation","enabled":false,"budgetUsd":"250.00","ttlHours":72,"teardownRule":"at_expiry","version":2}`
			if tc.enabled {
				body = strings.Replace(body, `"enabled":false`, `"enabled":true`, 1)
			}
			r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, body) },
				"swarm-templates", tc.cmd, "--template-id", tmpl)
			if r.Code != 0 || r.Req.Method != "POST" || r.Req.URL.Path != tc.path || !strings.Contains(r.Out, `"enabled"`) {
				t.Fatalf("%+v", r)
			}
		})
	}
}

// CLI-US-22.1: a security operator is refused (254).
func TestSwarmTemplateSecurityOperatorDenied(t *testing.T) {
	r := clitest.Run(t, clitest.Denied("admin:swarm-templates"), "swarm-templates", "disable-swarm-template", "--template-id", tmpl)
	if r.Code != 254 || !strings.Contains(r.Err, "(AccessDenied) when calling the DisableSwarmTemplate operation") || !strings.Contains(r.Err, "admin:swarm-templates") {
		t.Fatalf("%+v", r)
	}
}

// CLI-17: deleting a template that an open work order uses is a Conflict naming the work orders.
func TestDeleteSwarmTemplateOpenWorkOrder409(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 409, `{"code":"Conflict","message":"template is in use","requestId":"req_c","details":{"dependents":["wo_1"]}}`)
	}, "swarm-templates", "delete-swarm-template", "--template-id", tmpl)
	if r.Code != 254 || !strings.Contains(r.Err, "(Conflict)") || !strings.Contains(r.Err, "Dependents: wo_1") {
		t.Fatalf("%+v", r)
	}
}

// CLI-17: create takes the spec from a file.
func TestCreateSwarmTemplateSpecFile(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 201, `{"id":"`+tmpl+`","name":"cve","kind":"general","enabled":true,"budgetUsd":"10.00","ttlHours":1,"teardownRule":"at_expiry","version":1}`)
	},
		"swarm-templates", "create-swarm-template", "--name", "cve", "--kind", "general", "--budget-usd", "10.00", "--ttl-hours", "1",
		"--teardown-rule", "at_expiry", "--spec", `{"workers":2}`)
	if r.Code != 0 || !strings.Contains(string(r.ReqBody), `"spec":{"workers":2}`) || !strings.Contains(string(r.ReqBody), `"budgetUsd":"10.00"`) {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
}

// Final contract: a template's spec secrets entries are typed {pool, name, description?, delivery?, target?}
// and reach the API unchanged; the CLI never holds a value.
func TestCreateSwarmTemplateTypedSecrets(t *testing.T) {
	spec := `{"pools":[{"name":"fixers"}],"secrets":[{"pool":"fixers","name":"registry-token","description":"Registry read token","delivery":"file","target":"/etc/registry/token"},{"pool":"fixers","name":"scan-key"}]}`
	r := clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 201, `{"id":"`+tmpl+`","name":"cve","kind":"general","enabled":true,"budgetUsd":"10.00","ttlHours":1,"teardownRule":"at_expiry","version":1}`)
	},
		"swarm-templates", "create-swarm-template", "--name", "cve", "--kind", "general", "--budget-usd", "10.00", "--ttl-hours", "1",
		"--teardown-rule", "at_expiry", "--spec", spec)
	if r.Code != 0 || !strings.Contains(string(r.ReqBody), `"spec":`+spec) {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
}

const collabSpec = `{"pools":[{"name":"fixers"}],"collaboration":{"templateId":"tpl_01JD2Q3R4S5T6V7W8X9Y0Z0001"}}`
const tmplDoc = `{"id":"` + tmpl + `","name":"cve","kind":"general","enabled":true,"budgetUsd":"10.00","ttlHours":1,"teardownRule":"at_expiry","spec":` + collabSpec + `,"version":1}`

func createTemplate(t *testing.T, h http.HandlerFunc) clitest.Result {
	t.Helper()
	return clitest.Run(t, h, "swarm-templates", "create-swarm-template", "--name", "cve", "--kind", "general", "--budget-usd", "10.00", "--ttl-hours", "1",
		"--teardown-rule", "at_expiry", "--spec", collabSpec)
}

// CLI-US-35.1
func TestCreateSwarmTemplateSpecCollaborationUnchanged(t *testing.T) {
	r := createTemplate(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 201, tmplDoc) })
	if r.Code != 0 || !strings.Contains(string(r.ReqBody), `"spec":`+collabSpec) {
		t.Fatalf("%+v body=%s", r, r.ReqBody)
	}
	r = clitest.Run(t, func(w http.ResponseWriter, _ *http.Request) { clitest.Reply(w, 200, tmplDoc) },
		"swarm-templates", "describe-swarm-template", "--template-id", tmpl, "--output", "json")
	if r.Code != 0 || !strings.Contains(r.Out, `"collaboration"`) || !strings.Contains(r.Out, `"templateId": "tpl_01JD2Q3R4S5T6V7W8X9Y0Z0001"`) {
		t.Errorf("%+v", r)
	}
}

func TestSwarmTemplateSkeletonHasCollaboration(t *testing.T) {
	for _, op := range []string{"create-swarm-template", "update-swarm-template"} {
		r := clitest.Run(t, nil, "swarm-templates", op, "--generate-cli-skeleton")
		if r.Code != 0 || !strings.Contains(r.Out, `"collaboration"`) || !strings.Contains(r.Out, `"templateId"`) || len(r.Requests) != 0 {
			t.Errorf("%s: %+v", op, r)
		}
	}
}

// CLI-US-35.2
func TestSwarmTemplateCollaborationTemplateRefused254(t *testing.T) {
	r := createTemplate(t, func(w http.ResponseWriter, _ *http.Request) {
		clitest.Reply(w, 400, `{"code":"ValidationError","message":"invalid spec","requestId":"req_s","details":{"fields":[{"field":"spec.collaboration.templateId","reason":"collaboration_template_without_capability"}]}}`)
	})
	if r.Code != 254 || !strings.Contains(r.Err, "(ValidationError) when calling the CreateSwarmTemplate operation") || !strings.Contains(r.Err, "field: spec.collaboration.templateId reason: collaboration_template_without_capability") {
		t.Errorf("%+v", r)
	}
}

func TestSwarmTemplateCollaborationAccessDenied254(t *testing.T) {
	r := createTemplate(t, clitest.Denied("swarm_templates:write"))
	if r.Code != 254 || !strings.Contains(r.Err, "(AccessDenied) when calling the CreateSwarmTemplate operation") {
		t.Errorf("%+v", r)
	}
}
