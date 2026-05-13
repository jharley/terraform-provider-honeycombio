package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAcc_APIKeyResource(t *testing.T) {
	ctx := context.Background()
	c := testAccV2Client(t)
	env := testAccEnvironment(ctx, t, c)

	t.Run("ingest key", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 testAccPreCheckV2API(t),
			ProtoV5ProviderFactories: testAccProtoV5MuxServerFactory,
			Steps: []resource.TestStep{
				{
					Config: testAccConfigBasicAPIKeyTest("test key", "false", env.ID),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureAPIKeyExists(t, "honeycombio_api_key.test"),
						resource.TestCheckResourceAttrSet("honeycombio_api_key.test", "id"),
						resource.TestCheckResourceAttrSet("honeycombio_api_key.test", "secret"),
						resource.TestCheckResourceAttrSet("honeycombio_api_key.test", "key"),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "name", "test key"),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "type", "ingest"),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "environment_id", env.ID),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "disabled", "false"),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "permissions.0.create_datasets", "true"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.manage_boards"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.manage_columns"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.manage_markers"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.manage_recipients"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.manage_slos"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.manage_triggers"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.read_service_maps"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.run_queries"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.send_events"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.visible_team_members"),
					),
				},
				{ // update the name and disabled state — permissions require replace for ingest keys
					Config: testAccConfigBasicAPIKeyTest("updated test key", "true", env.ID),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureAPIKeyExists(t, "honeycombio_api_key.test"),
						resource.TestCheckResourceAttrSet("honeycombio_api_key.test", "id"),
						resource.TestCheckResourceAttrSet("honeycombio_api_key.test", "secret"),
						resource.TestCheckResourceAttrSet("honeycombio_api_key.test", "key"),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "name", "updated test key"),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "type", "ingest"),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "environment_id", env.ID),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "disabled", "true"),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "permissions.0.create_datasets", "true"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.manage_boards"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.manage_columns"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.manage_markers"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.manage_recipients"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.manage_slos"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.manage_triggers"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.read_service_maps"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.run_queries"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.send_events"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.visible_team_members"),
					),
				},
			},
		})
	})

	t.Run("configuration key", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 testAccPreCheckV2API(t),
			ProtoV5ProviderFactories: testAccProtoV5MuxServerFactory,
			Steps: []resource.TestStep{
				{
					Config: testAccConfigBasicConfigurationKeyTest("test config key", env.ID),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureAPIKeyExists(t, "honeycombio_api_key.test"),
						resource.TestCheckResourceAttrSet("honeycombio_api_key.test", "id"),
						resource.TestCheckResourceAttrSet("honeycombio_api_key.test", "secret"),
						resource.TestCheckResourceAttrSet("honeycombio_api_key.test", "key"),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "name", "test config key"),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "type", "configuration"),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "environment_id", env.ID),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "disabled", "false"),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "permissions.0.create_datasets", "true"),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "permissions.0.manage_boards", "true"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.manage_columns"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.manage_markers"),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "permissions.0.manage_triggers", "true"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.read_service_maps"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.run_queries"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.send_events"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.visible_team_members"),
					),
				},
				{ // update permissions in-place — no replacement required for configuration keys
					Config: testAccConfigUpdatedConfigurationKeyTest("test config key", env.ID),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureAPIKeyExists(t, "honeycombio_api_key.test"),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "permissions.0.create_datasets", "false"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.manage_boards"),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "permissions.0.manage_markers", "true"),
						resource.TestCheckNoResourceAttr("honeycombio_api_key.test", "permissions.0.manage_triggers"),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "permissions.0.run_queries", "true"),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "permissions.0.visible_team_members", "true"),
					),
				},
			},
		})
	})

	t.Run("default values", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 testAccPreCheckV2API(t),
			ProtoV5ProviderFactories: testAccProtoV5MuxServerFactory,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
resource "honeycombio_api_key" "test" {
  name = "test key"
  type = "ingest"

  environment_id = "%s"
}`, env.ID),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureAPIKeyExists(t, "honeycombio_api_key.test"),
						resource.TestCheckResourceAttrSet("honeycombio_api_key.test", "id"),
						resource.TestCheckResourceAttrSet("honeycombio_api_key.test", "secret"),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "name", "test key"),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "type", "ingest"),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "environment_id", env.ID),
						resource.TestCheckResourceAttr("honeycombio_api_key.test", "disabled", "false"),
					),
				},
			},
		})
	})

	t.Run("config-key-only permissions rejected on ingest key", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 testAccPreCheckV2API(t),
			ProtoV5ProviderFactories: testAccProtoV5MuxServerFactory,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
resource "honeycombio_api_key" "test" {
  name = "invalid ingest key"
  type = "ingest"

  environment_id = "%s"

  permissions {
    manage_boards = true
  }
}`, env.ID),
					ExpectError: regexp.MustCompile(`Invalid Permission for API Key Type`),
				},
			},
		})
	})
}

func testAccConfigBasicAPIKeyTest(name, disabled, envID string) string {
	return fmt.Sprintf(`
resource "honeycombio_api_key" "test" {
  name     = "%s"
  type     = "ingest"
  disabled = %s

  environment_id = "%s"

  permissions {
    create_datasets = true
  }
}`, name, disabled, envID)
}

func testAccConfigBasicConfigurationKeyTest(name, envID string) string {
	return fmt.Sprintf(`
resource "honeycombio_api_key" "test" {
  name = "%s"
  type = "configuration"

  environment_id = "%s"

  permissions {
    create_datasets = true
    manage_boards   = true
    manage_triggers = true
  }
}`, name, envID)
}

func testAccConfigUpdatedConfigurationKeyTest(name, envID string) string {
	return fmt.Sprintf(`
resource "honeycombio_api_key" "test" {
  name = "%s"
  type = "configuration"

  environment_id = "%s"

  permissions {
    create_datasets     = false
    manage_markers      = true
    run_queries         = true
    visible_team_members = true
  }
}`, name, envID)
}

func testAccEnsureAPIKeyExists(t *testing.T, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}

		client := testAccV2Client(t)
		_, err := client.APIKeys.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to fetch created API key: %s", err)
		}

		return nil
	}
}
