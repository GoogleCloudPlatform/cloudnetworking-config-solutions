project_id = ""
network    = ""
ingress_rules = [
  {
    name        = "allow-ssh-custom-ranges"
    description = "Allow SSH access from specific networks"
    priority    = 1000
    source_ranges = [
      "", # Source ranges such as "192.168.1.0/24" or "10.0.0.0/8"
    ]
    target_tags = ["ssh-allowed", "https-allowed"]
    allow = [{
      protocol = "tcp"
      ports    = ["22", "443"]
    }]
  }
]
