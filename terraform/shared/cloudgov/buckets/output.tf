output "ids" {
  value = {
    extract = module.s3_extract.bucket_id
    fetch   = module.s3_fetch.bucket_id
    meta    = module.s3_meta.bucket_id
    backup  = module.s3_backup.bucket_id
  }
}
