ephemeral "honeycombio_api_key" "example" {
  name = "ephemeral config key"
  type = "configuration"

  environment_id = var.environment_id

  permissions {
    create_datasets = true
    run_queries     = true
  }
}
