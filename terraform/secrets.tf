resource "google_secret_manager_secret" "github_client_secret" {
  secret_id = "bg_GITHUB_CLIENT"

  labels = {
    app : "bytegolf"
    env : "prod"
  }

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret" "github_secret_secret" {
  secret_id = "bg_GITHUB_SECRET"

  labels = {
    app : "bytegolf"
    env : "prod"
  }

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret" "github_state_secret" {
  secret_id = "bg_GITHUB_STATE"

  labels = {
    app : "bytegolf"
    env : "prod"
  }

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret" "jwt_secret" {
  secret_id = "bg_JWT_SECRET"

  labels = {
    app : "bytegolf"
    env : "prod"
  }

  replication {
    auto {}
  }
}