package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/honeycombio/terraform-provider-honeycombio/internal/helper/test"
)

func TestAcc_APIKeyEphemeralResource(t *testing.T) {
	ctx := context.Background()
	c := testAccV2Client(t)
	env := testAccEnvironment(ctx, t, c)
	keyName := test.RandomStringWithPrefix("test.ephemeral.", 20)

	// Track key IDs created during the test to verify Close deletes them
	var createdKeyIDs []string

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0), // ephemeral resources not supported prior to Terraform 1.10.x
		},
		PreCheck:                 testAccPreCheckV2API(t),
		ProtoV5ProviderFactories: testAccProtoV5MuxServerFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
ephemeral "honeycombio_api_key" "test" {
  name = "%s"
  type = "configuration"

  environment_id = "%s"

  permissions {
    create_datasets = true
    manage_boards   = true
  }
}`, keyName, env.ID),
				Check: func(_ *terraform.State) error {
					// Ephemeral resources don't appear in state, so we query
					// the API directly to capture the key ID for post-test cleanup verification.
					pager, err := c.APIKeys.List(ctx)
					if err != nil {
						return fmt.Errorf("failed to list API keys: %s", err)
					}
					for pager.HasNext() {
						keys, err := pager.Next(ctx)
						if err != nil {
							return fmt.Errorf("failed to fetch API keys page: %s", err)
						}
						for _, k := range keys {
							if k.Name != nil && *k.Name == keyName {
								createdKeyIDs = append(createdKeyIDs, k.ID)
							}
						}
					}
					return nil
				},
			},
		},
	})

	// After resource.Test returns, Close should have been called.
	// Verify the key was deleted.
	for _, id := range createdKeyIDs {
		_, err := c.APIKeys.Get(ctx, id)
		if err == nil {
			t.Errorf("expected ephemeral API key %s to be deleted after Close, but it still exists", id)
			// Clean up the leaked key
			c.APIKeys.Delete(ctx, id) //nolint:errcheck
		}
	}
}
