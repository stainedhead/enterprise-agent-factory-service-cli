// Package modelendpoints holds the overrides of the model-endpoints service: the catalog of OpenAI-compatible model
// endpoints (routes) a team gateway forwards to, their models, prices, credentials and federation, the provider
// switches, the provider statements, and a team's own endpoints (API contract section 4.13, v1.5.0, MEG-CLI-1 to
// MEG-CLI-4). Each command is one API call; credential values are read only from a file, standard input or a
// no-echo prompt.
package modelendpoints
