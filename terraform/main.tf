provider "google" {
  project = "squid-cloud"
  region  = "us-central1"
}

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.31.0"
    }
  }

  backend "gcs" {
    bucket = "bytegolf-tf-state"
    prefix = "terraform/state"
  }
}

locals {
  project       = "squid-cloud"
  backend_image = "gcr.io/squid-cloud/bytegolf-backend@sha256:16d4b2779db99290c8c3c4057f6d6c77e1d149be907e452aecab8d40dd3d2cd6"
  frontend_url  = "byte.golf"
  backend_url   = "api.byte.golf"
  cookie_name   = "bg-token"
}