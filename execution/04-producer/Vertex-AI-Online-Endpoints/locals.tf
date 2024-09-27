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
  endpoints          = [for file in fileset(local.config_folder_path, "*.yaml") : yamldecode(file("${local.config_folder_path}/${file}"))]

  endpoint_list = {
    for endpoint in local.endpoints : endpoint.display_name => {
      project     = endpoint.project
      name        = try(endpoint.name, var.name)
      description = try(endpoint.description, var.description)
      location    = endpoint.location
      region      = try(endpoint.region, var.region)
      labels      = try(endpoint.labels, var.labels)
      network     = endpoint.network
    }
  }

}