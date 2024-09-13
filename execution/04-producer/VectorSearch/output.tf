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

output "vector_search_instance_details" {
  description = "Display Vector search instance attributes including project ID, region, index name, index endpoint name."
  value = { for name, instance in module.vector_search :
    name => {
      "index_id" : instance.index_id,
      "index_endpoint_id" : instance.index_endpoint_id,
      "index_name" : instance.index_name,
      "index_endpoint_name" : instance.index_endpoint_name,
      "deployed_indexes" : instance.deployed_indexes,
      "deploy_index_name" : instance.deploy_index_name,
      "deploy_id" : instance.deploy_id,
      "private_endpoints" : instance.private_endpoints,
  } }
}
