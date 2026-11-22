// Package gen holds request and response types generated from the pinned OpenAPI snapshot
// (oapi-codegen, see the Makefile). They are not used to decode responses for output (D-C5);
// packages that read response members by hand reference those members here so a contract rename
// fails the build (FR-017, ADR-014).
package gen
