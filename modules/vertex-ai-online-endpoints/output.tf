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

output "endpoint_configuration" {
  value = {
    name         = var.name
    display_name = var.display_name
    description  = var.description
    location     = var.location
    region       = var.region
    labels       = var.labels
    network      = var.network
  }
  description = "Configuration details provided for the Vertex AI endpoint."
}

output "id" {
  value       = google_vertex_ai_endpoint.endpoint.id # Access directly from the resource
  description = "Identifier for the resource with format projects/{{project}}/locations/{{location}}/endpoints/{{name}}"
}

output "deployed_models" {
  value       = google_vertex_ai_endpoint.endpoint.deployed_models # Access directly 
  description = "The models deployed in this Endpoint."
}

output "etag" {
  value       = google_vertex_ai_endpoint.endpoint.etag
  description = "Used to perform consistent read-modify-write updates. If not set, a blind 'overwrite' update happens."
}

output "create_time" {
  value       = google_vertex_ai_endpoint.endpoint.create_time
  description = "Timestamp when this Endpoint was created."
}

output "update_time" {
  value       = google_vertex_ai_endpoint.endpoint.update_time
  description = "Timestamp when this Endpoint was last updated."
}

output "model_deployment_monitoring_job" {
  value       = google_vertex_ai_endpoint.endpoint.model_deployment_monitoring_job
  description = "Resource name of the Model Monitoring job associated with this Endpoint if monitoring is enabled by CreateModelDeploymentMonitoringJob. Format: projects/{project}/locations/{location}/modelDeploymentMonitoringJobs/{model_deployment_monitoring_job}"
}

output "terraform_labels" {
  value       = google_vertex_ai_endpoint.endpoint.terraform_labels
  description = "The combination of labels configured directly on the resource and default labels configured on the provider."
}

output "effective_labels" {
  value       = google_vertex_ai_endpoint.endpoint.effective_labels
  description = "All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services."
}