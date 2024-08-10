provider "google" {
  project = "squid-cloud"
  region  = "us-central1"
}

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "5.38.0"
    }
  }

  backend "gcs" {
    bucket = "bytegolf-tf-state"
    prefix = "terraform/state"
  }
}

locals {
  project       = "squid-cloud"
  backend_image = "us-central1-docker.pkg.dev/squid-cloud/bytegolf/backend:a2" # TODO: Update to local registry image
  frontend_url  = "byte.golf"
  backend_url   = "api.byte.golf"
  cookie_name   = "bg-token"
}

resource "google_storage_bucket" "frontend_bucket" {
  name          = "bytegolf-fe"
  location      = "US-CENTRAL1"
  force_destroy = false

  uniform_bucket_level_access = true
}