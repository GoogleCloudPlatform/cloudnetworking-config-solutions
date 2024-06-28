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

module "cloudsql" {
  source                        = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/cloudsql-instance?ref=v31.1.0"
  for_each                      = { for cloudsql in local.instance_list : cloudsql.name => cloudsql }
  project_id                    = each.value.project_id
  name                          = each.value.name
  database_version              = each.value.database_version
  region                        = each.value.region
  network_config                = each.value.network_config
  tier                          = each.value.tier
  availability_type             = each.value.availability_type
  activation_policy             = each.value.activation_policy
  backup_configuration          = each.value.backup_configuration
  collation                     = each.value.collation
  connector_enforcement         = each.value.connector_enforcement
  data_cache                    = each.value.data_cache
  databases                     = each.value.databases
  disk_autoresize_limit         = each.value.disk_autoresize_limit
  disk_size                     = each.value.disk_size
  disk_type                     = each.value.disk_type
  edition                       = each.value.edition
  encryption_key_name           = each.value.encryption_key_name
  flags                         = each.value.flags
  gcp_deletion_protection       = each.value.gcp_deletion_protection
  insights_config               = each.value.insights_config
  labels                        = each.value.labels
  maintenance_config            = each.value.maintenance_config
  prefix                        = each.value.prefix
  replicas                      = each.value.replicas
  root_password                 = each.value.root_password
  ssl                           = each.value.ssl
  terraform_deletion_protection = each.value.terraform_deletion_protection
  time_zone                     = each.value.time_zone
  users                         = each.value.users
}
