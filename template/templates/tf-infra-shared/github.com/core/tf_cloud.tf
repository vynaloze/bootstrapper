resource "tfe_workspace" "this" {
  for_each = toset(flatten([for k, v in var.tf_infra_repos : [for m in v.modules : "${k}-${m}"]]))

  organization   = var.tfc_organization
  name           = each.key
  execution_mode = "local"
}
