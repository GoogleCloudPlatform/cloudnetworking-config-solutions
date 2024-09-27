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

variable "project_id" {
  type        = string
  description = "The project ID of the Google Cloud project where the VPC will be created."
}

variable "network_name" {
  type        = string
  description = "Name of the VPC network to be created if var.create_network is marked as true or Name of the already existing network if var.create_network is false."
}

variable "create_network" {
  type        = bool
  default     = true
  description = "Variable to determine if a new network should be created or not."
}

variable "create_subnetwork" {
  type        = bool
  default     = true
  description = "Variable to determine if a new sub network should be created or not."
}

variable "psa_range_name" {
  type        = string
  default     = "psarange"
  description = "Variable to describe the name of the CIDR range required by the PSA/Service Networking."
}

variable "psa_range" {
  type        = string
  default     = "10.0.64.0/20"
  description = "Variable to describe the CIDR range required by the PSA/Service Networking."
}

variable "subnets" {
  description = "Subnet configuration."
  type = list(object({
    name                  = string
    ip_cidr_range         = string
    region                = string
    description           = optional(string)
    enable_private_access = optional(bool, true)
    flow_logs_config = optional(object({
      aggregation_interval = optional(string)
      filter_expression    = optional(string)
      flow_sampling        = optional(number)
      metadata             = optional(string)
      # only if metadata == "CUSTOM_METADATA"
      metadata_fields = optional(list(string))
    }))
    ipv6 = optional(object({
      access_type = optional(string, "INTERNAL")
      # this field is marked for internal use in the API documentation
      # enable_private_access = optional(string)
    }))
    secondary_ip_ranges = optional(map(string))

    iam = optional(map(list(string)), {})
    iam_bindings = optional(map(object({
      role    = string
      members = list(string)
      condition = optional(object({
        expression  = string
        title       = string
        description = optional(string)
      }))
    })), {})
    iam_bindings_additive = optional(map(object({
      member = string
      role   = string
      condition = optional(object({
        expression  = string
        title       = string
        description = optional(string)
      }))
    })), {})
  }))
  default  = []
  nullable = false
}

variable "region" {
  type        = string
  description = "Name of a Google Cloud region."
}

variable "delete_default_routes_on_create" {
  description = "Set to true to delete the default routes at creation time."
  type        = bool
  default     = true
}

variable "firewall_policy_enforcement_order" {
  description = "Order that Firewall Rules and Firewall Policies are evaluated. Can be either 'BEFORE_CLASSIC_FIREWALL' or 'AFTER_CLASSIC_FIREWALL'."
  type        = string
  nullable    = false
  default     = "AFTER_CLASSIC_FIREWALL"

  validation {
    condition     = var.firewall_policy_enforcement_order == "BEFORE_CLASSIC_FIREWALL" || var.firewall_policy_enforcement_order == "AFTER_CLASSIC_FIREWALL"
    error_message = "Enforcement order must be BEFORE_CLASSIC_FIREWALL or AFTER_CLASSIC_FIREWALL."
  }
}

variable "shared_vpc_host" {
  description = "Enable shared VPC for this project."
  type        = bool
  default     = true
}

variable "shared_vpc_service_projects" {
  description = "Shared VPC service projects to register with this host."
  type        = list(string)
  default     = []
}

variable "deletion_policy" {
  description = "The deletion policy for the service networking connection. Setting to ABANDON allows the resource to be abandoned rather than deleted. This will enable a successful terraform destroy when destroying CloudSQL instances. Use with care as it can lead to dangling resources."
  type        = string
  default     = ""
}

## Cloud NAT input variables

variable "create_nat" {
  type        = string
  description = "Set to true to create a Cloud NAT."
  default     = "true"
}

variable "nat_name" {
  type        = string
  description = "Name of the Cloud NAT to be created."
  default     = "internet-gateway"
}

variable "destination_range" {
  type        = string
  description = "The destination range of outgoing packets that this route applies to. Only IPv4 is supported."
  default     = "0.0.0.0/0"
}

variable "next_hop_gateway" {
  type        = string
  description = "URL to a gateway that should handle matching packets. Currently, you can only specify the internet gateway, using a full or partial valid URL."
  default     = "default-internet-gateway"
}

