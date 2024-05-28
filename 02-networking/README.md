## Prerequisites

<ol>
  <li>
      A Google Cloud project specifically set up to accommodate Google Cloud networking should exist.
  </li>
  <li>Following are the list of API's that should be enabled : <br>
      "iam.googleapis.com",<br>
      "compute.googleapis.com",<br>
      "logging.googleapis.com",<br>
      "monitoring.googleapis.com",<br>
      "servicenetworking.googleapis.com",<br>
  </li>
  <li> User or Service Account should have the following IAM permissions in the networking project ID : <br>
      <br> "roles/compute.networkAdmin",
      <br> "roles/compute.xpnAdmin",
      <br> If you want to connect to an HA VPN gateway that resides in a Google Cloud organization or project that you don't own, request the compute.vpnGateways.use permission from the owner.
  </li>
</ol>


<!-- BEGIN_TF_DOCS -->

## Providers

| Name | Version |
|------|---------|
| <a name="provider_google"></a> [google](#provider\_google) | 5.32.0 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_havpn"></a> [havpn](#module\_havpn) | github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/net-vpn-ha | v30.0.0 |
| <a name="module_nat"></a> [nat](#module\_nat) | github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/net-cloudnat | v30.0.0 |
| <a name="module_vpc_network"></a> [vpc\_network](#module\_vpc\_network) | ../modules/net-vpc | n/a |

## Resources

| Name | Type |
|------|------|
| [google_compute_route.default](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_route) | resource |
| [google_network_connectivity_service_connection_policy.policy](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/network_connectivity_service_connection_policy) | resource |
| [google_compute_network.vpc_network](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/compute_network) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_advertise_all_subnets"></a> [advertise\_all\_subnets](#input\_advertise\_all\_subnets) | Set to true if all subnets are required to be advertised. | `bool` | `false` | no |
| <a name="input_create_havpn"></a> [create\_havpn](#input\_create\_havpn) | Set to true to create a Cloud HA VPN. | `string` | `"false"` | no |
| <a name="input_create_nat"></a> [create\_nat](#input\_create\_nat) | Set to true to create a Cloud NAT. | `string` | `"true"` | no |
| <a name="input_create_network"></a> [create\_network](#input\_create\_network) | Variable to determine if a new network should be created or not. | `bool` | `true` | no |
| <a name="input_create_scp_policy"></a> [create\_scp\_policy](#input\_create\_scp\_policy) | Boolean flat to create a service connection policy. Set to true to create a service connection policy | `bool` | `false` | no |
| <a name="input_create_subnetwork"></a> [create\_subnetwork](#input\_create\_subnetwork) | Variable to determine if a new sub network should be created or not. | `bool` | `true` | no |
| <a name="input_delete_default_routes_on_create"></a> [delete\_default\_routes\_on\_create](#input\_delete\_default\_routes\_on\_create) | Set to true to delete the default routes at creation time. | `bool` | `true` | no |
| <a name="input_deletion_policy"></a> [deletion\_policy](#input\_deletion\_policy) | The deletion policy for the service networking connection. Setting to ABANDON allows the resource to be abandoned rather than deleted. This will enable a successful terraform destroy when destroying CloudSQL instances. Use with care as it can lead to dangling resources. | `string` | `""` | no |
| <a name="input_destination_range"></a> [destination\_range](#input\_destination\_range) | The destination range of outgoing packets that this route applies to. Only IPv4 is supported. | `string` | `"0.0.0.0/0"` | no |
| <a name="input_export_custom_routes"></a> [export\_custom\_routes](#input\_export\_custom\_routes) | Whether to export the custom routes to the peer network. | `bool` | `true` | no |
| <a name="input_firewall_policy_enforcement_order"></a> [firewall\_policy\_enforcement\_order](#input\_firewall\_policy\_enforcement\_order) | Order that Firewall Rules and Firewall Policies are evaluated. Can be either 'BEFORE\_CLASSIC\_FIREWALL' or 'AFTER\_CLASSIC\_FIREWALL'. | `string` | `"AFTER_CLASSIC_FIREWALL"` | no |
| <a name="input_ha_vpn_gateway1_name"></a> [ha\_vpn\_gateway1\_name](#input\_ha\_vpn\_gateway1\_name) | VPN Gateway name and prefix used for dependent resources. | `string` | `"vpn1"` | no |
| <a name="input_import_custom_routes"></a> [import\_custom\_routes](#input\_import\_custom\_routes) | Whether to import the custom routes to the peer network. | `bool` | `true` | no |
| <a name="input_nat_name"></a> [nat\_name](#input\_nat\_name) | Name of the Cloud NAT to be created. | `string` | `"internet-gateway"` | no |
| <a name="input_network_name"></a> [network\_name](#input\_network\_name) | Name of the VPC network to be created if var.create\_network is marked as true or Name of the already existing network if var.create\_network is false. | `string` | n/a | yes |
| <a name="input_next_hop_gateway"></a> [next\_hop\_gateway](#input\_next\_hop\_gateway) | URL to a gateway that should handle matching packets. Currently, you can only specify the internet gateway, using a full or partial valid URL. | `string` | `"default-internet-gateway"` | no |
| <a name="input_peer_gateways"></a> [peer\_gateways](#input\_peer\_gateways) | Configuration of the (external or GCP) peer gateway. | <pre>map(object({<br>    external = optional(object({<br>      redundancy_type = string<br>      interfaces      = list(string)<br>      description     = optional(string, "Terraform managed external VPN gateway")<br>    }))<br>    gcp = optional(string)<br>  }))</pre> | `{}` | no |
| <a name="input_project_id"></a> [project\_id](#input\_project\_id) | The project ID of the Google Cloud project where the VPC will be created. | `string` | n/a | yes |
| <a name="input_psa_range"></a> [psa\_range](#input\_psa\_range) | Variable to describe the CIDR range required by the PSA/Service Networking. | `string` | `"10.0.64.0/20"` | no |
| <a name="input_psa_range_name"></a> [psa\_range\_name](#input\_psa\_range\_name) | Variable to describe the name of the CIDR range required by the PSA/Service Networking. | `string` | `"psarange"` | no |
| <a name="input_region"></a> [region](#input\_region) | Name of a Google Cloud region. | `string` | n/a | yes |
| <a name="input_router1_asn"></a> [router1\_asn](#input\_router1\_asn) | ASN number required for the router1. | `number` | `64513` | no |
| <a name="input_scp_connection_limit"></a> [scp\_connection\_limit](#input\_scp\_connection\_limit) | Limit of the total number of connections to be allowed through the policy | `string` | `5` | no |
| <a name="input_service_class"></a> [service\_class](#input\_service\_class) | Allowed service class (static) | `string` | `"gcp-memorystore-redis"` | no |
| <a name="input_shared_vpc_host"></a> [shared\_vpc\_host](#input\_shared\_vpc\_host) | Enable shared VPC for this project. | `bool` | `true` | no |
| <a name="input_shared_vpc_service_projects"></a> [shared\_vpc\_service\_projects](#input\_shared\_vpc\_service\_projects) | Shared VPC service projects to register with this host. | `list(string)` | `[]` | no |
| <a name="input_subnets"></a> [subnets](#input\_subnets) | Subnet configuration. | <pre>list(object({<br>    name                  = string<br>    ip_cidr_range         = string<br>    region                = string<br>    description           = optional(string)<br>    enable_private_access = optional(bool, true)<br>    flow_logs_config = optional(object({<br>      aggregation_interval = optional(string)<br>      filter_expression    = optional(string)<br>      flow_sampling        = optional(number)<br>      metadata             = optional(string)<br>      # only if metadata == "CUSTOM_METADATA"<br>      metadata_fields = optional(list(string))<br>    }))<br>    ipv6 = optional(object({<br>      access_type = optional(string, "INTERNAL")<br>      # this field is marked for internal use in the API documentation<br>      # enable_private_access = optional(string)<br>    }))<br>    secondary_ip_ranges = optional(map(string))<br><br>    iam = optional(map(list(string)), {})<br>    iam_bindings = optional(map(object({<br>      role    = string<br>      members = list(string)<br>      condition = optional(object({<br>        expression  = string<br>        title       = string<br>        description = optional(string)<br>      }))<br>    })), {})<br>    iam_bindings_additive = optional(map(object({<br>      member = string<br>      role   = string<br>      condition = optional(object({<br>        expression  = string<br>        title       = string<br>        description = optional(string)<br>      }))<br>    })), {})<br>  }))</pre> | `[]` | no |
| <a name="input_subnets_for_scp_policy"></a> [subnets\_for\_scp\_policy](#input\_subnets\_for\_scp\_policy) | List of subnet names to apply the SCP policy to. | `list(string)` | <pre>[<br>  "default"<br>]</pre> | no |
| <a name="input_tunnel_1_bgp_peer_asn"></a> [tunnel\_1\_bgp\_peer\_asn](#input\_tunnel\_1\_bgp\_peer\_asn) | Peer BGP Autonomous System Number (ASN). | `number` | n/a | yes |
| <a name="input_tunnel_1_bgp_peer_ip_address"></a> [tunnel\_1\_bgp\_peer\_ip\_address](#input\_tunnel\_1\_bgp\_peer\_ip\_address) | Peer IP address of the BGP interface outside Google Cloud. Only IPv4 is supported. | `string` | n/a | yes |
| <a name="input_tunnel_1_gateway_interface"></a> [tunnel\_1\_gateway\_interface](#input\_tunnel\_1\_gateway\_interface) | The interface ID of the VPN gateway with which this VPN tunnel is associated. | `number` | `0` | no |
| <a name="input_tunnel_1_router_bgp_session_range"></a> [tunnel\_1\_router\_bgp\_session\_range](#input\_tunnel\_1\_router\_bgp\_session\_range) | IP address and range of the interface. | `string` | `"169.254.1.2/30"` | no |
| <a name="input_tunnel_1_shared_secret"></a> [tunnel\_1\_shared\_secret](#input\_tunnel\_1\_shared\_secret) | Shared secret used to set the secure session between the Cloud VPN gateway and the peer VPN gateway. Note: This property is sensitive and should be preserved carefully. | `string` | n/a | yes |
| <a name="input_tunnel_2_bgp_peer_asn"></a> [tunnel\_2\_bgp\_peer\_asn](#input\_tunnel\_2\_bgp\_peer\_asn) | Peer BGP Autonomous System Number (ASN). | `number` | n/a | yes |
| <a name="input_tunnel_2_bgp_peer_ip_address"></a> [tunnel\_2\_bgp\_peer\_ip\_address](#input\_tunnel\_2\_bgp\_peer\_ip\_address) | Peer IP address of the BGP interface outside Google Cloud. Only IPv4 is supported. | `string` | n/a | yes |
| <a name="input_tunnel_2_gateway_interface"></a> [tunnel\_2\_gateway\_interface](#input\_tunnel\_2\_gateway\_interface) | The interface ID of the VPN gateway with which this VPN tunnel is associated. | `number` | `1` | no |
| <a name="input_tunnel_2_router_bgp_session_range"></a> [tunnel\_2\_router\_bgp\_session\_range](#input\_tunnel\_2\_router\_bgp\_session\_range) | IP address and range of the interface. | `string` | `"169.254.2.2/30"` | no |
| <a name="input_tunnel_2_shared_secret"></a> [tunnel\_2\_shared\_secret](#input\_tunnel\_2\_shared\_secret) | Shared secret used to set the secure session between the Cloud VPN gateway and the peer VPN gateway. Note: This property is sensitive and should be preserved carefully. | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_name"></a> [name](#output\_name) | Name of the VPC network. |
| <a name="output_network_id"></a> [network\_id](#output\_network\_id) | Fully qualified network ID. |
| <a name="output_service_connection_policy_details"></a> [service\_connection\_policy\_details](#output\_service\_connection\_policy\_details) | Detailed information about each service connection policy |
| <a name="output_service_connection_policy_ids"></a> [service\_connection\_policy\_ids](#output\_service\_connection\_policy\_ids) | Map of service class to service connection policy IDs |
| <a name="output_subnet_ids"></a> [subnet\_ids](#output\_subnet\_ids) | List of fully qualified subnetwork IDs. |
| <a name="output_subnet_self_links_for_scp_policy"></a> [subnet\_self\_links\_for\_scp\_policy](#output\_subnet\_self\_links\_for\_scp\_policy) | The self-links of the subnets where the SCP policy is applied. |
| <a name="output_vpc_networks"></a> [vpc\_networks](#output\_vpc\_networks) | Complete details of the VPC network. |
<!-- END_TF_DOCS -->