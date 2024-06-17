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
      project_id                    = instance.project_id
      name                          = instance.name
      region                        = instance.region
      network_config                = instance.network_config
      database_version              = try(instance.database_version, var.database_version)
      tier                          = try(instance.tier, var.tier)
      availability_type             = try(instance.availability_type, var.availability_type)
      activation_policy             = try(instance.activation_policy, var.activation_policy)
      backup_configuration          = try(instance.backup_configuration, var.backup_configuration)
      collation                     = try(instance.collation, var.collation)
      connector_enforcement         = try(instance.connector_enforcement, var.connector_enforcement)
      data_cache                    = try(instance.data_cache, var.data_cache)
      databases                     = try(instance.databases, var.databases)
      disk_autoresize_limit         = try(instance.disk_autoresize_limit, var.disk_autoresize_limit)
      disk_size                     = try(instance.disk_size, var.disk_size)
      disk_type                     = try(instance.disk_type, var.disk_type)
      edition                       = try(instance.edition, var.edition)
      encryption_key_name           = try(instance.encryption, var.encryption_key_name)
      flags                         = try(instance.flags, var.flags)
      gcp_deletion_protection       = try(instance.gcp_deletion_protection, var.gcp_deletion_protection)
      insights_config               = try(instance.insights_config, var.insights_config)
      labels                        = try(instance.labels, var.labels)
      maintenance_config            = try(instance.maintenance_config, var.maintenance_config)
      prefix                        = try(instance.prefix, var.prefix)
      replicas                      = try(instance.replicas, var.replicas)
      root_password                 = try(instance.root_password, var.root_password)
      ssl                           = try(instance.ssl, var.ssl)
      terraform_deletion_protection = try(instance.terraform_deletion_protection, var.terraform_deletion_protection)
      time_zone                     = try(instance.timezone, var.time_zone)
      users                         = try(instance.users, var.users)
    }
  ])
}
