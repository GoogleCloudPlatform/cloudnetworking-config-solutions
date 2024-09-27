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
  config_folder_path = var.config_folder_path
  instances          = [for file in fileset(local.config_folder_path, "[^_]*.yaml") : yamldecode(file("${local.config_folder_path}/${file}"))]
  instance_list = flatten([
    for instance in try(local.instances, []) : {
      project_id                     = instance.project_id
      index_display_name             = instance.index_display_name
      region                         = instance.region
      index_endpoint_network         = try(instance.index_endpoint_network, var.index_endpoint_network)
      index_endpoint_display_name    = instance.index_endpoint_display_name
      deployed_index_id              = instance.deployed_index_id
      index_labels                   = try(instance.index_labels, var.index_labels)
      index_description              = try(instance.index_description, var.index_description)
      index_update_method            = try(instance.index_update_method, var.index_update_method)
      approximate_neighbors_count    = try(instance.approximate_neighbors_count)
      shard_size                     = try(instance.shard_size, var.shard_size)
      distance_measure_type          = try(instance.distance_measure_type, var.distance_measure_type)
      index_endpoint_description     = try(instance.index_endpoint_description, var.index_endpoint_description)
      index_endpoint_labels          = try(instance.index_endpoint_labels, var.index_endpoint_labels)
      tree_ah_config                 = try(instance.tree_ah_config, var.tree_ah_config)
      brute_force_config             = try(instance.brute_force_config, var.brute_force_config)
      deployed_display_name          = try(instance.deployed_display_name, var.deployed_display_name)
      reserved_ip_ranges             = try(instance.reserved_ip_ranges, var.reserved_ip_ranges)
      enable_access_logging          = try(instance.enable_access_logging, var.enable_access_logging)
      deployment_group               = try(instance.deployment_group, var.deployment_group)
      automatic_resources            = try(instance.automatic_resources, var.automatic_resources)
      dedicated_resources            = try(instance.dedicated_resources, var.dedicated_resources)
      deployed_index_auth_config     = try(instance.deployed_index_auth_config, var.deployed_index_auth_config)
      public_endpoint_enabled        = try(instance.public_endpoint_enabled, var.public_endpoint_enabled)
      private_service_connect_config = try(instance.private_service_connect_config, var.private_service_connect_config)
    }
  ])
}
