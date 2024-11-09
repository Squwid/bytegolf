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
  backend_image = "us-central1-docker.pkg.dev/squid-cloud/bytegolf/backend:a3"
  frontend_url  = "play.byte.golf"
  backend_url   = "api.play.byte.golf"
  cookie_name   = "bg-token"
}

