project_id = ""
network    = ""
egress_rules = {
  allow-egress-mrc = {
    deny = false
    rules = [{
      protocol = "tcp"
      ports    = ["6379"]
    }]
  }
}
