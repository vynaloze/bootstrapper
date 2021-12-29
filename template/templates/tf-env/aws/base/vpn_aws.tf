resource "aws_acm_certificate" "server" {
  private_key       = tls_private_key.server.private_key_pem
  certificate_body  = tls_locally_signed_cert.server.cert_pem
  certificate_chain = tls_self_signed_cert.ca.cert_pem

  tags = merge(local.common_tags, {
    Name = "${var.environment}-vpn-server",
  })

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_acm_certificate" "client" {
  private_key       = tls_private_key.client.private_key_pem
  certificate_body  = tls_locally_signed_cert.client.cert_pem
  certificate_chain = tls_self_signed_cert.ca.cert_pem

  tags = merge(local.common_tags, {
    Name = "${var.environment}-vpn-client",
  })

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_ec2_client_vpn_endpoint" "this" {
  client_cidr_block      = var.client_cidr_block
  server_certificate_arn = aws_acm_certificate.server.arn

  authentication_options {
    type                       = "certificate-authentication"
    root_certificate_chain_arn = aws_acm_certificate.client.arn
    # Even though it should be a _server_ certificate here ("If the client certificate has been issued by the same CA
    # as the server certificate, you can use the server certificate ARN for the client certificate ARN.":
    # https://docs.aws.amazon.com/vpn/latest/clientvpn-admin/cvpn-working-endpoints.html), we use _client_ certificate
    # to trigger recreate of the whole VPN configuration in case of triggering certificate recreation in Terraform.
    # Otherwise, old certificate would be still valid unless expired (after 10y) or revoked using CRL
    # (manually: CRL generation is not implemented in tls provider and CRL import is not implemented in aws provider).
  }

  connection_log_options {
    enabled = false
  }

  split_tunnel = true

  tags = merge(
    {
      Name = local.vpn_name
    },
    local.common_tags
  )
}

resource "aws_ec2_client_vpn_network_association" "this" {
  client_vpn_endpoint_id = aws_ec2_client_vpn_endpoint.this.id
  subnet_id              = module.vpc.private_subnets[0]
  security_groups        = [aws_security_group.vpn_access.id]
}

resource "aws_security_group" "vpn_access" {
  name        = local.vpn_name
  description = "access to VPN endpoint - available on the internet"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port   = 0
    protocol    = "-1"
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    from_port   = 0
    protocol    = "-1"
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.common_tags

  #checkov:skip=CKV2_AWS_5:this SG is attached to cVPN - skipping rule "Ensure that SGs are attached to EC2 instances or ENIs"
}

resource "aws_ec2_client_vpn_authorization_rule" "internet_access" {
  client_vpn_endpoint_id = aws_ec2_client_vpn_endpoint.this.id
  target_network_cidr    = "0.0.0.0/0"
  authorize_all_groups   = true
}
