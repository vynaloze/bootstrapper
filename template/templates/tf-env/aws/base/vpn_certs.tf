resource "tls_private_key" "ca" {
  algorithm = var.vpn_mtls_key_algorithm
  rsa_bits  = var.vpn_mtls_key_size
}

resource "tls_self_signed_cert" "ca" {
  key_algorithm   = var.vpn_mtls_key_algorithm
  private_key_pem = tls_private_key.ca.private_key_pem

  subject {
    common_name = "DevOps"
  }

  validity_period_hours = var.vpn_mtls_certificate_validity
  allowed_uses          = ["cert_signing", "crl_signing"]
  is_ca_certificate     = true
}

resource "tls_private_key" "server" {
  algorithm = var.vpn_mtls_key_algorithm
  rsa_bits  = var.vpn_mtls_key_size
}

resource "tls_cert_request" "server" {
  key_algorithm   = var.vpn_mtls_key_algorithm
  private_key_pem = tls_private_key.server.private_key_pem

  dns_names = ["${var.environment}-vpn-server"]

  subject {
    common_name = "${var.environment}-vpn-server"
  }
}

resource "tls_locally_signed_cert" "server" {
  cert_request_pem   = tls_cert_request.server.cert_request_pem
  ca_key_algorithm   = var.vpn_mtls_key_algorithm
  ca_private_key_pem = tls_private_key.ca.private_key_pem
  ca_cert_pem        = tls_self_signed_cert.ca.cert_pem

  validity_period_hours = var.vpn_mtls_certificate_validity
  allowed_uses          = ["key_encipherment", "digital_signature", "server_auth"]
  set_subject_key_id    = true
}

resource "tls_private_key" "client" {
  algorithm = var.vpn_mtls_key_algorithm
  rsa_bits  = var.vpn_mtls_key_size
}

resource "tls_cert_request" "client" {
  key_algorithm   = var.vpn_mtls_key_algorithm
  private_key_pem = tls_private_key.client.private_key_pem

  dns_names = ["${var.environment}-vpn-client"]

  subject {
    common_name = "${var.environment}-vpn-client"
  }
}

resource "tls_locally_signed_cert" "client" {
  cert_request_pem   = tls_cert_request.client.cert_request_pem
  ca_key_algorithm   = var.vpn_mtls_key_algorithm
  ca_private_key_pem = tls_private_key.ca.private_key_pem
  ca_cert_pem        = tls_self_signed_cert.ca.cert_pem

  validity_period_hours = var.vpn_mtls_certificate_validity
  allowed_uses          = ["digital_signature", "client_auth"]
  set_subject_key_id    = true
}

data "archive_file" "vpn_config" {
  type        = "zip"
  output_path = "${path.module}/${var.environment}.zip"

  source {
    content  = tls_self_signed_cert.ca.cert_pem
    filename = "${var.environment}-ca.crt"
  }

  source {
    content  = tls_private_key.client.private_key_pem
    filename = "${var.environment}-client.key"
  }

  source {
    content  = tls_locally_signed_cert.client.cert_pem
    filename = "${var.environment}-client.crt"
  }

  source {
    content = templatefile("${path.module}/vpn_config.tpl", {
      endpoint            = trimprefix(aws_ec2_client_vpn_endpoint.this.dns_name, "*.")
      vpc_dns_resolver_ip = cidrhost(var.vpc_cidr, 2)
      dns_domain          = local.domain_name
      environment         = var.environment
    })
    filename = "${var.environment}.ovpn"
  }
}

resource "aws_s3_bucket_object" "devops_storage_vpn_config" {
  bucket = aws_s3_bucket.devops_storage.id
  key    = "vpn/${basename(data.archive_file.vpn_config.output_path)}"
  source = data.archive_file.vpn_config.output_path
}
