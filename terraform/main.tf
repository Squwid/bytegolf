provider "google" {
  project = "squid-cloud"
  region  = "us-central1"
}

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.47.0"
    }
  }

  backend "gcs" {
    bucket = "bytegolf-tf-state"
    prefix = "terraform/state"
  }
}

locals {
  env            = terraform.workspace
  project        = "squid-cloud"
  frontend_image = "gcr.io/squid-cloud/bytegolf-frontend@sha256:bed79ebd2680c0ad85d05e9a7cddf0ab81fb729807a78de6de4bd04b1cf0476a"
  backend_image  = "gcr.io/squid-cloud/bytegolf-backend@sha256:16d4b2779db99290c8c3c4057f6d6c77e1d149be907e452aecab8d40dd3d2cd6"

  frontend_addr  = terraform.workspace == "prod" ? "https://byte.golf" : "https://${terraform.workspace}.byte.golf"
  backend_addr   = terraform.workspace == "prod" ? "https://api.byte.golf" : "https://${terraform.workspace}.api.byte.golf"
  frontend_url   = terraform.workspace == "prod" ? "byte.golf" : "${terraform.workspace}.byte.golf"
  backend_url    = terraform.workspace == "prod" ? "api.byte.golf" : "${terraform.workspace}.api.byte.golf"
  cookie_name    = terraform.workspace == "prod" ? "bg-token" : "bg-token-${terraform.workspace}"
}