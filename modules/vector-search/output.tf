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

output "index_name" {
  value       = google_vertex_ai_index.index.name
  description = "The resource name of the Index."
}

output "index_id" {
  value       = google_vertex_ai_index.index.id
  description = "An identifier for the index resource."
}
output "deploy_index_name" {
  value       = google_vertex_ai_index_endpoint_deployed_index.basic_deployed_index.name
  description = "The name of the DeployedIndex resource."
}

output "deploy_id" {
  value       = google_vertex_ai_index_endpoint_deployed_index.basic_deployed_index.id
  description = "An identifier for the deployed index resource."
}

output "private_endpoints" {
  value       = google_vertex_ai_index_endpoint_deployed_index.basic_deployed_index.private_endpoints
  description = "Provides paths for users to send requests directly to the deployed index services running on Cloud via private services access."
}

output "deployed_indexes" {
  value       = google_vertex_ai_index.index.deployed_indexes
  description = "The pointers to DeployedIndexes created from this Index. An Index can be only deleted if all its DeployedIndexes had been undeployed first."
}

output "index_endpoint_id" {
  value       = google_vertex_ai_index_endpoint.index_endpoint.id
  description = "An identifier for the index endpoint resource with format."
}

output "index_endpoint_name" {
  value       = google_vertex_ai_index_endpoint.index_endpoint.name
  description = "The resource name of the Index."
}

output "public_endpoint_domain_name" {
  value       = google_vertex_ai_index_endpoint.index_endpoint.public_endpoint_domain_name
  description = "If publicEndpointEnabled is true, this field will be populated with the domain name to use for this index endpoint."
}
