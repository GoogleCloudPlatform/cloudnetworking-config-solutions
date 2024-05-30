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

output "instances_self_links" {
  description = "List of self-links for compute instances"
  value       = [for instance in module.vm : instance.self_link] # Correctly access the module output
}


output "external_ips" {
  description = "Instance external IP addresses."
  value = { for instance in module.vm :
    instance.id => instance.instance.network_interface[0].access_config[0].nat_ip
    if length(instance.instance.network_interface[0].access_config) > 0
  }
  sensitive = true
}

output "id" {
  description = "Fully qualified instance id."
  value       = [for instance in module.vm : instance.instance.id]
  sensitive   = true
}

output "internal_ips" {
  description = "Instance internal IP addresses."
  value = { for instance in module.vm :
    instance.id => instance.instance.network_interface[0].network_ip
  }
  sensitive = true
}

output "vm_instances" {
  description = "Map of VM instance information"
  value = {
    for instance in local.instances_self_links :
    replace(instance, "https://www.googleapis.com/compute/v1/projects/", "") => {
      name       = element(split("/", instance), length(split("/", instance)) - 1)
      self_link  = instance
      zone       = element(split("/zones/", instance), 1)
      image      = local.instance_map[element(split("/", instance), length(split("/", instance)) - 1)].image
      subnetwork = local.instance_map[element(split("/", instance), length(split("/", instance)) - 1)].subnetwork
      network    = local.instance_map[element(split("/", instance), length(split("/", instance)) - 1)].network
    }
  }
}
