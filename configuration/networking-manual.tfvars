psc_endpoints = [
  {
    producer_instance_project_id = ""
    endpoint_project_id          = ""
    target                       = "" # in format "projects/xxx-tp/regions/xx-central1/serviceAttachments/gkedpm-xxx"
    subnetwork_name              = "subnetwork-1"
    network_name                 = "network-1"
    ip_address_literal           = "10.128.0.26"
    region                       = "" # Example : us-central1
  },
  {
    producer_instance_project_id = ""
    endpoint_project_id          = ""
    subnetwork_name              = "subnetwork-2"
    network_name                 = "network-2"
    ip_address_literal           = "10.128.0.27"
    region                       = ""                  # Example : us-central2
    producer_instance_name       = "psc-instance-name" # Can only be used for CloudSQL
  }
]