project_id = ""
network    = ""
egress_rules = {
  allow-egress-cloudsql = {
    deny = false
    rules = [{
      protocol = "tcp"
      ports    = ["3306"]
    }]
  }
}