## Cloud HA VPN input variables

variable "create_havpn" {
  type        = string
  description = "Set to true to create a Cloud HA VPN."
  default     = "false"
}

variable "ha_vpn_gateway1_name" {
  description = "VPN Gateway name and prefix used for dependent resources."
  type        = string
  default     = "vpn1"
}

variable "peer_gateways" {
  description = "Configuration of the (external or GCP) peer gateway."
  type = map(object({
    external = optional(object({
      redundancy_type = string
      interfaces      = list(string)
      description     = optional(string, "Terraform managed external VPN gateway")
    }))
    gcp = optional(string)
  }))
  nullable = false
  default  = {}
  validation {
    condition = alltrue([
      for k, v in var.peer_gateways : (v.external != null) != (v.gcp != null)
    ])
    error_message = "Peer gateway configuration must define exactly one between `external` and `gcp`."
  }
}

variable "router1_asn" {
  type        = number
  description = "ASN number required for the router1."
  default     = 64513
}


variable "advertise_all_subnets" {
  type        = bool
  description = "Set to true if all subnets are required to be advertised."
  default     = false
}

variable "export_custom_routes" {
  type        = bool
  description = "Whether to export the custom routes to the peer network."
  default     = true
}

variable "import_custom_routes" {
  type        = bool
  description = "Whether to import the custom routes to the peer network."
  default     = true
}

variable "tunnel_1_bgp_peer_ip_address" {
  type        = string
  description = "Peer IP address of the BGP interface outside Google Cloud. Only IPv4 is supported."
}

variable "tunnel_1_bgp_peer_asn" {
  type        = number
  description = "Peer BGP Autonomous System Number (ASN)."
}

variable "tunnel_1_router_bgp_session_range" {
  type        = string
  description = "IP address and range of the interface."
  default     = "169.254.1.2/30"
}

variable "tunnel_1_shared_secret" {
  type        = string
  description = "Shared secret used to set the secure session between the Cloud VPN gateway and the peer VPN gateway. Note: This property is sensitive and should be preserved carefully."
}

variable "tunnel_1_gateway_interface" {
  type        = number
  description = "The interface ID of the VPN gateway with which this VPN tunnel is associated."
  default     = 0
}

variable "tunnel_2_bgp_peer_ip_address" {
  type        = string
  description = "Peer IP address of the BGP interface outside Google Cloud. Only IPv4 is supported."
}

variable "tunnel_2_bgp_peer_asn" {
  type        = number
  description = "Peer BGP Autonomous System Number (ASN)."
}

variable "tunnel_2_router_bgp_session_range" {
  type        = string
  description = "IP address and range of the interface."
  default     = "169.254.2.2/30"
}

variable "tunnel_2_shared_secret" {
  type        = string
  description = "Shared secret used to set the secure session between the Cloud VPN gateway and the peer VPN gateway. Note: This property is sensitive and should be preserved carefully."
}

variable "tunnel_2_gateway_interface" {
  type        = number
  description = "The interface ID of the VPN gateway with which this VPN tunnel is associated."
  default     = 1
}

variable "create_scp_policy" {
  type        = bool
  description = "Boolean flat to create a service connection policy. Set to true to create a service connection policy"
  default     = false
}

variable "subnets_for_scp_policy" {
  type        = list(string)
  description = "List of subnet names to apply the SCP policy to."
  default     = [""]
}

variable "scp_connection_limit" {
  type        = string
  default     = 5
  description = "Limit of the total number of connections to be allowed through the policy"
}

variable "service_class" {
  type        = string
  default     = "gcp-memorystore-redis"
  description = "Allowed service class (static)"
}


##Interconnect

variable "interconnect_project_id" {
  description = "The ID of the project in which the resource(physical connection at colocation facilitity) belongs."
  type        = string
  default     = ""
}

variable "first_interconnect_name" {
  description = "Name of the first interconnect object. This will be used to populate the URL of the underlying Interconnect object that this attachment's traffic will traverse through."
  type        = string
  default     = ""
}

variable "second_interconnect_name" {
  description = "Name of the second interconnect object. This will be used to populate the URL of the underlying Interconnect object that this attachment's traffic will traverse through."
  type        = string
  default     = ""
}

