
data "cloudfoundry_domain" "public" {
  name = "app.cloud.gov"
}

data "cloudfoundry_space" "app_space" {
  # name = var.cf_space_guid
  name = local.cf_space
  org  = var.cf_org_guid
}


######################################
# DATABASES
######################################
module "databases" {
  source              = "../shared/cloudgov/databases"
  cf_org              = local.cf_org
  cf_space            = local.cf_space
  queue_db_plan_name  = "micro-psql"
  search_db_plan_name = "micro-psql"
  work_db_plan_name   = "micro-psql"
}


######################################
# S3 BUCKETS
######################################
module "buckets" {
  source   = "../shared/cloudgov/buckets"
  cf_org   = local.cf_org
  cf_space = local.cf_space
}

######################################
# SERVICES
######################################
module "admin" {
  source = "../shared/services/admin"
  # disk_quota = 256
  # memory = 128
  # instances = 1
  space_name = data.cloudfoundry_space.app_space.name
  app_space_id = data.cloudfoundry_space.app_space.id
  domain_id = data.cloudfoundry_domain.public.id
  databases = module.databases.ids
  buckets = module.buckets.ids
}

module "entree" {
  source = "../shared/services/entree"
  # disk_quota = 256
  # memory = 128
  # instances = 1
  space_name = data.cloudfoundry_space.app_space.name
  app_space_id = data.cloudfoundry_space.app_space.id
  domain_id = data.cloudfoundry_domain.public.id
  databases = module.databases.ids
}

module "extract" {
  source = "../shared/services/extract"
  # disk_quota = 256
  # memory = 128
  # instances = 1
  space_name = data.cloudfoundry_space.app_space.name
  app_space_id = data.cloudfoundry_space.app_space.id
  domain_id = data.cloudfoundry_domain.public.id
  databases = module.databases.ids
  buckets = module.buckets.ids
}

module "fetch" {
  source = "../shared/services/fetch"
  # disk_quota = 256
  # memory = 128
  # instances = 1
  space_name = data.cloudfoundry_space.app_space.name
  app_space_id = data.cloudfoundry_space.app_space.id
  domain_id = data.cloudfoundry_domain.public.id
  databases = module.databases.ids
  buckets = module.buckets.ids
}

module "pack" {
  source = "../shared/services/pack"
  # disk_quota = 256
  # memory = 128
  # instances = 1
  space_name = data.cloudfoundry_space.app_space.name
  app_space_id = data.cloudfoundry_space.app_space.id
  domain_id = data.cloudfoundry_domain.public.id
  databases = module.databases.ids
  buckets = module.buckets.ids
}

module "walk" {
  source = "../shared/services/walk"
  # disk_quota = 256
  # memory = 128
  # instances = 1
  space_name = data.cloudfoundry_space.app_space.name
  app_space_id = data.cloudfoundry_space.app_space.id
  domain_id = data.cloudfoundry_domain.public.id
  databases = module.databases.ids
  buckets = module.buckets.ids
}
