package validation

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	v2client "github.com/honeycombio/terraform-provider-honeycombio/client/v2"
)

// configKeyOnlyPermissions are the permission attribute names that are only valid
// on configuration API keys.
var configKeyOnlyPermissions = []string{
	"manage_boards",
	"manage_columns",
	"manage_markers",
	"manage_recipients",
	"manage_slos",
	"manage_triggers",
	"read_service_maps",
	"run_queries",
	"send_events",
	"visible_team_members",
}

// APIKeyConfigKeyOnlyPermissions returns a validator.List that prevents
// configuration-key-only permissions from being enabled on ingest API keys.
func APIKeyConfigKeyOnlyPermissions() validator.List {
	return configKeyOnlyPermissionsValidator{}
}

var _ validator.List = configKeyOnlyPermissionsValidator{}

type configKeyOnlyPermissionsValidator struct{}

func (v configKeyOnlyPermissionsValidator) Description(_ context.Context) string {
	return "Configuration-key-only permissions may not be set on ingest API keys."
}

func (v configKeyOnlyPermissionsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v configKeyOnlyPermissionsValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var keyType types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("type"), &keyType)...)
	if resp.Diagnostics.HasError() || keyType.IsUnknown() || keyType.ValueString() != v2client.APIKeyTypeIngest {
		return
	}

	elements := req.ConfigValue.Elements()
	if len(elements) == 0 {
		return
	}

	elem, ok := elements[0].(types.Object)
	if !ok {
		return
	}

	attrs := elem.Attributes()
	for _, name := range configKeyOnlyPermissions {
		val, exists := attrs[name]
		if !exists {
			continue
		}
		boolVal, ok := val.(types.Bool)
		if !ok || boolVal.IsNull() || boolVal.IsUnknown() || !boolVal.ValueBool() {
			continue
		}
		resp.Diagnostics.AddAttributeError(
			req.Path.AtListIndex(0).AtName(name),
			"Invalid Permission for API Key Type",
			fmt.Sprintf("The %q permission is only supported on configuration API keys.", name),
		)
	}
}
