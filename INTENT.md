# Intent

## Purpose
Give people and automation a command line to the Enterprise Agent Factory Service that anyone who uses
the AWS CLI already knows how to drive: `afs <service> <operation> --options`, profiles, SSO login,
`--output`, `--query` and pagination that behave the same way.

## Where it fits
One of four children of [enterprise-agent-factory-service](https://github.com/stainedhead/enterprise-agent-factory-service).
It is a client of the REST API in `enterprise-agent-factory-service-restapi` and follows its OpenAPI
contract. The web console in `enterprise-agent-factory-service-website` offers the same operations in a
browser; every console action has an `afs` equivalent.

## Goals
- Every API operation available as one CLI operation, named predictably from the API.
- AWS CLI conventions for options, configuration, output, errors and exit codes.
- Respect the role model: the API's decision is final, and the CLI explains a denial clearly.
- Safe with credentials and secrets: nothing sensitive in output, logs or shell history.
- Usable by agents and pipelines as well as people (stable JSON output, API keys, non-interactive mode).

## Non-goals
- Calling Kubernetes or AWS directly. Everything goes through the API.
- Client-side workflows that bundle several changes; those belong in the API or in scripts.
- Enforcing access control; that is the API's job.
