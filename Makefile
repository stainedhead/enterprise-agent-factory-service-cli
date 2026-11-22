.PHONY: all build test race vet lint fmt fmt-check tidy vuln cover gate clean generate generate-check contract-check contract-sync docs

VERSION ?= dev
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
PKG     := github.com/stainedhead/enterprise-agent-factory-service-cli/internal/buildinfo
LDFLAGS := -s -w -X $(PKG).Version=$(VERSION) -X $(PKG).Commit=$(COMMIT) -X $(PKG).Date=$(DATE)
BINS    := afs

GOLANGCI_LINT_VERSION ?= v2.15.0
GOVULNCHECK_VERSION   ?= v1.8.0

all: gate

build:
	@mkdir -p bin
	@for b in $(BINS); do CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/$$b ./cmd/$$b || exit 1; done

test:
	go test ./...

race:
	go test -race ./...

vet:
	go vet ./...

# Build golangci-lint with this module's toolchain (go.mod), so it can analyze Go 1.27 code.
lint:
	GOTOOLCHAIN=$(shell go env GOVERSION) go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run

fmt:
	gofmt -w .

fmt-check:
	@out="$$(gofmt -l .)"; if [ -n "$$out" ]; then echo "$$out"; exit 1; fi

tidy:
	go mod tidy

vuln:
	go run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) ./...

# Coverage gate: internal/... (generated files and the golden test helper excluded) must stay at or above 90%.
COVER_MIN ?= 90
cover:
	@go test -coverpkg=./internal/... -coverprofile=cover.out ./internal/... >/dev/null
	@grep -v -e '_gen.go' -e '/gen/' -e '/internal/golden/' -e '/internal/clitest/' cover.out > cover.filtered.out
	@total=$$(go tool cover -func=cover.filtered.out | awk '/^total:/ {gsub("%","",$$3); print $$3}'); \
	  echo "internal/... coverage: $$total% (minimum $(COVER_MIN)%)"; \
	  awk -v t="$$total" -v m="$(COVER_MIN)" 'BEGIN { exit !(t+0 >= m+0) }' || { echo "coverage below $(COVER_MIN)%"; exit 1; }
	@rm -f cover.out cover.filtered.out

# The snapshot must match api/SOURCE, and the generated spec must name the same checksum.
contract-check:
	@want=$$(sed -n 's/^sha256: //p' api/SOURCE); got=$$(sha256sum api/openapi.yaml | cut -d' ' -f1); \
	  [ "$$want" = "$$got" ] || { echo "api/openapi.yaml sha256 $$got != api/SOURCE $$want"; exit 1; }
	@grep -q "SnapshotSHA256 = \"$$(sed -n 's/^sha256: //p' api/SOURCE)\"" internal/spec/operations_gen.go || { echo "internal/spec is not generated from the pinned snapshot"; exit 1; }

# Pin the snapshot: make contract-sync RESTAPI_DIR=../enterprise-agent-factory-service-restapi
# or RESTAPI_REF=<branch|tag|commit> (fetched with gh).
RESTAPI_REPO ?= stainedhead/enterprise-agent-factory-service-restapi
contract-sync:
	@if [ -n "$(RESTAPI_DIR)" ]; then \
	  cp "$(RESTAPI_DIR)/api/openapi.yaml" api/openapi.yaml; \
	  ref=$$(git -C "$(RESTAPI_DIR)" rev-parse --abbrev-ref HEAD); commit=$$(git -C "$(RESTAPI_DIR)" rev-parse HEAD); \
	elif [ -n "$(RESTAPI_REF)" ]; then \
	  gh api -H "Accept: application/vnd.github.raw" "repos/$(RESTAPI_REPO)/contents/api/openapi.yaml?ref=$(RESTAPI_REF)" > api/openapi.yaml || exit 1; \
	  ref=$(RESTAPI_REF); commit=$$(gh api "repos/$(RESTAPI_REPO)/commits/$(RESTAPI_REF)" --jq .sha); \
	else echo "set RESTAPI_DIR or RESTAPI_REF"; exit 1; fi; \
	printf 'repo: %s\nref: %s\ncommit: %s\nsha256: %s\nsynced_at: %s\n' "$(RESTAPI_REPO)" "$$ref" "$$commit" "$$(sha256sum api/openapi.yaml | cut -d' ' -f1)" "$$(date -u +%Y-%m-%dT%H:%M:%SZ)" > api/SOURCE
	$(MAKE) generate

OAPI_CODEGEN_VERSION ?= v2.8.2
# oapi-codegen needs a newer toolchain than the base one; use the module's.
GOTOOLCHAIN_PIN ?= $(shell go env GOVERSION)

generate:
	go run ./tools/afsgen -in api/openapi.yaml -out internal/spec/operations_gen.go
	GOTOOLCHAIN=$(GOTOOLCHAIN_PIN) go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION) -config internal/client/gen/oapi-codegen.yaml api/openapi.yaml
	gofmt -w internal/spec/operations_gen.go internal/client/gen/types.gen.go

# Regenerate into a temp dir and compare with the committed files.
generate-check:
	@d=$$(mktemp -d); \
	go run ./tools/afsgen -in api/openapi.yaml -out $$d/operations_gen.go && \
	sed "s#^output:.*#output: $$d/types.gen.go#" internal/client/gen/oapi-codegen.yaml > $$d/cfg.yaml && \
	GOTOOLCHAIN=$(GOTOOLCHAIN_PIN) go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION) -config $$d/cfg.yaml api/openapi.yaml && \
	gofmt -w $$d/operations_gen.go $$d/types.gen.go && \
	diff -q $$d/operations_gen.go internal/spec/operations_gen.go && diff -q $$d/types.gen.go internal/client/gen/types.gen.go; \
	rc=$$?; rm -f $$d/*.go $$d/cfg.yaml; rmdir $$d; exit $$rc

docs:
	go run ./tools/docgen

# The full quality gate. Must pass before every commit.
gate: tidy contract-check generate-check fmt-check vet lint race cover build
	./bin/afs --version

clean:
	$(RM) -r ./bin ./dist
