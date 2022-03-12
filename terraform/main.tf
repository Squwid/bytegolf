provider "google" {
  project = "squid-cloud"
  region  = "us-central1"
}

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.13.0"
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
  backend_container = "gcr.io/squid-cloud/bytegolf-backend@sha256:dff7b98627d9784eaa050aa4a0c04daa638647d9c093b42b0100135a0a43bc46"
  frontend_addr     = "https://byte.golf"
  backend_addr      = "https://api.byte.golf"
  cookie_name       = "bg-token"
}