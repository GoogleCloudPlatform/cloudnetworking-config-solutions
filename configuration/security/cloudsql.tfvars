project_id = ""
network    = ""
egress_rules = {
  allow-egress = {
    deny = false
    rules = [{
      protocol = "tcp"
      ports    = ["3306"]
    }]
  }
}
