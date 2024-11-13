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

  # Decode YAML configuration files and build the list of instances
  instances = [for file in fileset(local.config_folder_path, "[^_]*.yaml") : yamldecode(file("${local.config_folder_path}/${file}"))]

  # Instance list with PSC and network configuration fields
  instance_list = flatten([
    for instance in try(local.instances, []) : {
      cluster_id                  = instance.cluster_id
      cluster_display_name        = instance.cluster_display_name
      project_id                  = instance.project_id
      region                      = instance.region
      network_id                  = instance.network_id
      primary_instance            = instance.primary_instance
      database_version            = try(instance.database_version, var.database_version)
      allocated_ip_range          = try(instance.allocated_ip_range, var.allocated_ip_range)
      cluster_labels              = try(instance.cluster_labels, var.cluster_labels)
      cluster_initial_user        = try(instance.cluster_initial_user, var.cluster_initial_user)
      read_pool_instance          = try(instance.read_pool_instance, var.read_pool_instance)
      automated_backup_policy     = try(instance.automated_backup_policy, var.automated_backup_policy)
      cluster_encryption_key_name = try(instance.cluster_encryption_key_name, var.cluster_encryption_key_name)

      # PSC configuration fields
      psc_enabled                   = try(instance.psc_enabled, var.psc_enabled)
      psc_allowed_consumer_projects = try(instance.psc_allowed_consumer_projects, var.psc_allowed_consumer_projects)
    }
  ])

  # Map for alloydb_network_config keyed by cluster_display_name
  alloydb_network_config = {
    for instance in local.instance_list :
    instance.cluster_display_name => {
      network_self_link = instance.psc_enabled ? null : instance.network_id
      psc_config = instance.psc_enabled ? {
        psc_enabled                   = false,
        psc_allowed_consumer_projects = instance.psc_allowed_consumer_projects
      } : null
    }
  }
}
