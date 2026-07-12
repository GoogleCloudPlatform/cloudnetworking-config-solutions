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

module "vm" {
  for_each = local.instance_map
  source   = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/compute-vm?ref=v38.2.0"

  # Basic VM Configuration
  name             = each.value.name
  project_id       = each.value.project_id
  zone             = each.value.zone
  can_ip_forward   = each.value.can_ip_forward
  hostname         = each.value.hostname
  enable_display   = each.value.enable_display
  description      = each.value.description
  instance_type    = each.value.instance_type
  min_cpu_platform = each.value.min_cpu_platform
  tags             = each.value.tags
  labels           = each.value.labels

  # Boot Disk Configuration
  boot_disk = each.value.boot_disk

  # Network Configuration
  network_interfaces = [{
    network    = each.value.network
    subnetwork = each.value.subnetwork
  }]

  # Attached Disks
  attached_disks = [
    for disk in each.value.attached_disks :
    {
      name              = disk.name
      source            = disk.source
      device_name       = try(disk.device_name, null)
      auto_delete       = try(disk.auto_delete, var.attached_disk_defaults.auto_delete)
      mode              = try(disk.options.mode, var.attached_disk_defaults.mode)
      replica_zone      = try(disk.options.replica_zone, var.attached_disk_defaults.replica_zone)
      type              = try(disk.options.type, var.attached_disk_defaults.type)
      snapshot_schedule = disk.snapshot_schedule
    }
    if disk.source_type != "image"
  ]

  # Additional Options
  metadata                    = each.value.metadata
  network_attached_interfaces = each.value.network_attached_interfaces
  options                     = each.value.options
  scratch_disks               = each.value.scratch_disks
  shielded_config             = each.value.shielded_config
  snapshot_schedules          = each.value.snapshot_schedules
  tag_bindings                = each.value.tag_bindings
  service_account             = each.value.service_account
}