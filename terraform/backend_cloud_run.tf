resource "google_cloud_run_service" "backend_service" {
  name     = "${local.env}-bytegolf-backend"
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
      service_account_name = google_service_account.backend.email
      container_concurrency = 20
      timeout_seconds = 30

      containers {
        image = local.backend_image

        resources {
          # requests = {
          #   memory = "256Mi"
          #   cpu = "1000m"
          # }

          limits = {
            memory = "128Mi"
            cpu    = "1000m"
          }
        }

        ports {
          container_port = "8080"
        }

        env {
          name  = "BG_ENV"
          value = local.env
        }

        env {
          name  = "BG_FRONTEND_ADDR"
          value = local.frontend_addr
        }

        env {
          name  = "BG_BACKEND_ADDR"
          value = local.backend_addr
        }

        env {
          name  = "BG_COOKIE_NAME"
          value = local.cookie_name
        }

        env {
          name = "GITHUB_CLIENT"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.github_client_secret.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "GITHUB_SECRET"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.github_secret_secret.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "GITHUB_STATE"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.github_state_secret.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "JDOODLE_CLIENT"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.jdoodle_client_secret.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "JDOODLE_SECRET"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.jdoodle_secret_secret.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "JWT_SECRET"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.jwt_secret.secret_id
              key  = "latest"
            }
          }
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
    google_secret_manager_secret.jdoodle_client_secret,
    google_secret_manager_secret.jdoodle_secret_secret,
    google_secret_manager_secret.jwt_secret,
    google_service_account.backend
  ]
}

data "google_iam_policy" "noauth_backend" {
  binding {
    role    = "roles/run.invoker"
    members = ["allUsers"]
  }
}

resource "google_cloud_run_service_iam_policy" "backend_noauth" {
  location = google_cloud_run_service.backend_service.location
  project  = google_cloud_run_service.backend_service.project
  service  = google_cloud_run_service.backend_service.name

  policy_data = data.google_iam_policy.noauth_backend.policy_data
}

resource "google_cloud_run_domain_mapping" "backend" {
  location = "us-central1"
  name     = local.backend_url

  metadata {
    namespace = local.project
  }

  spec {
    route_name = google_cloud_run_service.backend_service.name
  }

  depends_on = [
    google_cloud_run_service.backend_service
  ]
}


#################################################
#                SERVICE ACCOUNT                #
#################################################

resource "google_service_account" "backend" {
  account_id   = "bg-backend-${local.env}"
  display_name = "Backend Service Account - ${local.env}"
}

# TODO: Find a better way to limit access to specific secrets that should be unaccessable 
# between environments
resource "google_project_iam_member" "backend_secret_accessor" {
  project = local.project
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${resource.google_service_account.backend.email}"
}

resource "google_project_iam_member" "backend_firebase_admin" {
  project = local.project
  role    = "roles/firebase.admin"
  member  = "serviceAccount:${resource.google_service_account.backend.email}"
}
