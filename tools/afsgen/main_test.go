package main

import (
	"os"
	"strings"
	"testing"
)

func TestGenerateMini(t *testing.T) {
	raw, err := os.ReadFile("testdata/mini.yaml")
	if err != nil {
		t.Fatal(err)
	}
	a, err := Generate(raw)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := Generate(raw)
	if string(a) != string(b) {
		t.Fatal("generation is not deterministic")
	}
	s := strings.Join(strings.Fields(string(a)), " ")
	for _, want := range []string{
		`ContractVersion = "9.8.7"`,
		`ID: "CreateThing"`, `ID: "DeleteThing"`, `ID: "ListThings"`, `ID: "PutThing"`,
		`Service: "things", Command: "list-things"`,
		`Paginated: true`, `ResultKey: "rows"`, `CSV: true`,
		`Produces: []string{"application/json", "text/csv"}`,
		`IfMatch: IfMatchOptional`, `IfMatch: IfMatchRequired`, `Idempotent: true`,
		`Success: []int{201, 202}`, `Success: []int{204}`,
		`Ref: "Money"`, `Nullable: true`, `Items: &Schema{Type: "string"}`,
		`Type: "any"`, `Max: f64(5)`, `MaxLen: i64(9)`, `Enum: []string{"a", "b"}`,
		`Service: "", Command: ""`,
	} {
		if !strings.Contains(s, want) {
			t.Errorf("output lacks %s", want)
		}
	}
	if strings.Index(s, `ID: "CreateThing"`) > strings.Index(s, `ID: "ListThings"`) {
		t.Error("not sorted by operationId")
	}
	if strings.Contains(s, "maxResults") || strings.Contains(s, "If-Match") {
		t.Error("excluded parameters leaked")
	}
}

func TestGenerateRejectsBadCommandName(t *testing.T) {
	_, err := Generate([]byte("openapi: 3.1.0\ninfo: {version: 1.0.0}\npaths:\n  /x:\n    get:\n      operationId: X\n      parameters: []\n      responses: {}\n      x-afs-cli: BadName\n"))
	if err == nil {
		t.Fatal("want error")
	}
	_, err = Generate([]byte("openapi: 3.1.0\ninfo: {version: 1.0.0}\npaths:\n  /x:\n    get:\n      operationId: X\n      parameters: []\n      responses: {}\n"))
	if err == nil {
		t.Fatal("want error for missing x-afs-cli")
	}
}
