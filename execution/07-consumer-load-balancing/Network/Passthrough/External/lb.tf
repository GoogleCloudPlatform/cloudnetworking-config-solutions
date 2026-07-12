# Copyright 2025 Google LLC
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

module "nlb_passthrough_ext" {
  for_each = local.lb_map
  source   = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/net-lb-ext?ref=v39.2.0"

  # Required variables mapped from locals
  project_id = each.value.project_id
  region     = each.value.region
  name       = each.value.name

  # Backend configuration mapped from locals (includes constructed group self_link)
  backends = each.value.backends

  # Optional configurations mapped from locals (YAML/variables)
  description            = each.value.description
  labels                 = each.value.labels
  backend_service_config = each.value.backend_service_config

  # Health Check Configuration
  # If health_check_name is provided in YAML/locals, use the existing HC by name/self-link
  # Otherwise, use the auto-created health check config from locals
  health_check        = each.value.health_check_name
  health_check_config = each.value.health_check_name == null ? each.value.health_check_config : null

  # Forwarding Rules Configuration mapped from locals
  forwarding_rules_config = each.value.forwarding_rules_config
}