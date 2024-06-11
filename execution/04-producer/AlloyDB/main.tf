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

module "alloy_db" {
  source                      = "GoogleCloudPlatform/alloy-db/google"
  version                     = "~> 2.2.0"
  for_each                    = { for alloydb in local.instance_list : alloydb.cluster_display_name => alloydb }
  project_id                  = each.value.project_id
  cluster_id                  = each.value.cluster_id
  cluster_display_name        = each.value.cluster_display_name
  cluster_location            = each.value.region
  network_self_link           = each.value.network_id
  allocated_ip_range          = each.value.allocated_ip_range
  database_version            = each.value.database_version
  cluster_labels              = each.value.cluster_labels
  cluster_initial_user        = each.value.cluster_initial_user
  primary_instance            = each.value.primary_instance
  read_pool_instance          = each.value.read_pool_instance
  automated_backup_policy     = each.value.automated_backup_policy
  cluster_encryption_key_name = each.value.cluster_encryption_key_name
}
