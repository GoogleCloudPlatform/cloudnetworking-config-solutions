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

module "vector_search" {
  source                         = "../../../modules/vector-search"
  for_each                       = { for instances in local.instance_list : instances.index_display_name => instances }
  project_id                     = each.value.project_id
  region                         = each.value.region
  index_labels                   = each.value.index_labels
  index_display_name             = each.value.index_display_name
  deployed_index_id              = each.value.deployed_index_id
  index_description              = each.value.index_description
  index_update_method            = each.value.index_update_method
  index_endpoint_display_name    = each.value.index_endpoint_display_name
  index_endpoint_description     = each.value.index_endpoint_description
  index_endpoint_labels          = each.value.index_endpoint_labels
  index_endpoint_network         = each.value.index_endpoint_network
  tree_ah_config                 = each.value.tree_ah_config
  brute_force_config             = each.value.brute_force_config
  deployed_display_name          = each.value.deployed_display_name
  reserved_ip_ranges             = each.value.reserved_ip_ranges
  enable_access_logging          = each.value.enable_access_logging
  deployment_group               = each.value.deployment_group
  automatic_resources            = each.value.automatic_resources
  dedicated_resources            = each.value.dedicated_resources
  deployed_index_auth_config     = each.value.deployed_index_auth_config
  public_endpoint_enabled        = each.value.public_endpoint_enabled
  private_service_connect_config = each.value.private_service_connect_config
}
