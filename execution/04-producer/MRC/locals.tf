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

  # Create a map directly for the instance list
  instance_list = {
    for instance in local.instances : instance.redis_cluster_name => {
      project_id                  = instance.project_id
      shard_count                 = try(instance.shard_count, var.shard_count)
      network_id                  = instance.network_id
      region                      = try(instance.region, var.region)
      deletion_protection_enabled = try(instance.deletion_protection_enabled, var.deletion_protection_enabled)
      replica_count               = try(instance.replica_count, var.replica_count)
    }
  }

}
