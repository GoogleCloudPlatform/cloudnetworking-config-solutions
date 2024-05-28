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

locals {
  network_name    = try(module.vpc_network.name, "")
  network_id      = try(module.vpc_network.id, "")
  nat_router_name = "${var.nat_name}-route"
  // For generating list of subnet IDs to be displayed as output.
  subnet_ids = [for subnetwork_link in data.google_compute_network.vpc_network.subnetworks_self_links : trimprefix(subnetwork_link, "https://www.googleapis.com/compute/v1/")]
  // Subnets for SCP
  subnet_self_links_for_scp_policy = [
    for subnet in module.vpc_network.subnets :
    subnet.self_link
    if contains(var.subnets_for_scp_policy, subnet.name)
  ]
}
