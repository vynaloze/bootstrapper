resource "tfe_workspace" "this" {
  for_each = var.tf_infra_repos

  organization   = var.tfc_organization
  name           = each.key
  execution_mode = "local"
}
