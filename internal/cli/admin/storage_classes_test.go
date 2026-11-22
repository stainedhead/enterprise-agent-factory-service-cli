package admin_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDoneStorageClasses(t *testing.T) {
	clitest.CheckOperations(t, "storage-classes", "ListStorageClassesAdmin", "PutStorageClass", "DeleteStorageClass")
}

// CLI-1: deleting a class in use is a 409 (254).
func TestDeleteStorageClassInUse409(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			clitest.Reply(w, 200, `{"items":[{"name":"gp3-standard","version":2}]}`)
			return
		}
		clitest.Reply(w, 409, `{"code":"Conflict","message":"in use by 3 volumes","requestId":"r"}`)
	}, "admin", "delete-storage-class", "--name", "gp3-standard")
	if r.Code != 254 || !strings.Contains(r.Err, "(Conflict)") {
		t.Fatalf("%+v", r)
	}
}
