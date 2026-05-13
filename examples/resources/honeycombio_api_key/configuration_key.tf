resource "honeycombio_api_key" "development" {
  name = "Dev Terraform Key"
  type = "configuration"

  environment_id = var.environment_id

  permissions {
    create_datasets = true
    manage_boards   = true
    manage_columns  = true
    manage_slos     = true
    manage_triggers = true
  }
}
