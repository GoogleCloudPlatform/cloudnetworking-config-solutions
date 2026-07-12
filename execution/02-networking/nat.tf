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

module "nat" {
  count          = var.create_nat ? 1 : 0
  source         = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/net-cloudnat?ref=v36.2.0"
  project_id     = var.project_id
  region         = var.region
  name           = var.nat_name
  router_network = local.network_id
}

resource "google_compute_route" "default" {
  count            = var.create_nat ? 1 : 0
  name             = local.nat_router_name
  project          = var.project_id
  dest_range       = var.destination_range
  network          = local.network_id
  next_hop_gateway = var.next_hop_gateway
}