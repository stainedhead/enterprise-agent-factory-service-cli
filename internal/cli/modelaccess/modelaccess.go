// Package modelaccess holds what agents and pools share about model access (API contract section 4.4 *Models*,
// MOD-3, MOD-8): --models takes the word all or one or more model ids, --default-model-id names the default
// model, and --model-id, the 1.0.1 spelling, stays for now with a deprecation line.
package modelaccess

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
)

// All is the word of --models that means every model the team may use, now and later (OR-D13).
const All = "all"

// ModelsDesc is the help text of --models.
const ModelsDesc = "Models the agent may call: the word all (every model the team may use, now and later) or one or more model ids (a fixed selection). Sent members replace the stored ones."

// DefaultModelDesc is the help text of --default-model-id.
const DefaultModelDesc = "The model the runtime starts with; it must be one of the agent's models. Without --models on an update, only the default changes."

// DeprecatedModelID is the help text and warning of --model-id.
const DeprecatedModelID = "use --default-model-id (and --models to choose the set) instead."

// RoutesDesc is the help text of --routes (v1.5.0).
const RoutesDesc = "Model routes the agent may call: the word all (every route of the team, now and later) or one or more route slugs (a fixed selection). Folded into models.routeAccess and models.routes."

// DefaultRouteDesc is the help text of --default-route (v1.5.0).
const DefaultRouteDesc = "The route of the default model; --default-model-id then names the route's model. Or give --default-model-id <route>/<model> alone."

// Options are the local options the model override adds.
var Options = []cli.Option{
	{Flag: "--default-model-id", In: "local", Kind: cli.String, Desc: DefaultModelDesc},
	{Flag: "--routes", In: "local", Kind: cli.List, Elem: "string", Desc: RoutesDesc},
	{Flag: "--default-route", In: "local", Kind: cli.String, Desc: DefaultRouteDesc},
}

// Override returns the override that gives an operation with a models member (create and update of agents
// and pools) the options of MOD-8. Merge it with the operation's own override.
func Override(op string) cli.Override {
	return cli.Override{
		Op:                op,
		Extra:             Options,
		Kinds:             map[string]cli.Kind{"models": cli.List},
		Desc:              map[string]string{"models": ModelsDesc},
		DeprecatedOptions: map[string]string{"modelId": DeprecatedModelID},
		Build:             Build,
	}
}

// Build turns --models and --default-model-id into the models member {access, modelIds, defaultModelId}.
// A models member that is already an object (from --cli-input-json) is kept, and --default-model-id is
// folded into it. Every failure is a usage error (252) before any call.
func Build(p *cli.Parsed, parts *cli.RequestParts) error {
	if parts.Body == nil {
		return nil
	}
	if p.Has("default-model-id") {
		if _, ok := parts.Body.Get("modelId"); ok {
			return errors.New("--model-id is the deprecated name of --default-model-id; give one of them")
		}
	}
	obj := cli.NewObject()
	if raw, ok := parts.Body.Get("models"); ok {
		var words []string
		if err := json.Unmarshal(raw, &words); err == nil {
			conv, err := fromWords(words)
			if err != nil {
				return err
			}
			obj = conv
		} else {
			parsed, err := cli.ParseObject(raw)
			if err != nil {
				return errors.New("--models takes the word all or model ids")
			}
			obj = parsed
		}
	}
	if p.Has("default-model-id") {
		b, _ := json.Marshal(p.Value("default-model-id"))
		obj.Set("defaultModelId", b)
	}
	if p.Has("routes") {
		if err := routesFromWords(p.Values["routes"], obj); err != nil {
			return err
		}
	}
	if p.Has("default-route") {
		b, _ := json.Marshal(p.Value("default-route"))
		obj.Set("defaultRoute", b)
	}
	if obj.Len() > 0 {
		b, _ := json.Marshal(obj)
		parts.Body.Set("models", b)
	}
	return nil
}

// fromWords reads the words of --models: all, or a list of model ids.
func fromWords(words []string) (*cli.Object, error) {
	obj := cli.NewObject()
	if len(words) == 0 {
		return nil, errors.New("--models takes the word all or at least one model id")
	}
	for _, w := range words {
		if w == All {
			if len(words) > 1 {
				return nil, fmt.Errorf("--models %s cannot be combined with model ids", All)
			}
			obj.Set("access", json.RawMessage(`"all"`))
			return obj, nil
		}
		if w == "" {
			return nil, errors.New("--models: a model id cannot be empty")
		}
	}
	obj.Set("access", json.RawMessage(`"selected"`))
	ids, _ := json.Marshal(words)
	obj.Set("modelIds", ids)
	return obj, nil
}

// routesFromWords reads the words of --routes into obj: all, or a list of route slugs (v1.5.0).
func routesFromWords(words []string, obj *cli.Object) error {
	if len(words) == 0 {
		return errors.New("--routes takes the word all or at least one route slug")
	}
	for _, w := range words {
		if w == All {
			if len(words) > 1 {
				return fmt.Errorf("--routes %s cannot be combined with route slugs", All)
			}
			obj.Set("routeAccess", json.RawMessage(`"all"`))
			return nil
		}
		if w == "" {
			return errors.New("--routes: a route slug cannot be empty")
		}
	}
	obj.Set("routeAccess", json.RawMessage(`"selected"`))
	b, _ := json.Marshal(words)
	obj.Set("routes", b)
	return nil
}
