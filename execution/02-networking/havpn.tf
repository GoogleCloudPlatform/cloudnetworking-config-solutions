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

module "havpn" {
  count         = var.create_havpn ? 1 : 0
  source        = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/net-vpn-ha?ref=v36.2.0"
  project_id    = var.project_id
  region        = var.region
  network       = local.network_id
  name          = var.ha_vpn_gateway1_name
  peer_gateways = var.peer_gateways
  router_config = {
    asn = var.router1_asn
    custom_advertise = {
      all_subnets = var.advertise_all_subnets
      ip_ranges = {
        "${var.psa_range}" = "${var.psa_range_name}"
      }
    }
  }
  tunnels = {
    remote-0 = {
      bgp_peer = {
        address = var.tunnel_1_bgp_peer_ip_address
        asn     = var.tunnel_2_bgp_peer_asn
      }
      bgp_session_range     = var.tunnel_1_router_bgp_session_range
      shared_secret         = var.tunnel_1_shared_secret
      vpn_gateway_interface = var.tunnel_1_gateway_interface
    }
    remote-1 = {
      bgp_peer = {
        address = var.tunnel_2_bgp_peer_ip_address
        asn     = var.tunnel_2_bgp_peer_asn
      }
      bgp_session_range     = var.tunnel_2_router_bgp_session_range
      shared_secret         = var.tunnel_2_shared_secret
      vpn_gateway_interface = var.tunnel_2_gateway_interface
    }
  }
}