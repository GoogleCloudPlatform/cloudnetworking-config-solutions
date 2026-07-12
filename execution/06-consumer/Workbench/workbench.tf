# Copyright 2025 Google LLC
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


module "workbench_instance" {
  for_each = local.workbench_map
  source   = "GoogleCloudPlatform/vertex-ai/google//modules/workbench"
  version  = "~> 0.2"

  name       = each.value.name
  project_id = each.value.project_id
  location   = each.value.location # Module internally maps this to 'zone'

  machine_type     = try(each.value.gce_setup.machine_type, var.machine_type)
  service_accounts = try(each.value.gce_setup.service_accounts, var.service_accounts)
  metadata         = merge(try(each.value.gce_setup.metadata, {}), try(var.metadata_configs, {}), var.metadata)
  instance_owners  = try(each.value.gce_setup.instance_owners, var.instance_owners)
  labels           = try(each.value.gce_setup.labels, var.labels)

  vm_image = {
    project = try(each.value.gce_setup.vm_image.project, var.vm_image_project)
    family  = try(each.value.gce_setup.vm_image.family, var.vm_image_family)
    name    = try(each.value.gce_setup.vm_image.name, var.vm_image_name)
  }

  boot_disk_type    = try(each.value.gce_setup.boot_disk_type, var.boot_disk_type)
  boot_disk_size_gb = try(each.value.gce_setup.boot_disk_size_gb, var.boot_disk_size_gb)

  data_disks = [
    {
      disk_size_gb    = try(each.value.gce_setup.data_disks[0].disk_size_gb, var.data_disk_size)
      disk_type       = try(each.value.gce_setup.data_disks[0].disk_type, var.data_disk_type)
      disk_encryption = try(each.value.gce_setup.data_disks[0].disk_encryption, var.disk_encryption_default)
    }
  ]

  network_interfaces = [
    {
      network          = each.value.gce_setup.network_interfaces[0].network
      subnet           = each.value.gce_setup.network_interfaces[0].subnet
      nic_type         = try(each.value.gce_setup.network_interfaces[0].nic_type, var.nic_type)
      internal_ip_only = try(each.value.gce_setup.network_interfaces[0].internal_ip_only, var.internal_ip_only)
    }
  ]

  disable_public_ip    = try(each.value.gce_setup.disable_public_ip, var.disable_public_ip_default)
  disable_proxy_access = try(each.value.gce_setup.disable_proxy_access, var.disable_proxy_access_default)
  tags                 = try(each.value.gce_setup.tags, var.network_tags)

  shielded_instance_config = {
    enable_secure_boot          = try(each.value.gce_setup.enable_secure_boot, var.enable_secure_boot_default)
    enable_vtpm                 = try(each.value.gce_setup.enable_vtpm, var.enable_vtpm_default)
    enable_integrity_monitoring = try(each.value.gce_setup.enable_integrity_monitoring, var.enable_integrity_monitoring_default)
  }
}