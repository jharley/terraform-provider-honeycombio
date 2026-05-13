package validation_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	"github.com/honeycombio/terraform-provider-honeycombio/internal/helper/validation"
)

// permissionsAttrTypes defines the attribute types of the permissions block used in tests.
var permissionsAttrTypes = map[string]attr.Type{
	"create_datasets":      types.BoolType,
	"manage_boards":        types.BoolType,
	"manage_columns":       types.BoolType,
	"manage_markers":       types.BoolType,
	"manage_recipients":    types.BoolType,
	"manage_slos":          types.BoolType,
	"manage_triggers":      types.BoolType,
	"read_service_maps":    types.BoolType,
	"run_queries":          types.BoolType,
	"send_events":          types.BoolType,
	"visible_team_members": types.BoolType,
}

// apiKeyTestSchema is a minimal schema matching the permissions + type attributes of the API key resource.
var apiKeyTestSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"type": schema.StringAttribute{Required: true},
	},
	Blocks: map[string]schema.Block{
		"permissions": schema.ListNestedBlock{
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"create_datasets":      schema.BoolAttribute{Optional: true},
					"manage_boards":        schema.BoolAttribute{Optional: true},
					"manage_columns":       schema.BoolAttribute{Optional: true},
					"manage_markers":       schema.BoolAttribute{Optional: true},
					"manage_recipients":    schema.BoolAttribute{Optional: true},
					"manage_slos":          schema.BoolAttribute{Optional: true},
					"manage_triggers":      schema.BoolAttribute{Optional: true},
					"read_service_maps":    schema.BoolAttribute{Optional: true},
					"run_queries":          schema.BoolAttribute{Optional: true},
					"send_events":          schema.BoolAttribute{Optional: true},
					"visible_team_members": schema.BoolAttribute{Optional: true},
				},
			},
		},
	},
}

// makePermissions constructs a types.List for the permissions block from the given attribute map.
func makePermissions(t *testing.T, attrs map[string]attr.Value) types.List {
	t.Helper()
	defaults := map[string]attr.Value{
		"create_datasets":      types.BoolValue(false),
		"manage_boards":        types.BoolValue(false),
		"manage_columns":       types.BoolValue(false),
		"manage_markers":       types.BoolValue(false),
		"manage_recipients":    types.BoolValue(false),
		"manage_slos":          types.BoolValue(false),
		"manage_triggers":      types.BoolValue(false),
		"read_service_maps":    types.BoolValue(false),
		"run_queries":          types.BoolValue(false),
		"send_events":          types.BoolValue(false),
		"visible_team_members": types.BoolValue(false),
	}
	for k, v := range attrs {
		defaults[k] = v
	}
	obj := types.ObjectValueMust(permissionsAttrTypes, defaults)
	return types.ListValueMust(types.ObjectType{AttrTypes: permissionsAttrTypes}, []attr.Value{obj})
}

// makeAPIKeyConfig builds a tfsdk.Config with the given key type and permissions list value.
func makeAPIKeyConfig(ctx context.Context, t *testing.T, keyType string, perms types.List) tfsdk.Config {
	t.Helper()

	permsTFValue, err := perms.ToTerraformValue(ctx)
	if err != nil {
		t.Fatalf("ToTerraformValue for permissions: %s", err)
	}

	return tfsdk.Config{
		Schema: apiKeyTestSchema,
		Raw: tftypes.NewValue(
			apiKeyTestSchema.Type().TerraformType(ctx),
			map[string]tftypes.Value{
				"type":        tftypes.NewValue(tftypes.String, keyType),
				"permissions": permsTFValue,
			},
		),
	}
}

