resource "honeycombio_environment" "example" {
  name  = "example"
  color = "purple"
}

ephemeral "honeycombio_api_key" "example" {
  name = "example tf key"
  type = "configuration"

  environment_id = honeycombio_environment.example.id

  permissions {
    create_datasets   = true
    manage_columns    = true
    manage_boards     = true
    manage_recipients = true
    manage_triggers   = true
    manage_slos       = true
  }
}

provider "honeycombio" {
  alias = "example"

  api_key = ephemeral.honeycombio_api_key.example.key
}

// now we can use our key-scoped provider to make resources inside our new environment
resource "honeycombio_dataset" "example" {
  provider = honeycombio.example

  name = "example"
}
