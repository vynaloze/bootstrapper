variable "tf_infra_repos" {
  description = "Terraform infrastructure repositories (containing root modules) to create, keyed by repo name"

  type = map(object({
    default_branch = string
    strict         = bool
    build_checks   = list(string)
  }))
}

variable "misc_repos" {
  description = "Miscellaneous repositories to create, keyed by repo name"

  type = map(object({
    default_branch = string
    strict         = bool
    build_checks   = list(string)
  }))
}

variable "repo_owner" {
  description = "Owner of the repositories"

  type = string
}

# tflint-ignore: terraform_unused_declarations
variable "repo_user" {
  description = "Repository access username (or empty)"
  default     = ""

  type = string
}

variable "repo_password" {
  description = "Repository access token or password"

  type      = string
  sensitive = true
}