variable "ic_router_name" {
  description = "Name of the interconnect router."
  type        = string
  default     = "interconnect-router"
}

variable "ic_router_bgp_asn" {
  description = "Local BGP Autonomous System Number (ASN). Must be an RFC6996 private ASN, either 16-bit or 32-bit. The value will be fixed for this router resource."
  type        = string
  default     = ""
}

variable "ic_router_advertise_mode" {
  description = "User-specified flag to indicate which mode to use for advertisement. Default value is DEFAULT. Possible values are: DEFAULT, CUSTOM"
  type        = string
  default     = "CUSTOM"
}

variable "ic_router_advertise_groups" {
  description = "User-specified list of prefix groups to advertise in custom mode. This field can only be populated if advertiseMode is CUSTOM and is advertised to all peers of the router."
  type        = list(string)
  default     = ["ALL_SUBNETS"]
}

variable "user_specified_ip_range" {
  description = "User-specified list of individual IP ranges to advertise in custom mode. This range specifies google private api address."
  type        = list(string)
  default     = ["199.36.154.8/30"]
}

## First Vlan Attachment

variable "create_interconnect" {
  type        = string
  description = "Set to true to create google cloud resources for setting up dedicated interconnect."
  default     = "false"
}

variable "first_va_name" {
  description = "The name of the first interconnect attachment"
  type        = string
  default     = "dedicated-ic-vlan-attachment-3"
}

variable "first_va_description" {
  description = "The description of the first interconnect attachment"
  type        = string
  default     = "interconnect-a vlan attachment 0"
}

variable "first_va_asn" {
  description = "(Required) Local BGP Autonomous System Number (ASN). Must be an RFC6996 private ASN, either 16-bit or 32-bit. The value will be fixed for this router resource."
  type        = string
  default     = ""
}

variable "first_va_bandwidth" {
  description = "Provisioned bandwidth capacity for the first interconnect attachment."
  type        = string
  default     = "BPS_1G"
}

variable "first_va_bgp_range" {
  description = "Up to 16 candidate prefixes that can be used to restrict the allocation of cloudRouterIpAddress and customerRouterIpAddress for this attachment. All prefixes must be within link-local address space (169.254.0.0/16) and must be /29 or shorter (/28, /27, etc)."
  type        = string
  default     = ""
}

variable "first_vlan_tag" {
  description = "The IEEE 802.1Q VLAN tag for this attachment, in the range 2-4094."
  type        = number
  default     = null
}

variable "create_first_vc_router" {
  description = "Select 'true' to create a separate router for this VLAN attachment, or 'false' to use the current router configuration."
  type        = bool
  default     = false
}

## Second Vlan Attachment

variable "second_va_name" {
  description = "The name of the Second interconnect attachment."
  type        = string
  default     = "dedicated-ic-vlan-attachment-4"
}

variable "second_va_description" {
  description = "The description of the second interconnect attachment"
  type        = string
  default     = "interconnect-b vlan attachment 1"
}

variable "second_va_asn" {
  description = "(Required) Local BGP Autonomous System Number (ASN). Must be an RFC6996 private ASN, either 16-bit or 32-bit. The value will be fixed for this router resource."
  type        = string
  default     = ""
}

variable "second_va_bandwidth" {
  description = "Provisioned bandwidth capacity for the second interconnect attachment."
  type        = string
  default     = "BPS_1G"
}

variable "second_va_bgp_range" {
  description = "Up to 16 candidate prefixes that can be used to restrict the allocation of cloudRouterIpAddress and customerRouterIpAddress for this attachment. All prefixes must be within link-local address space (169.254.0.0/16) and must be /29 or shorter (/28, /27, etc)."
  type        = string
  default     = ""
}

variable "second_vlan_tag" {
  description = "The IEEE 802.1Q VLAN tag for this attachment, in the range 2-4094."
  type        = number
  default     = null
}

variable "create_second_vc_router" {
  description = "Select 'true' to create a separate router for this VLAN attachment, or 'false' to use the current router configuration."
  type        = bool
  default     = false
}

variable "admin_enabled" {
  description = "Whether the VLAN attachment is enabled."
  type        = bool
  default     = true
}
