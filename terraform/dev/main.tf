
data "cloudfoundry_domain" "public" {
  name = "app.cloud.gov"
}

data "cloudfoundry_space" "app_space" {
  # name = var.cf_space_guid
  name = local.cf_space
  org  = var.cf_org_guid
}

resource "cloudfoundry_route" "admin_route" {
  space    = data.cloudfoundry_space.app_space.id
  domain   = data.cloudfoundry_domain.public.id
  hostname = "jemison-admin"
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
  admin_route_id = cloudfoundry_route.admin_route.id
  app_space_id = data.cloudfoundry_space.app_space.id
  db_queue_id = module.databases.db_queue_id
  db_search_id = module.databases.db_search_id
  db_work_id = module.databases.db_work_id
  app_zip = "blargh"
} 