package admin_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/clitest"
)

func TestCommandsDoneCatalog(t *testing.T) {
	clitest.CheckOperations(t, "catalog", "ListCatalogSystems", "GetCatalogSystem", "CreateCatalogSystem", "UpdateCatalogSystem", "DeleteCatalogSystem")
}

// CLI-1: deleting a built-in system is a 409 (254).
func TestDeleteBuiltinCatalogSystem409(t *testing.T) {
	r := clitest.Run(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			clitest.Reply(w, 200, `{"id":"ext_1","version":1,"builtin":true}`)
			return
		}
		clitest.Reply(w, 409, `{"code":"Conflict","message":"built-in systems cannot be deleted","requestId":"r"}`)
	}, "admin", "delete-external-system", "--system-id", "ext_1")
	if r.Code != 254 || !strings.Contains(r.Err, "(Conflict) when calling the DeleteCatalogSystem operation: built-in systems cannot be deleted") {
		t.Fatalf("%+v", r)
	}
}
