package v2

import (
	"time"
)

type Environment struct {
	ID          string               `jsonapi:"primary,environments"`
	Name        string               `jsonapi:"attr,name"`
	Slug        string               `jsonapi:"attr,slug"`
	Description *string              `jsonapi:"attr,description,omitempty"`
	Color       *string              `jsonapi:"attr,color,omitempty"`
	Settings    *EnvironmentSettings `jsonapi:"attr,settings,omitempty"`
}

type EnvironmentSettings struct {
	DeleteProtected *bool `json:"delete_protected" jsonapi:"attr,delete_protected,omitempty"`
}

type Team struct {
	ID   string `jsonapi:"primary,teams"`
	Name string `jsonapi:"attr,name"`
	Slug string `jsonapi:"attr,slug"`
}

type Timestamps struct {
	CreatedAt time.Time `jsonapi:"attr,created,rfc3339,omitempty"`
	UpdatedAt time.Time `jsonapi:"attr,updated,rfc3339,omitempty"`
}

type AuthMetadata struct {
	ID         string      `jsonapi:"primary,api-keys"`
	Name       string      `jsonapi:"attr,name"`
	KeyType    string      `jsonapi:"attr,key_type"`
	Disabled   bool        `jsonapi:"attr,disabled"`
	Scopes     []string    `jsonapi:"attr,scopes"`
	Timestamps *Timestamps `jsonapi:"attr,timestamps"`
	Team       *Team       `jsonapi:"relation,team"`
}

const (
	APIKeyTypeIngest        = "ingest"
	APIKeyTypeConfiguration = "configuration"
)

type APIKey struct {
	ID          string             `jsonapi:"primary,api-keys,omitempty"`
	Name        *string            `jsonapi:"attr,name,omitempty"`
	KeyType     string             `jsonapi:"attr,key_type,omitempty"`
	Disabled    *bool              `jsonapi:"attr,disabled,omitempty"`
	Secret      string             `jsonapi:"attr,secret,omitempty"`
	Permissions *APIKeyPermissions `jsonapi:"attr,permissions,omitempty"`
	Timestamps  *Timestamps        `jsonapi:"attr,timestamps,omitempty"`
	Environment *Environment       `jsonapi:"relation,environment"`
}

// APIKeyPermissions are the permissions granted to an API Key.
//
//	See: https://docs.honeycomb.io/api/permissions for full details.
type APIKeyPermissions struct {
	CreateDatasets bool `jsonapi:"attr,create_datasets"`

	// the following are only supported on configuration keys
	ManageBoards       *bool `jsonapi:"attr,manage_boards,omitempty"`
	ManageColumns      *bool `jsonapi:"attr,manage_columns,omitempty"`
	ManageMarkers      *bool `jsonapi:"attr,manage_markers,omitempty"`
	ManageRecipients   *bool `jsonapi:"attr,manage_recipients,omitempty"`
	ManageSLOs         *bool `jsonapi:"attr,manage_slos,omitempty"`
	ManageTriggers     *bool `jsonapi:"attr,manage_triggers,omitempty"`
	ReadServiceMaps    *bool `jsonapi:"attr,read_service_maps,omitempty"`
	RunQueries         *bool `jsonapi:"attr,run_queries,omitempty"`
	SendEvents         *bool `jsonapi:"attr,send_events,omitempty"`
	VisibleTeamMembers *bool `jsonapi:"attr,visible_team_members,omitempty"`
}
