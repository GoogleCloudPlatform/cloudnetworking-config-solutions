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
  config_folder_path = var.config_folder_path
  clusters = [
    for file in fileset(local.config_folder_path, "[^_]*.yaml") : yamldecode(file("${local.config_folder_path}/${file}"))
  ]

  cluster_list = flatten([
    for cluster in try(local.clusters, []) : [
      {
        project_id                                  = cluster.project_id
        name                                        = cluster.name
        region                                      = try(cluster.region, var.region)
        zones                                       = try(cluster.zones, var.zones)
        network                                     = cluster.network
        subnetwork                                  = cluster.subnetwork
        description                                 = try(cluster.description, var.description)
        regional                                    = try(cluster.regional, var.regional)
        network_project_id                          = try(cluster.network_project_id, var.network_project_id)
        kubernetes_version                          = try(cluster.kubernetes_version, var.kubernetes_version)
        master_authorized_networks                  = try(cluster.master_authorized_networks, var.master_authorized_networks)
        enable_vertical_pod_autoscaling             = try(cluster.enable_vertical_pod_autoscaling, var.enable_vertical_pod_autoscaling)
        horizontal_pod_autoscaling                  = try(cluster.horizontal_pod_autoscaling, var.horizontal_pod_autoscaling)
        http_load_balancing                         = try(cluster.http_load_balancing, var.http_load_balancing)
        service_external_ips                        = try(cluster.service_external_ips, var.service_external_ips)
        datapath_provider                           = try(cluster.datapath_provider, var.datapath_provider)
        maintenance_start_time                      = try(cluster.maintenance_start_time, var.maintenance_start_time)
        maintenance_exclusions                      = try(cluster.maintenance_exclusions, var.maintenance_exclusions)
        maintenance_end_time                        = try(cluster.maintenance_end_time, var.maintenance_end_time)
        maintenance_recurrence                      = try(cluster.maintenance_recurrence, var.maintenance_recurrence)
        ip_range_pods                               = cluster.ip_range_pods
        additional_ip_range_pods                    = try(cluster.additional_ip_range_pods, var.additional_ip_range_pods)
        ip_range_services                           = cluster.ip_range_services
        stack_type                                  = try(cluster.stack_type, var.stack_type)
        node_pools                                  = try(cluster.node_pools, var.node_pools)
        windows_node_pools                          = try(cluster.windows_node_pools, var.windows_node_pools)
        node_pools_labels                           = try(cluster.node_pools_labels, var.node_pools_labels)
        node_pools_resource_labels                  = try(cluster.node_pools_resource_labels, var.node_pools_resource_labels)
        node_pools_metadata                         = try(cluster.node_pools_metadata, var.node_pools_metadata)
        node_pools_linux_node_configs_sysctls       = try(cluster.node_pools_linux_node_configs_sysctls, var.node_pools_linux_node_configs_sysctls)
        enable_cost_allocation                      = try(cluster.enable_cost_allocation, var.enable_cost_allocation)
        resource_usage_export_dataset_id            = try(cluster.resource_usage_export_dataset_id, var.resource_usage_export_dataset_id)
        enable_network_egress_export                = try(cluster.enable_network_egress_export, var.enable_network_egress_export)
        enable_resource_consumption_export          = try(cluster.enable_resource_consumption_export, var.enable_resource_consumption_export)
        cluster_autoscaling                         = try(cluster.cluster_autoscaling, var.cluster_autoscaling)
        node_pools_taints                           = try(cluster.node_pools_taints, var.node_pools_taints)
        node_pools_tags                             = try(cluster.node_pools_tags, var.node_pools_tags)
        node_pools_oauth_scopes                     = try(cluster.node_pools_oauth_scopes, var.node_pools_oauth_scopes)
        network_tags                                = try(cluster.network_tags, var.network_tags)
        stub_domains                                = try(cluster.stub_domains, var.stub_domains)
        upstream_nameservers                        = try(cluster.upstream_nameservers, var.upstream_nameservers)
        non_masquerade_cidrs                        = try(cluster.non_masquerade_cidrs, var.non_masquerade_cidrs)
        ip_masq_resync_interval                     = try(cluster.ip_masq_resync_interval, var.ip_masq_resync_interval)
        ip_masq_link_local                          = try(cluster.ip_masq_link_local, var.ip_masq_link_local)
        configure_ip_masq                           = try(cluster.configure_ip_masq, var.configure_ip_masq)
        logging_service                             = try(cluster.logging_service, var.logging_service)
        monitoring_service                          = try(cluster.monitoring_service, var.monitoring_service)
        create_service_account                      = try(cluster.create_service_account, var.create_service_account)
        grant_registry_access                       = try(cluster.grant_registry_access, var.grant_registry_access)
        registry_project_ids                        = try(cluster.registry_project_ids, var.registry_project_ids)
        service_account                             = try(cluster.service_account, var.service_account)
        service_account_name                        = try(cluster.service_account_name, var.service_account_name)
        boot_disk_kms_key                           = try(cluster.boot_disk_kms_key, var.boot_disk_kms_key)
        issue_client_certificate                    = try(cluster.issue_client_certificate, var.issue_client_certificate)
        cluster_ipv4_cidr                           = try(cluster.cluster_ipv4_cidr, var.cluster_ipv4_cidr)
        cluster_resource_labels                     = try(cluster.cluster_resource_labels, var.cluster_resource_labels)
        dns_cache                                   = try(cluster.dns_cache, var.dns_cache)
        authenticator_security_group                = try(cluster.authenticator_security_group, var.authenticator_security_group)
        identity_namespace                          = try(cluster.identity_namespace, var.identity_namespace)
        enable_mesh_certificates                    = try(cluster.enable_mesh_certificates, var.enable_mesh_certificates)
        release_channel                             = try(cluster.release_channel, var.release_channel)
        gateway_api_channel                         = try(cluster.gateway_api_channel, var.gateway_api_channel)
        add_cluster_firewall_rules                  = try(cluster.add_cluster_firewall_rules, var.add_cluster_firewall_rules)
        add_master_webhook_firewall_rules           = try(cluster.add_master_webhook_firewall_rules, var.add_master_webhook_firewall_rules)
        firewall_priority                           = try(cluster.firewall_priority, var.firewall_priority)
        firewall_inbound_ports                      = try(cluster.firewall_inbound_ports, var.firewall_inbound_ports)
        add_shadow_firewall_rules                   = try(cluster.add_shadow_firewall_rules, var.add_shadow_firewall_rules)
        shadow_firewall_rules_priority              = try(cluster.shadow_firewall_rules_priority, var.shadow_firewall_rules_priority)
        shadow_firewall_rules_log_config            = try(cluster.shadow_firewall_rules_log_config, var.shadow_firewall_rules_log_config)
        enable_confidential_nodes                   = try(cluster.enable_confidential_nodes, var.enable_confidential_nodes)
        enable_cilium_clusterwide_network_policy    = try(cluster.enable_cilium_clusterwide_network_policy, var.enable_cilium_clusterwide_network_policy)
        security_posture_mode                       = try(cluster.security_posture_mode, var.security_posture_mode)
        security_posture_vulnerability_mode         = try(cluster.security_posture_vulnerability_mode, var.security_posture_vulnerability_mode)
        disable_default_snat                        = try(cluster.disable_default_snat, var.disable_default_snat)
        notification_config_topic                   = try(cluster.notification_config_topic, var.notification_config_topic)
        notification_filter_event_type              = try(cluster.notification_filter_event_type, var.notification_filter_event_type)
        deletion_protection                         = try(cluster.deletion_protection, var.deletion_protection)
        enable_tpu                                  = try(cluster.enable_tpu, var.enable_tpu)
        network_policy                              = try(cluster.network_policy, var.network_policy)
        network_policy_provider                     = try(cluster.network_policy_provider, var.network_policy_provider)
        initial_node_count                          = try(cluster.initial_node_count, var.initial_node_count)
        remove_default_node_pool                    = try(cluster.remove_default_node_pool, var.remove_default_node_pool)
        filestore_csi_driver                        = try(cluster.filestore_csi_driver, var.filestore_csi_driver)
        disable_legacy_metadata_endpoints           = try(cluster.disable_legacy_metadata_endpoints, var.disable_legacy_metadata_endpoints)
        default_max_pods_per_node                   = try(cluster.default_max_pods_per_node, var.default_max_pods_per_node)
        database_encryption                         = try(cluster.database_encryption, var.database_encryption)
        enable_shielded_nodes                       = try(cluster.enable_shielded_nodes, var.enable_shielded_nodes)
        enable_binary_authorization                 = try(cluster.enable_binary_authorization, var.enable_binary_authorization)
        node_metadata                               = try(cluster.node_metadata, var.node_metadata)
        cluster_dns_provider                        = try(cluster.cluster_dns_provider, var.cluster_dns_provider)
        cluster_dns_scope                           = try(cluster.cluster_dns_scope, var.cluster_dns_scope)
        cluster_dns_domain                          = try(cluster.cluster_dns_domain, var.cluster_dns_domain)
        gce_pd_csi_driver                           = try(cluster.gce_pd_csi_driver, var.gce_pd_csi_driver)
        gke_backup_agent_config                     = try(cluster.gke_backup_agent_config, var.gke_backup_agent_config)
        gcs_fuse_csi_driver                         = try(cluster.gcs_fuse_csi_driver, var.gcs_fuse_csi_driver)
        stateful_ha                                 = try(cluster.stateful_ha, var.stateful_ha)
        timeouts                                    = try(cluster.timeouts, var.timeouts)
        monitoring_enable_managed_prometheus        = try(cluster.monitoring_enable_managed_prometheus, var.monitoring_enable_managed_prometheus)
        monitoring_enable_observability_metrics     = try(cluster.monitoring_enable_observability_metrics, var.monitoring_enable_observability_metrics)
        monitoring_observability_metrics_relay_mode = try(cluster.monitoring_observability_metrics_relay_mode, var.monitoring_observability_metrics_relay_mode)
        monitoring_enabled_components               = try(cluster.monitoring_enabled_components, var.monitoring_enabled_components)
        logging_enabled_components                  = try(cluster.logging_enabled_components, var.logging_enabled_components)
        enable_kubernetes_alpha                     = try(cluster.enable_kubernetes_alpha, var.enable_kubernetes_alpha)
        config_connector                            = try(cluster.config_connector, var.config_connector)
        enable_intranode_visibility                 = try(cluster.enable_intranode_visibility, var.enable_intranode_visibility)
        enable_l4_ilb_subsetting                    = try(cluster.enable_l4_ilb_subsetting, var.enable_l4_ilb_subsetting)
        fleet_project                               = try(cluster.fleet_project, var.fleet_project)
        enable_private_endpoint                     = try(cluster.enable_private_endpoint, var.enable_private_endpoint)
        enable_private_nodes                        = try(cluster.enable_private_nodes, var.enable_private_nodes)
        master_ipv4_cidr_block                      = try(cluster.master_ipv4_cidr_block, var.master_ipv4_cidr_block)
        http_load_balancing                         = try(cluster.http_load_balancing, var.http_load_balancing)
        network_policy                              = try(cluster.network_policy, var.network_policy)
        horizontal_pod_autoscaling                  = try(cluster.horizontal_pod_autoscaling, var.horizontal_pod_autoscaling)
        filestore_csi_driver                        = try(cluster.filestore_csi_driver, var.filestore_csi_driver)
      }
    ]
  ])

  # Move cluster_map assignment outside of flatten block
  cluster_map = { for cluster in local.cluster_list : cluster.name => cluster }
}
