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

output "name" {
  description = "Name of the VPC network."
  value       = local.network_name
}

output "network_id" {
  description = "Fully qualified network ID."
  value       = local.network_id
}

output "subnet_ids" {
  description = "List of fully qualified subnetwork IDs."
  value       = local.subnet_ids
}

output "vpc_networks" {
  description = "Complete details of the VPC network."
  value       = module.vpc_network
}
