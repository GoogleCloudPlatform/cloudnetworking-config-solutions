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

resource "google_network_connectivity_service_connection_policy" "policy" {
  count         = var.create_scp_policy ? 1 : 0 # Conditional creation
  name          = "SCP-${local.network_name}-${var.service_class}"
  location      = var.region
  network       = local.network_id
  service_class = var.service_class                              # Static service class
  description   = "Policy for ${var.service_class} connectivity" # Use the static service class
  project       = var.project_id

  psc_config {
    subnetworks = local.subnet_self_links_for_scp_policy
    limit       = var.scp_connection_limit
  }
  depends_on = [module.vpc_network]
}