package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type APIKeyResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
	EnvironmentID types.String `tfsdk:"environment_id"`
	Disabled      types.Bool   `tfsdk:"disabled"`
	Permissions   types.List   `tfsdk:"permissions"` // APIKeyPermissionModel
	Secret        types.String `tfsdk:"secret"`
	Key           types.String `tfsdk:"key"`
}

type APIKeyPermissionModel struct {
	CreateDatasets     types.Bool `tfsdk:"create_datasets"`
	ManageBoards       types.Bool `tfsdk:"manage_boards"`
	ManageColumns      types.Bool `tfsdk:"manage_columns"`
	ManageMarkers      types.Bool `tfsdk:"manage_markers"`
	ManageRecipients   types.Bool `tfsdk:"manage_recipients"`
	ManageSLOs         types.Bool `tfsdk:"manage_slos"`
	ManageTriggers     types.Bool `tfsdk:"manage_triggers"`
	ReadServiceMaps    types.Bool `tfsdk:"read_service_maps"`
	RunQueries         types.Bool `tfsdk:"run_queries"`
	SendEvents         types.Bool `tfsdk:"send_events"`
	VisibleTeamMembers types.Bool `tfsdk:"visible_team_members"`
}

var APIKeyPermissionsAttrType = map[string]attr.Type{
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
