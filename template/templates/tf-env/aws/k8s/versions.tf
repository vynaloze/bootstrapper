terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 3.16"
    }
    local = {
      source  = "hashicorp/local"
      version = ">= 1.4"
    }
    null = {
      source  = "hashicorp/null"
      version = ">= 2.1"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 2.3"
    }
    template = {
      source  = "hashicorp/template"
      version = ">= 2.1"
    }
  }
  required_version = ">= 0.13"
}
