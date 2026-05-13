package provider

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/honeycombio/terraform-provider-honeycombio/client"
	v2client "github.com/honeycombio/terraform-provider-honeycombio/client/v2"
	"github.com/honeycombio/terraform-provider-honeycombio/internal/helper"
	"github.com/honeycombio/terraform-provider-honeycombio/internal/models"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ ephemeral.EphemeralResource              = &apiKeyEphemeralResource{}
	_ ephemeral.EphemeralResourceWithConfigure = &apiKeyEphemeralResource{}
	_ ephemeral.EphemeralResourceWithClose     = &apiKeyEphemeralResource{}
)

type apiKeyEphemeralResource struct {
	client *v2client.Client
}

func NewAPIKeyEphemeralResource() ephemeral.EphemeralResource {
	return &apiKeyEphemeralResource{}
}

func (*apiKeyEphemeralResource) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key"
}

func (r *apiKeyEphemeralResource) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	w := getClientFromEphemeralResourceRequest(&req)
	if w == nil {
		return
	}

	c, err := w.V2Client()
	if err != nil || c == nil {
		resp.Diagnostics.AddError("Failed to configure client", err.Error())
		return
	}
	r.client = c
}

func (*apiKeyEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a short-lived configuration API key that is deleted when the Terraform operation completes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the API Key.",
				MarkdownDescription: "The ID of the API Key.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the API key.",
				MarkdownDescription: "The name of the API key.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"type": schema.StringAttribute{
				Required:            true,
				Description:         "The type of API key. Currently only 'configuration' is supported.",
				MarkdownDescription: "The type of API key. Currently only `configuration` is supported.",
				Validators: []validator.String{
					stringvalidator.OneOf(v2client.APIKeyTypeConfiguration),
				},
			},
			"environment_id": schema.StringAttribute{
				Required:            true,
				Description:         "The Environment ID the API key is scoped to.",
				MarkdownDescription: "The Environment ID the API key is scoped to.",
			},
			"disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Whether the API key is disabled.",
				MarkdownDescription: "Whether the API key is disabled.",
			},
			"key": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         "The API key formatted for use based on its type.",
				MarkdownDescription: "The API key formatted for use based on its type.",
			},
			"secret": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         "The secret portion of the API Key.",
				MarkdownDescription: "The secret portion of the API Key.",
			},
		},
		Blocks: map[string]schema.Block{
			"permissions": schema.ListNestedBlock{
				Description:         "A configuration block setting what actions the API key can perform.",
				MarkdownDescription: "A configuration block setting what actions the API key can perform.",
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"create_datasets": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Allow this key to create missing datasets when sending telemetry.",
							MarkdownDescription: "Allow this key to create missing datasets when sending telemetry.",
						},
						"manage_boards": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Allow this key to manage boards.",
							MarkdownDescription: "Allow this key to manage boards.",
						},
						"manage_columns": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Allow this key to manage columns.",
							MarkdownDescription: "Allow this key to manage columns.",
						},
						"manage_markers": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Allow this key to manage markers.",
							MarkdownDescription: "Allow this key to manage markers.",
						},
						"manage_recipients": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Allow this key to manage recipients.",
							MarkdownDescription: "Allow this key to manage recipients.",
						},
						"manage_slos": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Allow this key to manage SLOs.",
							MarkdownDescription: "Allow this key to manage SLOs.",
						},
						"manage_triggers": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Allow this key to manage triggers.",
							MarkdownDescription: "Allow this key to manage triggers.",
						},
						"read_service_maps": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Allow this key to read service maps.",
							MarkdownDescription: "Allow this key to read service maps.",
						},
						"run_queries": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Allow this key to run queries.",
							MarkdownDescription: "Allow this key to run queries.",
						},
						"send_events": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Allow this key to send events.",
							MarkdownDescription: "Allow this key to send events.",
						},
						"visible_team_members": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Allow this key to view team members.",
							MarkdownDescription: "Allow this key to view team members.",
						},
					},
				},
			},
		},
	}
}

func (r *apiKeyEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var plan models.APIKeyResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newKey := &v2client.APIKey{
		Name:        plan.Name.ValueStringPointer(),
		KeyType:     plan.Type.ValueString(),
		Environment: &v2client.Environment{ID: plan.EnvironmentID.ValueString()},
		Disabled:    plan.Disabled.ValueBoolPointer(),
		Permissions: expandAPIKeyPermissions(ctx, plan.Type.ValueString(), plan.Permissions, &resp.Diagnostics),
	}
	if resp.Diagnostics.HasError() {
		return
	}

	key, err := r.client.APIKeys.Create(ctx, newKey)
	if helper.AddDiagnosticOnError(&resp.Diagnostics, "Creating Honeycomb Ephemeral API Key", err) {
		return
	}

	// Store the key ID in private state for Close to delete.
	// SetKey requires valid JSON, so marshal the ID first.
	idJSON, err := json.Marshal(key.ID)
	if helper.AddDiagnosticOnError(&resp.Diagnostics, "Creating Honeycomb Ephemeral API Key", err) {
		return
	}
	resp.Diagnostics.Append(resp.Private.SetKey(ctx, "id", idJSON)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.APIKeyResourceModel
	state.ID = types.StringValue(key.ID)
	state.Name = types.StringValue(*key.Name)
	state.Type = types.StringValue(key.KeyType)
	state.EnvironmentID = types.StringValue(key.Environment.ID)
	state.Disabled = types.BoolValue(*key.Disabled)
	state.Secret = types.StringValue(key.Secret)

	switch key.KeyType {
	case v2client.APIKeyTypeConfiguration:
		state.Key = types.StringValue(key.Secret)
	default:
		resp.Diagnostics.AddError(
			"Unknown API Key Type",
			"API Key Type "+key.KeyType+" is not supported.",
		)
	}

	if !plan.Permissions.IsNull() {
		state.Permissions = flattenAPIKeyPermissions(ctx, key.KeyType, key.Permissions, &resp.Diagnostics)
	} else {
		state.Permissions = types.ListNull(types.ObjectType{AttrTypes: models.APIKeyPermissionsAttrType})
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, state)...)
}

func (r *apiKeyEphemeralResource) Close(ctx context.Context, req ephemeral.CloseRequest, resp *ephemeral.CloseResponse) {
	idBytes, diags := req.Private.GetKey(ctx, "id")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var keyID string
	if err := json.Unmarshal(idBytes, &keyID); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading API Key ID",
			"Could not read the API Key ID from private state: "+err.Error(),
		)
		return
	}

	err := r.client.APIKeys.Delete(ctx, keyID)
	var detailedErr client.DetailedError
	if err != nil {
		if errors.As(err, &detailedErr) {
			// if not found consider it deleted -- so don't error
			if !detailedErr.IsNotFound() {
				resp.Diagnostics.Append(helper.NewDetailedErrorDiagnostic(
					"Error Deleting Honeycomb Ephemeral API Key",
					&detailedErr,
				))
			}
		} else {
			resp.Diagnostics.AddError(
				"Error Deleting Honeycomb Ephemeral API Key",
				"Could not delete API Key ID "+keyID+": "+err.Error(),
			)
		}
	}
}
