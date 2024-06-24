project_id = ""
network    = ""
egress_rules = {
  allow-egress-alloydb = {
    deny = false
    rules = [{
      protocol = "tcp"
      ports    = ["5432"]
    }]
  }
}
