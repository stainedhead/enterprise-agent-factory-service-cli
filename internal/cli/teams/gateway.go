package teams

import (
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/cli"
	"github.com/stainedhead/enterprise-agent-factory-service-cli/internal/output"
)

func init() { registrations = append(registrations, registerGateway) }

// gatewayColumns are the table columns of a team's gateway settings (screen 60, v1.5.0).
var gatewayColumns = []output.Column{
	{Header: "Team", Path: "teamId"}, {Header: "Deployment", Path: "deployment"},
	{Header: "MayAddEndpoints", Path: "mayAddModelEndpoints"}, {Header: "State", Path: "gatewayState"},
	{Header: "Url", Path: "gatewayUrl", Null: "-"}, {Header: "Version", Path: "version"},
}

// selectionColumns are the rows of a team's opted-in service endpoints.
var selectionColumns = []output.Column{
	{Header: "EndpointId", Path: "endpointId"}, {Header: "Slug", Path: "slug"}, {Header: "SelectedBy", Path: "selectedBy"},
}

func registerGateway(r *cli.Registry) {
	r.Override(cli.Override{Op: "PutTeamGatewaySettings", Columns: gatewayColumns,
		Notes: []string{"shared sends the team's agents through the service's shared gateway; per_namespace runs the team's own gateway in its namespace (agents roll to it once it is ready). mayAddModelEndpoints lets the team's admins declare their own endpoints."},
		Examples: []cli.Example{
			{Title: "Give a team its own gateway", Command: "afs teams put-team-gateway-settings --team-id " + exTeam + " --deployment per_namespace"},
			{Title: "Let a team declare its own model endpoints", Command: "afs teams put-team-gateway-settings --team-id " + exTeam + " --may-add-model-endpoints"},
		}})
	r.Override(cli.Override{Op: "PutTeamModelEndpoints", Columns: selectionColumns, ResultKey: "selections",
		Notes: []string{"Replaces the team's opted-in service endpoints: an endpoint left out is deselected. Endpoints that need extra terms, or whose provider is selected by service admins, need a service admin (needs_service_admin). No If-Match is sent unless --if-match is given."},
		Examples: []cli.Example{
			{Title: "Opt the team in to two service routes", Command: "afs teams put-team-model-endpoints --team-id " + exTeam + " --endpoint-ids mep_01JD2Q3R4S5T6V7W8X9Y0Z0001 mep_01JD2Q3R4S5T6V7W8X9Y0Z0003"},
		}})
}
