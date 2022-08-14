resource "google_cloud_run_service" "backend_service" {
  name     = "${local.env}-bytegolf-backend"
  location = "us-central1"

  metadata {
    annotations = {
      "run.googleapis.com/ingress" : "internal-and-cloud-load-balancing"
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

      containers {
        image = local.backend_container

        resources {
          # requests = {
          #   memory = "256Mi"
          #   cpu = "1000m"
          # }

          limits = {
            memory = "256Mi"
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

data "google_iam_policy" "noauth" {
  binding {
    role    = "roles/run.invoker"
    members = ["allUsers"]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location = google_cloud_run_service.backend_service.location
  project  = google_cloud_run_service.backend_service.project
  service  = google_cloud_run_service.backend_service.name

  policy_data = data.google_iam_policy.noauth.policy_data
}

#################################################
#                      NEG                      #
#################################################

resource "google_compute_backend_service" "backend_service" {
  provider    = google
  name        = "${local.env}-bytegolf-backend-service"
  description = "Backend service for Bytegolf Backend ${local.env}"
  enable_cdn  = false

  backend {
    group = google_compute_region_network_endpoint_group.backend_neg.id
  }
}

resource "google_compute_region_network_endpoint_group" "backend_neg" {
  name                  = "${local.env}-backend-neg"
  network_endpoint_type = "SERVERLESS"
  region                = "us-central1"

  cloud_run {
    service = google_cloud_run_service.backend_service.name
  }
}

#################################################
#                SERVICE ACCOUNT                #
#################################################

resource "google_service_account" "backend" {
  account_id   = "bg-backend-${local.env}"
  display_name = "Backend Service Account - ${local.env}"
}

resource "google_project_iam_binding" "backend_secret_accessor" {
  project = local.project
  role    = "roles/secretmanager.secretAccessor"

  members = ["serviceAccount:${resource.google_service_account.backend.email}"]
}

resource "google_project_iam_binding" "firestore_service_agent" {
  project = local.project
  role    = "roles/firestore.serviceAgent"

  members = ["serviceAccount:${resource.google_service_account.backend.email}"]
}