func Test_APIKeyConfigKeyOnlyPermissions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	emptyPerms := types.ListValueMust(types.ObjectType{AttrTypes: permissionsAttrTypes}, []attr.Value{})

	tests := []struct {
		name        string
		configValue types.List
		config      tfsdk.Config
		expectError bool
	}{
		{
			name:        "null list — skipped",
			configValue: types.ListNull(types.ObjectType{AttrTypes: permissionsAttrTypes}),
			config:      tfsdk.Config{},
			expectError: false,
		},
		{
			name:        "unknown list — skipped",
			configValue: types.ListUnknown(types.ObjectType{AttrTypes: permissionsAttrTypes}),
			config:      tfsdk.Config{},
			expectError: false,
		},
		{
			name:        "configuration key — config-only permissions allowed",
			configValue: makePermissions(t, map[string]attr.Value{"manage_boards": types.BoolValue(true)}),
			config:      makeAPIKeyConfig(ctx, t, "configuration", makePermissions(t, map[string]attr.Value{"manage_boards": types.BoolValue(true)})),
			expectError: false,
		},
		{
			name:        "ingest key — empty permissions",
			configValue: emptyPerms,
			config:      makeAPIKeyConfig(ctx, t, "ingest", emptyPerms),
			expectError: false,
		},
		{
			name:        "ingest key — create_datasets only",
			configValue: makePermissions(t, map[string]attr.Value{"create_datasets": types.BoolValue(true)}),
			config:      makeAPIKeyConfig(ctx, t, "ingest", makePermissions(t, map[string]attr.Value{"create_datasets": types.BoolValue(true)})),
			expectError: false,
		},
		{
			name:        "ingest key — config-only permission set to false",
			configValue: makePermissions(t, map[string]attr.Value{"manage_boards": types.BoolValue(false)}),
			config:      makeAPIKeyConfig(ctx, t, "ingest", makePermissions(t, map[string]attr.Value{"manage_boards": types.BoolValue(false)})),
			expectError: false,
		},
		{
			name:        "ingest key — manage_boards set to true",
			configValue: makePermissions(t, map[string]attr.Value{"manage_boards": types.BoolValue(true)}),
			config:      makeAPIKeyConfig(ctx, t, "ingest", makePermissions(t, map[string]attr.Value{"manage_boards": types.BoolValue(true)})),
			expectError: true,
		},
		{
			name: "ingest key — multiple config-only permissions set",
			configValue: makePermissions(t, map[string]attr.Value{
				"manage_boards":   types.BoolValue(true),
				"manage_triggers": types.BoolValue(true),
				"run_queries":     types.BoolValue(true),
			}),
			config: makeAPIKeyConfig(ctx, t, "ingest", makePermissions(t, map[string]attr.Value{
				"manage_boards":   types.BoolValue(true),
				"manage_triggers": types.BoolValue(true),
				"run_queries":     types.BoolValue(true),
			})),
			expectError: true,
		},
		{
			name:        "unknown key type — skipped",
			configValue: makePermissions(t, map[string]attr.Value{"manage_boards": types.BoolValue(true)}),
			config: func() tfsdk.Config {
				perms := makePermissions(t, map[string]attr.Value{"manage_boards": types.BoolValue(true)})
				permsTFValue, _ := perms.ToTerraformValue(ctx)
				return tfsdk.Config{
					Schema: apiKeyTestSchema,
					Raw: tftypes.NewValue(
						apiKeyTestSchema.Type().TerraformType(ctx),
						map[string]tftypes.Value{
							"type":        tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
							"permissions": permsTFValue,
						},
					),
				}
			}(),
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := validator.ListRequest{
				Path:           path.Root("permissions"),
				PathExpression: path.MatchRoot("permissions"),
				ConfigValue:    tc.configValue,
				Config:         tc.config,
			}
			resp := validator.ListResponse{}

			validation.APIKeyConfigKeyOnlyPermissions().ValidateList(ctx, req, &resp)

			assert.Equal(t, tc.expectError, resp.Diagnostics.HasError(),
				"unexpected error state: %s", resp.Diagnostics)
		})
	}
}
