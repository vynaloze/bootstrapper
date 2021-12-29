locals {
  devops_storage_s3_bucket_name = "${local.account_id}-${var.environment}-devops-storage"
}

resource "aws_s3_bucket" "devops_storage" {
  bucket = local.devops_storage_s3_bucket_name
  acl    = "private"

  versioning {
    enabled = true
  }

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }

  tags = merge(local.common_tags, {
    Name = local.devops_storage_s3_bucket_name
  })
  #checkov:skip=CKV_AWS_18:do not require access logging
  #checkov:skip=CKV_AWS_52:do not require MFA delete
  #checkov:skip=CKV_AWS_144:do not require cross-region replication
  #checkov:skip=CKV_AWS_145:do not require SSE-KMS (SSE-S3 is enough)
}

resource "aws_s3_bucket_public_access_block" "devops_storage" {
  bucket = aws_s3_bucket.devops_storage.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}
