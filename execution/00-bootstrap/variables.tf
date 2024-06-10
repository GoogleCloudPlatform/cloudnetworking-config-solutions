
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

variable "bootstrap_project_id" {
  type        = string
  description = "Google Cloud Project ID which will be used to create the service account and Google Cloud storage buckets."
}
variable "network_hostproject_id" {
  type        = string
  description = "Google Cloud Project ID for the networking host project to be used to create networking and security resources."
}
variable "network_serviceproject_id" {
  type        = string
  description = "Google Cloud Project ID to be used to create Google Cloud resources like consumer and producer services."
}
variable "gcs_bucket_name" {
  type        = string
  description = "Name of the Google Cloud storage bucket."
  default     = "terraform-state"
}
variable "versioning" {
  type        = bool
  description = "The Goocle Cloud storage bucket versioning."
  default     = true
}
variable "gcs_bucket_location" {
  description = "Location of the Google Cloud storage bucket."
  type        = string
  default     = "EU"
}
variable "organization_sa_name" {
  type        = string
  description = "Name of the service account to create for organization stage."
  default     = "organization-stage-sa"
}
variable "organization_stage_administrator" {
  type        = list(string)
  description = "List of Members to be granted an IAM role. e.g. (group:my-group@example.com),(user:my-user@example.com)"
}
variable "networking_sa_name" {
  type        = string
  description = "Name of the service account to create for networking stage."
  default     = "networking-stage-sa"
}
variable "networking_stage_administrator" {
  type        = list(string)
  description = "List of Members to be granted an IAM role. e.g. (group:my-group@example.com),(user:my-user@example.com)"
}
variable "security_sa_name" {
  type        = string
  description = "Name of the service account to create for security stage."
  default     = "security-stage-sa"
}
variable "security_stage_administrator" {
  type        = list(string)
  description = "List of Members to be granted an IAM role. e.g. (group:my-group@example.com),(user:my-user@example.com)"
}
variable "producer_sa_name" {
  type        = string
  description = "Name of the service account to create for producer stage."
  default     = "producer-stage-sa"
}
variable "producer_stage_administrator" {
  type        = list(string)
  description = "List of Members to be granted an IAM role. e.g. (group:my-group@example.com),(user:my-user@example.com)"
}
variable "networking_manual_sa_name" {
  type        = string
  description = "Name of the service account to create for networking manual stage."
  default     = "networking-manual-stage-sa"
}
variable "networking_manual_stage_administrator" {
  type        = list(string)
  description = "List of Members to be granted an IAM role. e.g. (group:my-group@example.com),(user:my-user@example.com)"
}
variable "consumer_sa_name" {
  type        = string
  description = "Name of the service account to create for consumer stage."
  default     = "consumer-stage-sa"
}
variable "consumer_stage_administrator" {
  type        = list(string)
  description = "List of Members to be granted an IAM role. e.g. (group:my-group@example.com),(user:my-user@example.com)"
}

