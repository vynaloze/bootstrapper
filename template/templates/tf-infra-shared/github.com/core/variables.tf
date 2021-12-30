variable "tf_infra_repos" {
  description = "Terraform infrastructure repositories (containing root modules) to create, keyed by repo name"

  type = map(object({
    modules = list(string)

    default_branch = string
    strict         = bool
    build_checks   = list(string)
  }))
}

variable "tf_module_repos" {
  description = "Terraform module repositories (containing reusable modules) to create, keyed by repo name"

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

variable "tfc_organization" {
  description = "Terraform Cloud organization name"

  type = string
}
