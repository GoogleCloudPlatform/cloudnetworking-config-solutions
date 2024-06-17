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
  config_folder_path   = var.config_folder_path
  instances            = [for file in fileset(local.config_folder_path, "[^_]*.yaml") : yamldecode(file("${local.config_folder_path}/${file}"))]
  instances_self_links = [for instance in module.vm : instance.self_link]

  instance_list = flatten([
    for instance in try(local.instances, []) : {
      project_id                  = instance.project_id
      name                        = instance.name
      region                      = instance.region
      zone                        = instance.zone
      network                     = instance.network
      image                       = try(instance.image, var.image)
      subnetwork                  = instance.subnetwork
      can_ip_forward              = try(instance.can_ip_forward, var.can_ip_forward)
      hostname                    = try(instance.hostname, var.hostname)
      enable_display              = try(instance.enable_display, var.enable_display)
      description                 = try(instance.description, var.description)
      instance_type               = try(instance.instance_type, var.instance_type)
      min_cpu_platform            = try(instance.min_cpu_platform, var.min_cpu_platform)
      tags                        = try(instance.tags, var.tags)
      labels                      = try(instance.labels, var.labels)
      metadata                    = try(instance.metadata, var.metadata)
      network_attached_interfaces = try(instance.network_attached_interfaces, var.network_attached_interfaces)
      options                     = try(instance.options, var.options)
      scratch_disks               = try(instance.scratch_disks, var.scratch_disks)
      shielded_config             = try(instance.shielded_config, var.shielded_config)
      snapshot_schedules          = try(instance.snapshot_schedules, var.snapshot_schedules)
      tag_bindings                = try(instance.tag_bindings, var.tag_bindings)
      tag_bindings_firewall       = try(instance.tag_bindings_firewall, var.tag_bindings_firewall)

      # Service Account
      service_account = {
        auto_create = try(instance.service_account.auto_create, var.service_account.auto_create)
        email       = try(instance.service_account.email, var.service_account.email)
        scopes      = try(instance.service_account.scopes, var.service_account.scopes)
      }

      # Boot Disk Configuration
      boot_disk = {
        auto_delete       = try(instance.boot_disk.auto_delete, var.boot_disk.auto_delete)
        snapshot_schedule = try(instance.boot_disk.snapshot_schedule, var.boot_disk.snapshot_schedule)
        source            = try(instance.boot_disk.source, var.boot_disk.source)
        initialize_params = {
          image = instance.image
          size  = try(instance.boot_disk.initialize_params.size, var.boot_disk.initialize_params.size)
          type  = try(instance.boot_disk.initialize_params.type, var.boot_disk.initialize_params.type)
        }
        use_independent_disk = try(instance.boot_disk.use_independent_disk, var.boot_disk.use_independent_disk)
      }

      # Attached Disks
      attached_disks = try(instance.attached_disks, var.attached_disks)

    }
  ])

  # Move instance_map assignment outside of flatten block
  instance_map = { for instance in local.instance_list : instance.name => instance }
}