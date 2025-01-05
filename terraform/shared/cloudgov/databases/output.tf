output ids {
  value = {
    queues = module.db_queues.instance_id
    search = module.db_search.instance_id
    work = module.db_work.instance_id
  }
}