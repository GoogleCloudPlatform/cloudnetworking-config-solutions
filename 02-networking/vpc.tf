# Copyright 2024 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

module "vpc_network" {
  source                            = "../modules/net-vpc"
  name                              = var.network_name
  vpc_create                        = var.create_network
  project_id                        = var.project_id
  delete_default_routes_on_create   = var.delete_default_routes_on_create
  firewall_policy_enforcement_order = var.firewall_policy_enforcement_order
  deletion_policy                   = var.deletion_policy
  psa_config = {
    export_routes = var.export_custom_routes
    import_routes = var.import_custom_routes
    ranges = {
      "${var.psa_range_name}" = "${var.psa_range}"
    }
  }
  shared_vpc_host             = var.shared_vpc_host
  shared_vpc_service_projects = var.shared_vpc_service_projects
  subnets                     = var.create_subnetwork ? var.subnets : null
}

data "google_compute_network" "vpc_network" {
  name    = module.vpc_network.name
  project = var.project_id
}
