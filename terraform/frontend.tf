resource "google_cloud_run_service" "frontend_service" {
  name     = "bytegolf-frontend"
  location = "us-central1"

  metadata {
    annotations = {
      "run.googleapis.com/ingress" : "all"
    }
  }

  template {
    metadata {
      annotations = {
        "autoscaling.knative.dev/minScale" = "0"
        "autoscaling.knative.dev/maxScale" = "3"
      }
    }

    spec {
      service_account_name  = google_service_account.backend.email
      container_concurrency = 20
      timeout_seconds       = 30

      containers {
        image = "squwid/bgcs-site-proxy:v0.3"

        resources {
          limits = {
            memory = "256Mi"
            cpu    = "1000m"
          }
        }

        ports {
          container_port = "8000"
        }

        env {
          name  = "BGCS_BUCKET"
          value = google_storage_bucket.frontend_bucket.name
        }

        env {
          name  = "BGCS_NOT_FOUND_FILE"
          value = "index.html"
        }

        env {
          name  = "BGCS_DEFAULT_FILE"
          value = "index.html"
        }
      }
    }
  }

  # Split Traffic - https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_run_service#example-usage---cloud-run-service-traffic-split
  traffic {
    percent         = 100
    latest_revision = true
  }
  autogenerate_revision_name = true

  depends_on = [
    google_secret_manager_secret.github_client_secret,
    google_secret_manager_secret.github_secret_secret,
    google_secret_manager_secret.github_state_secret,
    google_secret_manager_secret.jwt_secret,
    google_service_account.backend
  ]
}

data "google_iam_policy" "noauth_frontend" {
  binding {
    role    = "roles/run.invoker"
    members = ["allUsers"]
  }
}

resource "google_cloud_run_service_iam_policy" "frontend_noauth" {
  location = google_cloud_run_service.frontend_service.location
  project  = google_cloud_run_service.frontend_service.project
  service  = google_cloud_run_service.frontend_service.name

  policy_data = data.google_iam_policy.noauth_frontend.policy_data
}

resource "google_cloud_run_domain_mapping" "frontend" {
  location = "us-central1"
  name     = local.frontend_url

  metadata {
    namespace = local.project
  }

  spec {
    route_name = google_cloud_run_service.frontend_service.name
  }

  depends_on = [
    google_cloud_run_service.frontend_service
  ]
}


resource "google_storage_bucket" "frontend_bucket" {
  name          = "bytegolf-fe"
  location      = "US-CENTRAL1"
  force_destroy = false

  uniform_bucket_level_access = true
}

#################################################
#                SERVICE ACCOUNT                #
#################################################

resource "google_service_account" "frontend" {
  account_id   = "bg-frontend"
  display_name = "Bytegolf Frontend Service Account"
}

resource "google_storage_bucket_iam_member" "bucket_object_viewer" {
  bucket = google_storage_bucket.frontend_bucket.name
  role   = "roles/storage.objectViewer"
  member = "serviceAccount:${google_service_account.frontend.email}"
}