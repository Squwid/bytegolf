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
  env               = terraform.workspace
  project           = "squid-cloud"
  backend_container = "gcr.io/squid-cloud/bytegolf-backend@sha256:69a9ee3bc4c0d27346973809587ece52ec639938151d9a4573602e880a9f211a"
  frontend_addr     = terraform.workspace == "prod" ? "https://byte.golf" : "https://${terraform.workspace}.byte.golf"
  backend_addr      = terraform.workspace == "prod" ? "https://api.byte.golf" : "https://${terraform.workspace}.api.byte.golf"
  cookie_name       = "bg-token"
}