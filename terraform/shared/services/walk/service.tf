resource "cloudfoundry_app" "walk" {
  name                 = "walk"
  space                = var.app_space_id # data.cloudfoundry_space.app_space.id
  buildpacks            = ["https://github.com/cloudfoundry/apt-buildpack", "https://github.com/cloudfoundry/binary-buildpack.git"]
  path                 = "${path.module}/../app.tar.gz"
  source_code_hash     = filesha256("${path.module}/../app.tar.gz")
  disk_quota           = var.disk_quota
  memory               = var.memory
  instances            = var.instances
  strategy             = "rolling"
  timeout              = 200
  health_check_type    = "port"
  health_check_timeout = 180
  health_check_http_endpoint = "/api/heartbeat"

  service_binding {
    service_instance = var.databases.queues
  }

  service_binding {
    service_instance = var.buckets.fetch
  }
}
