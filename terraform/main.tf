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
  backend_container = "gcr.io/squid-cloud/bytegolf-backend@sha256:dfb2710e0cf0fcc86c5a8052b337679aef441e28238cd005e7802b3e7fc7bcd5"
  frontend_addr     = terraform.workspace == "prod" ? "https://byte.golf" : "https://${terraform.workspace}.byte.golf"
  backend_addr      = terraform.workspace == "prod" ? "https://api.byte.golf" : "https://${terraform.workspace}.api.byte.golf"
  cookie_name       = terraform.workspace == "prod" ? "bg-token" : "bg-token-${terraform.workspace}"
}