resource "google_secret_manager_secret" "github_client_secret" {
  secret_id = "${local.env}_GITHUB_CLIENT"

  labels = {
    app : "bytegolf"
    env : "${local.env}"
  }

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret" "github_secret_secret" {
  secret_id = "${local.env}_GITHUB_SECRET"

  labels = {
    app : "bytegolf"
    env : "${local.env}"
  }

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret" "github_state_secret" {
  secret_id = "${local.env}_GITHUB_STATE"

  labels = {
    app : "bytegolf"
    env : "${local.env}"
  }

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret" "jwt_secret" {
  secret_id = "${local.env}_JWT_SECRET"

  labels = {
    app : "bytegolf"
    env : "${local.env}"
  }

  replication {
    automatic = true
  }
}