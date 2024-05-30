
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

output "storage_bucket_name" {
  description = "Google Cloud storage bucket name."
  value       = module.google_storage_bucket.name
}
output "organization_stage_email" {
  description = "Organization stage service account iam email."
  value       = module.organization.iam_email
}
output "networking_stage_email" {
  description = "Networking stage service account iam email."
  value       = module.networking.iam_email
}
output "security_stage_email" {
  description = "Security stage service account iam email."
  value       = module.security.iam_email
}
output "producer_stage_email" {
  description = "Producer stage service account iam email."
  value       = module.producer.iam_email
}
output "networking_manual_stage_email" {
  description = "Networking manual stage service account iam email."
  value       = module.networking_manual.iam_email
}
output "consumer_stage_email" {
  description = "Consumer stage service account iam email."
  value       = module.consumer.iam_email
}

