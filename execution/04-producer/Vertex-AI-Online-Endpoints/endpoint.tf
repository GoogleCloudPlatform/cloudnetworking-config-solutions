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

module "vertex_endpoints" {
  for_each = local.endpoint_list

  source             = "../../../modules/vertex-ai-online-endpoints/"
  config_folder_path = var.config_folder_path
  display_name       = each.key

  name    = each.value.name
  project = each.value.project

  description = each.value.description
  location    = each.value.location

  region  = each.value.region
  labels  = each.value.labels
  network = each.value.network
}