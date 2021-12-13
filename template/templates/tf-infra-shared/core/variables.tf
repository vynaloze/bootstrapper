variable "tf_infra_repos" {
  description = "Terraform infrastructure repositories (containing root modules) to create, keyed by repo name"

  type = map(object({
    default_branch = string
    strict         = bool
  }))
}

variable "tfc_org_name" {
  description = "Terraform Cloud organization name"

  type = string
}

variable "repo_owner" {
  description = "Owner of the repositories"

  type = string
}

variable "repo_user" {
  description = "Repository access username (or empty)"

  type = string
}

variable "repo_password" {
  description = "Repository access token or password"

  type      = string
  sensitive = true
}

