resource "google_cloud_run_service" "frontend_service" {
  name     = "${local.env}-bytegolf-frontend"
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
      service_account_name = google_service_account.frontend.email
      container_concurrency = 20
      timeout_seconds = 30


      containers {
        image = local.frontend_image

        resources {
          limits = {
            memory = "128Mi"
            cpu    = "1000m"
          }
        }
      }
    }
  }


  traffic {
    percent         = 100
    latest_revision = true
  }
  autogenerate_revision_name = true
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

#################################################
#                SERVICE ACCOUNT                #
#################################################

resource "google_service_account" "frontend" {
  account_id   = "bg-frontend-${local.env}"
  display_name = "Frontend Service Account - ${local.env}"
}

# TODO: Find a better way to limit access to specific secrets that should be unaccessable 
# between environments
resource "google_project_iam_member" "frontend_secret_accessor" {
  project = local.project
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${resource.google_service_account.frontend.email}"
}
