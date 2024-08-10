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
  backend_image = "crccheck/hello-world"
  frontend_url  = "byte.golf"
  backend_url   = "api.byte.golf"
  cookie_name   = "bg-token"
}

