/**
 * Copyright 2025 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

module "network_firewall_policy" {
  for_each      = local.instance_map
  source        = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/net-firewall-policy?ref=v40.2.0"
  region        = each.value.region
  name          = each.value.name
  parent_id     = each.value.parent_id
  attachments   = each.value.attachments
  description   = each.value.description
  egress_rules  = each.value.egress_rules
  ingress_rules = each.value.ingress_rules
}
