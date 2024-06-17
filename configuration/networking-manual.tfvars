psc_endpoints = [
  {
    endpoint_project_id          = "project-for-endpoint"
    producer_instance_project_id = "project-producer"
    producer_instance_name       = "sql-1"
    subnetwork_name              = "subnetwork-1"
    network_name                 = "network-1"
    ip_address_literal           = "10.128.0.50"
  },
  {
    endpoint_project_id          = "project-for-endpoint"
    producer_instance_project_id = "project-producer"
    producer_instance_name       = "sql-2"
    subnetwork_name              = "subnetwork-2"
    network_name                 = "network-2"
    ip_address_literal           = ""
  },
  {
    endpoint_project_id          = "project-for-endpoint"
    producer_instance_project_id = "project-producer"
    producer_instance_name       = "sql-3"
    subnetwork_name              = "subnetwork-3"
    network_name                 = "network-3"
    ip_address_literal           = ""
  }
]
