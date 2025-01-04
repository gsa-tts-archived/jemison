#################################################################
# POSTGRES
#################################################################

module db_queues {
  source = "github.com/gsa-tts/terraform-cloudgov//database?ref=v0.9.1"
  cf_org_name      = var.cf_org
  cf_space_name    = var.cf_space
  name             = "queuesdb"
  tags             = ["rds"]
  rds_plan_name    = var.queue_db_plan_name
}


module db_search {
  source = "github.com/gsa-tts/terraform-cloudgov//database?ref=v0.9.1"
  cf_org_name      = var.cf_org
  cf_space_name    = var.cf_space
  name             = "searchdb"
  tags             = ["rds"]
  rds_plan_name    = var.search_db_plan_name
}


module db_work {
  source = "github.com/gsa-tts/terraform-cloudgov//database?ref=v0.9.1"
  cf_org_name      = var.cf_org
  cf_space_name    = var.cf_space
  name             = "workdb"
  tags             = ["rds"]
  rds_plan_name    = var.work_db_plan_name
}
