resource "google_storage_bucket" "frontend_bucket" {
  name                        = "bytegolf-frontend-${local.env}"
  location                    = "us-central1"
  storage_class               = "STANDARD"
  uniform_bucket_level_access = false

  website {
    main_page_suffix = "index.html"
    not_found_page = "index.html"
  }

  force_destroy = true # Turn this to false if you dont want terraform to delete objects in the bucket on delete
}

resource "google_storage_bucket_iam_binding" "public_access_to_bucket" {
  bucket = google_storage_bucket.frontend_bucket.name
  role = "roles/storage.objectViewer"
  members = [
    "allUsers"
  ]
}

#################################################
#                Backend Bucket                 #
#################################################

resource "google_compute_backend_bucket" "frontend_bucket" {
  name        = "${local.env}-frontend-bucket"
  description = "Backend bucket for Bytegolf Frontend ${local.env}"
  bucket_name = google_storage_bucket.frontend_bucket.name
  enable_cdn  = false
}