variable disk_quota {
  type        = number
  description = "app disk allocation (in MB)"
  default     = 256
}

variable memory {
  type        = number
  description = "app ram allocation (in MB)"
  default     = 128
}

variable instances {
  type        = number
  description = "number of app instances"
  default     = 1
}

variable gin_debug_level {
  type = string
  default = "release"
}

variable gitref {
  type        = string
  description = "gitref for the specific version of app that you want to use"
  default     = "refs/heads/main"
  # You can also specify a specific commit, eg "7487f882903b9e834a5133a883a88b16fb8b16c9"
}


variable app_space_id {
  type = string
}

variable domain_id {
  type = string
}

variable db_queue_id {
  type = string
}

variable db_work_id {
  type = string
}

variable db_search_id {
  type = string
}

variable space_name {
  type = string
}