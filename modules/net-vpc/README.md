# VPC module

This module allows creation and management of VPC networks including subnetworks and subnetwork IAM bindings, and most features and options related to VPCs and subnets.

## Examples

<!-- BEGIN TOC -->
- [Examples](#examples)
  - [Simple VPC](#simple-vpc)
  - [Subnet Options](#subnet-options)
  - [Subnet IAM](#subnet-iam)
  - [Peering](#peering)
  - [Shared VPC](#shared-vpc)
  - [Private Service Networking](#private-service-networking)
  - [Private Service Networking with peering routes and peered Cloud DNS domains](#private-service-networking-with-peering-routes-and-peered-cloud-dns-domains)
  - [Subnets for Private Service Connect, Proxy-only subnets](#subnets-for-private-service-connect-proxy-only-subnets)
  - [DNS Policies](#dns-policies)
  - [Subnet Factory](#subnet-factory)
  - [Custom Routes](#custom-routes)
  - [Policy Based Routes](#policy-based-routes)
  - [Private Google Access routes](#private-google-access-routes)
  - [Allow Firewall Policy to be evaluated before Firewall Rules](#allow-firewall-policy-to-be-evaluated-before-firewall-rules)
  - [IPv6](#ipv6)
- [Variables](#variables)
- [Outputs](#outputs)
<!-- END TOC -->

### Simple VPC

```hcl
module "vpc" {
  source     = "./fabric/modules/net-vpc"
  project_id = var.project_id
  name       = "my-network"
  subnets = [
    {
      ip_cidr_range = "10.0.0.0/24"
      name          = "production"
      region        = "europe-west1"
      secondary_ip_ranges = {
        pods     = "172.16.0.0/20"
        services = "192.168.0.0/24"
      }
    },
    {
      ip_cidr_range = "10.0.16.0/24"
      name          = "production"
      region        = "europe-west2"
    }
  ]
}
# tftest modules=1 resources=5 inventory=simple.yaml e2e
```

### Subnet Options

```hcl
module "vpc" {
  source     = "./fabric/modules/net-vpc"
  project_id = var.project_id
  name       = "my-network"
  subnets = [
    # simple subnet
    {
      name          = "simple"
      region        = "europe-west1"
      ip_cidr_range = "10.0.0.0/24"
    },
    # custom description and PGA disabled
    {
      name                  = "no-pga"
      region                = "europe-west1"
      ip_cidr_range         = "10.0.1.0/24",
      description           = "Subnet b"
      enable_private_access = false
    },
    # secondary ranges
    {
      name          = "with-secondary-ranges"
      region        = "europe-west1"
      ip_cidr_range = "10.0.2.0/24"
      secondary_ip_ranges = {
        a = "192.168.0.0/24"
        b = "192.168.1.0/24"
      }
    },
    # enable flow logs
    {
      name          = "with-flow-logs"
      region        = "europe-west1"
      ip_cidr_range = "10.0.3.0/24"
      flow_logs_config = {
        flow_sampling        = 0.5
        aggregation_interval = "INTERVAL_10_MIN"
      }
    }
  ]
}
# tftest modules=1 resources=7 inventory=subnet-options.yaml e2e
```

### Subnet IAM

Subnet IAM variables follow our general interface, with extra keys/members for the subnet to which each binding will be applied.

```hcl
module "vpc" {
  source     = "./fabric/modules/net-vpc"
  project_id = var.project_id
  name       = "my-network"
  subnets = [
    {
      name          = "subnet-1"
      region        = "europe-west1"
      ip_cidr_range = "10.0.1.0/24"
      iam = {
        "roles/compute.networkUser" = [
          "group:${var.group_email}"
        ]
      }
      iam_bindings = {
        subnet-1-iam = {
          members = ["group:${var.group_email}"]
          role    = "roles/compute.networkUser"
          condition = {
            expression = "resource.matchTag('123456789012/env', 'prod')"
            title      = "test_condition"
          }
        }
      }
    },
    {
      name          = "subnet-2"
      region        = "europe-west1"
      ip_cidr_range = "10.0.2.0/24"
      iam_bindings_additive = {
        subnet-2-iam = {
          member = "group:${var.group_email}"
          role   = "roles/compute.networkUser"
          subnet = "europe-west1/subnet-2"
        }
      }
    }
  ]
}
# tftest modules=1 resources=8 inventory=subnet-iam.yaml e2e
```

### Peering

A single peering can be configured for the VPC, so as to allow management of simple scenarios, and more complex configurations like hub and spoke by defining the peering configuration on the spoke VPCs. Care must be taken so as a single peering is created/changed/destroyed at a time, due to the specific behaviour of the peering API calls.

If you only want to create the "local" side of the peering, use `peering_create_remote_end` to `false`. This is useful if you don't have permissions on the remote project/VPC to create peerings.

```hcl
module "vpc-hub" {
  source     = "./fabric/modules/net-vpc"
  project_id = var.project_id
  name       = "vpc-hub"
  subnets = [{
    ip_cidr_range = "10.0.0.0/24"
    name          = "subnet-1"
    region        = "europe-west1"
  }]
}

module "vpc-spoke-1" {
  source     = "./fabric/modules/net-vpc"
  project_id = var.project_id
  name       = "vpc-spoke1"
  subnets = [{
    ip_cidr_range = "10.0.1.0/24"
    name          = "subnet-2"
    region        = "europe-west1"
  }]
  peering_config = {
    peer_vpc_self_link = module.vpc-hub.self_link
    import_routes      = true
  }
}
# tftest modules=2 resources=10 inventory=peering.yaml
```

### Shared VPC

[Shared VPC](https://cloud.google.com/vpc/docs/shared-vpc) is a project-level functionality which enables a project to share its VPCs with other projects. The `shared_vpc_host` variable is here to help with rapid prototyping, we recommend leveraging the project module for production usage.

```hcl

module "service-project" {
  source          = "./fabric/modules/project"
  billing_account = var.billing_account_id
  name            = "prj1"
  prefix          = var.prefix
  parent          = var.folder_id
  services = [
    "cloudresourcemanager.googleapis.com",
    "compute.googleapis.com",
    "iam.googleapis.com",
    "serviceusage.googleapis.com"
  ]
}

module "vpc-host" {
  source     = "./fabric/modules/net-vpc"
  project_id = var.project_id
  name       = "my-host-network"
  subnets = [
    {
      ip_cidr_range = "10.0.0.0/24"
      name          = "subnet-1"
      region        = "europe-west1"
      secondary_ip_ranges = {
        pods     = "172.16.0.0/20"
        services = "192.168.0.0/24"
      }
      iam = {
        "roles/compute.networkUser" = [
          "serviceAccount:${var.service_account.email}"
        ]
        "roles/compute.securityAdmin" = [
          "serviceAccount:${var.service_account.email}"
        ]
      }
    }
  ]
  shared_vpc_host = true
  shared_vpc_service_projects = [
    module.service-project.project_id
  ]
}
# tftest modules=2 resources=13 inventory=shared-vpc.yaml e2e
```

### Private Service Networking

```hcl
module "vpc" {
  source     = "./fabric/modules/net-vpc"
  project_id = var.project_id
  name       = "my-network"
  subnets = [
    {
      ip_cidr_range = "10.0.0.0/24"
      name          = "production"
      region        = "europe-west1"
    }
  ]
  psa_config = {
    ranges = { myrange = "10.0.1.0/24" }
  }
}
# tftest modules=1 resources=7 inventory=psa.yaml e2e
```

### Private Service Networking with peering routes and peered Cloud DNS domains

Custom routes can be optionally exported/imported through the peering formed with the Google managed PSA VPC.

```hcl
module "vpc" {
  source     = "./fabric/modules/net-vpc"
  project_id = var.project_id
  name       = "my-network"
  subnets = [
    {
      ip_cidr_range = "10.0.0.0/24"
      name          = "production"
      region        = "europe-west1"
    }
  ]
  psa_config = {
    ranges         = { myrange = "10.0.1.0/24" }
    export_routes  = true
    import_routes  = true
    peered_domains = ["gcp.example.com."]
  }
}
# tftest modules=1 resources=8 inventory=psa-routes.yaml e2e
```

### Subnets for Private Service Connect, Proxy-only subnets

Along with common private subnets module supports creation more service specific subnets for the following purposes:

- [Proxy-only subnets](https://cloud.google.com/load-balancing/docs/proxy-only-subnets) for Regional HTTPS Internal HTTPS Load Balancers
- [Private Service Connect](https://cloud.google.com/vpc/docs/private-service-connect#psc-subnets) subnets

```hcl
module "vpc" {
  source     = "./fabric/modules/net-vpc"
  project_id = var.project_id
  name       = "my-network"

  subnets_proxy_only = [
    {
      ip_cidr_range = "10.0.1.0/24"
      name          = "regional-proxy"
      region        = "europe-west1"
      active        = true
    },
    {
      ip_cidr_range = "10.0.4.0/24"
      name          = "global-proxy"
      region        = "australia-southeast2"
      active        = true
      global        = true
    }
  ]
  subnets_psc = [
    {
      ip_cidr_range = "10.0.3.0/24"
      name          = "psc"
      region        = "europe-west1"
    }
  ]
}
# tftest modules=1 resources=6 inventory=proxy-only-subnets.yaml e2e
```

### DNS Policies

```hcl
module "vpc" {
  source     = "./fabric/modules/net-vpc"
  project_id = var.project_id
  name       = "my-network"
  dns_policy = {
    inbound = true
    outbound = {
      private_ns = ["10.0.0.1"]
      public_ns  = ["8.8.8.8"]
    }
  }
  subnets = [
    {
      ip_cidr_range = "10.0.0.0/24"
      name          = "production"
      region        = "europe-west1"
    }
  ]
}
# tftest modules=1 resources=5 inventory=dns-policies.yaml e2e
```

### Subnet Factory

The `net-vpc` module includes a subnet factory (see [Resource Factories](../../blueprints/factories/)) for the massive creation of subnets leveraging one configuration file per subnet. The factory also supports proxy-only and PSC subnets via the `purpose` attribute. The `name` attribute is optional and defaults to the file name, allowing to use the same name for subnets in different regions.

```hcl
module "vpc" {
  source     = "./fabric/modules/net-vpc"
  project_id = var.project_id
  name       = "my-network"
  factories_config = {
    subnets_folder = "config/subnets"
  }
}
# tftest modules=1 resources=10 files=subnet-simple,subnet-simple-2,subnet-detailed,subnet-proxy,subnet-proxy-global,subnet-psc inventory=factory.yaml
```

```yaml
# tftest-file id=subnet-simple path=config/subnets/subnet-simple.yaml
name: simple
region: europe-west4
ip_cidr_range: 10.0.1.0/24
```

```yaml
# tftest-file id=subnet-simple-2 path=config/subnets/subnet-simple-2.yaml
name: simple
region: europe-west8
ip_cidr_range: 10.0.2.0/24
```

```yaml
# tftest-file id=subnet-detailed path=config/subnets/subnet-detailed.yaml
region: europe-west1
description: Sample description
ip_cidr_range: 10.0.0.0/24
# optional attributes
enable_private_access: false  # defaults to true
iam:
  roles/compute.networkUser:
    - group:lorem@example.com
    - serviceAccount:fbz@prj.iam.gserviceaccount.com
    - user:foobar@example.com
secondary_ip_ranges:          # map of secondary ip ranges
  secondary-range-a: 192.168.0.0/24
flow_logs_config:             # enable, set to empty map to use defaults
  aggregation_interval: "INTERVAL_5_SEC"
  flow_sampling: 0.5
  metadata: "INCLUDE_ALL_METADATA"
```

```yaml
# tftest-file id=subnet-proxy path=config/subnets/subnet-proxy.yaml
region: europe-west4
ip_cidr_range: 10.1.0.0/24
proxy_only: true
```

```yaml
# tftest-file id=subnet-proxy-global path=config/subnets/subnet-proxy-global.yaml
region: australia-southeast2
ip_cidr_range: 10.4.0.0/24
proxy_only: true
global: true
```

```yaml
# tftest-file id=subnet-psc path=config/subnets/subnet-psc.yaml
region: europe-west4
ip_cidr_range: 10.2.0.0/24
psc: true
```

### Custom Routes

VPC routes can be configured through the `routes` variable.

```hcl
locals {
  route_types = {
    gateway    = "global/gateways/default-internet-gateway"
    instance   = "zones/europe-west1-b/test"
    ip         = "192.168.0.128"
    ilb        = "regions/europe-west1/forwardingRules/test"
    vpn_tunnel = "regions/europe-west1/vpnTunnels/foo"
  }
}

module "vpc" {
  source     = "./fabric/modules/net-vpc"
  for_each   = local.route_types
  project_id = var.project_id
  name       = "my-network-with-route-${replace(each.key, "_", "-")}"
  routes = {
    next-hop = {
      description   = "Route to internal range."
      dest_range    = "192.168.128.0/24"
      tags          = null
      next_hop_type = each.key
      next_hop      = each.value
    }
    gateway = {
      dest_range    = "0.0.0.0/0",
      priority      = 100
      tags          = ["tag-a"]
      next_hop_type = "gateway",
      next_hop      = "global/gateways/default-internet-gateway"
    }
  }
  create_googleapis_routes = null
}
# tftest modules=5 resources=15 inventory=routes.yaml
```

### Policy Based Routes

Policy based routes can be configured through the `policy_based_routes` variable.

```hcl
module "vpc" {
  source     = "./fabric/modules/net-vpc"
  project_id = var.project_id
  name       = "my-vpc"
  policy_based_routes = {
    skip-pbr-for-nva = {
      use_default_routing = true
      priority            = 100
      target = {
        tags = ["nva"]
      }
    }
    send-all-to-nva = {
      next_hop_ilb_ip = "10.0.0.253"
      priority        = 101
      filter = {
        src_range  = "10.0.0.0/8"
        dest_range = "0.0.0.0/0"
      }
      target = {
        interconnect_attachment = "europe-west8"
      }
    }
  }
  create_googleapis_routes = null
}
# tftest modules=1 resources=3 inventory=pbr.yaml
```

### Private Google Access routes

By default the VPC module creates IPv4 routes for the [Private Google Access ranges](https://cloud.google.com/vpc/docs/configure-private-google-access#config-routing). This behavior can be controlled through the `create_googleapis_routes` variable:

```hcl
module "vpc" {
  source     = "./fabric/modules/net-vpc"
  project_id = var.project_id
  name       = "my-vpc"
  create_googleapis_routes = {
    restricted   = false
    restricted-6 = true
    private      = false
    private-6    = true
  }
}
# tftest modules=1 resources=3 inventory=googleapis.yaml e2e
```

### Allow Firewall Policy to be evaluated before Firewall Rules

```hcl
module "vpc" {
  source                            = "./fabric/modules/net-vpc"
  project_id                        = var.project_id
  name                              = "my-network"
  firewall_policy_enforcement_order = "BEFORE_CLASSIC_FIREWALL"
  subnets = [
    {
      ip_cidr_range = "10.0.0.0/24"
      name          = "production"
      region        = "europe-west1"
      secondary_ip_ranges = {
        pods     = "172.16.0.0/20"
        services = "192.168.0.0/24"
      }
    },
    {
      ip_cidr_range = "10.0.16.0/24"
      name          = "production"
      region        = "europe-west2"
    }
  ]
}
# tftest modules=1 resources=5 inventory=firewall_policy_enforcement_order.yaml e2e
```

### IPv6

A non-overlapping private IPv6 address space can be configured for the VPC via the `ipv6_config` variable. If an internal range is not specified, a unique /48 ULA prefix from the `fd20::/20` range is assigned.

```hcl
module "vpc" {
  source     = "./fabric/modules/net-vpc"
  project_id = var.project_id
  name       = "my-network"
  ipv6_config = {
    # internal_range is optional
    enable_ula_internal = true
    # internal_range      = "fd20:6b2:27e5::/48"
  }
  subnets = [
    {
      ip_cidr_range = "10.0.0.0/24"
      name          = "test"
      region        = "europe-west1"
      ipv6          = {}
    },
    {
      ip_cidr_range = "10.0.1.0/24"
      name          = "test"
      region        = "europe-west3"
      ipv6 = {
        access_type = "EXTERNAL"
      }
    }
  ]
}
# tftest modules=1 resources=5 inventory=ipv6.yaml e2e
```
<!-- BEGIN_TF_DOCS -->

## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.7.0 |
| <a name="requirement_google"></a> [google](#requirement\_google) | >= 5.11.0, < 6.0.0 |
| <a name="requirement_google-beta"></a> [google-beta](#requirement\_google-beta) | >= 5.11.0, < 6.0.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_google"></a> [google](#provider\_google) | >= 5.11.0, < 6.0.0 |
| <a name="provider_google-beta"></a> [google-beta](#provider\_google-beta) | >= 5.11.0, < 6.0.0 |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_name"></a> [name](#input\_name) | The name of the network being created. | `string` | n/a | yes |
| <a name="input_project_id"></a> [project\_id](#input\_project\_id) | The ID of the project where this VPC will be created. | `string` | n/a | yes |
| <a name="input_auto_create_subnetworks"></a> [auto\_create\_subnetworks](#input\_auto\_create\_subnetworks) | Set to true to create an auto mode subnet, defaults to custom mode. | `bool` | `false` | no |
| <a name="input_create_googleapis_routes"></a> [create\_googleapis\_routes](#input\_create\_googleapis\_routes) | Toggle creation of googleapis private/restricted routes. Disabled when vpc creation is turned off, or when set to null. | <pre>object({<br>    private      = optional(bool, true)<br>    private-6    = optional(bool, false)<br>    restricted   = optional(bool, true)<br>    restricted-6 = optional(bool, false)<br>  })</pre> | `{}` | no |
| <a name="input_delete_default_routes_on_create"></a> [delete\_default\_routes\_on\_create](#input\_delete\_default\_routes\_on\_create) | Set to true to delete the default routes at creation time. | `bool` | `false` | no |
| <a name="input_deletion_policy"></a> [deletion\_policy](#input\_deletion\_policy) | The deletion policy for the service networking connection. Setting to ABANDON allows the resource to be abandoned rather than deleted. This will enable a successful terraform destroy when destroying CloudSQL instances. Use with care as it can lead to dangling resources. | `string` | `""` | no |
| <a name="input_description"></a> [description](#input\_description) | An optional description of this resource (triggers recreation on change). | `string` | `"Terraform-managed."` | no |
| <a name="input_dns_policy"></a> [dns\_policy](#input\_dns\_policy) | DNS policy setup for the VPC. | <pre>object({<br>    inbound = optional(bool)<br>    logging = optional(bool)<br>    outbound = optional(object({<br>      private_ns = list(string)<br>      public_ns  = list(string)<br>    }))<br>  })</pre> | `null` | no |
| <a name="input_factories_config"></a> [factories\_config](#input\_factories\_config) | Paths to data files and folders that enable factory functionality. | <pre>object({<br>    subnets_folder = string<br>  })</pre> | `null` | no |
| <a name="input_firewall_policy_enforcement_order"></a> [firewall\_policy\_enforcement\_order](#input\_firewall\_policy\_enforcement\_order) | Order that Firewall Rules and Firewall Policies are evaluated. Can be either 'BEFORE\_CLASSIC\_FIREWALL' or 'AFTER\_CLASSIC\_FIREWALL'. | `string` | `"AFTER_CLASSIC_FIREWALL"` | no |
| <a name="input_ipv6_config"></a> [ipv6\_config](#input\_ipv6\_config) | Optional IPv6 configuration for this network. | <pre>object({<br>    enable_ula_internal = optional(bool)<br>    internal_range      = optional(string)<br>  })</pre> | `{}` | no |
| <a name="input_mtu"></a> [mtu](#input\_mtu) | Maximum Transmission Unit in bytes. The minimum value for this field is 1460 (the default) and the maximum value is 1500 bytes. | `number` | `null` | no |
| <a name="input_peering_config"></a> [peering\_config](#input\_peering\_config) | VPC peering configuration. | <pre>object({<br>    peer_vpc_self_link = string<br>    create_remote_peer = optional(bool, true)<br>    export_routes      = optional(bool)<br>    import_routes      = optional(bool)<br>  })</pre> | `null` | no |
| <a name="input_policy_based_routes"></a> [policy\_based\_routes](#input\_policy\_based\_routes) | Policy based routes, keyed by name. | <pre>map(object({<br>    description         = optional(string, "Terraform-managed.")<br>    labels              = optional(map(string))<br>    priority            = optional(number)<br>    next_hop_ilb_ip     = optional(string)<br>    use_default_routing = optional(bool, false)<br>    filter = optional(object({<br>      ip_protocol = optional(string)<br>      dest_range  = optional(string)<br>      src_range   = optional(string)<br>    }), {})<br>    target = optional(object({<br>      interconnect_attachment = optional(string)<br>      tags                    = optional(list(string))<br>    }), {})<br>  }))</pre> | `{}` | no |
| <a name="input_psa_config"></a> [psa\_config](#input\_psa\_config) | The Private Service Access configuration for Service Networking. | <pre>object({<br>    ranges         = map(string)<br>    export_routes  = optional(bool, false)<br>    import_routes  = optional(bool, false)<br>    peered_domains = optional(list(string), [])<br>  })</pre> | `null` | no |
| <a name="input_routes"></a> [routes](#input\_routes) | Network routes, keyed by name. | <pre>map(object({<br>    description   = optional(string, "Terraform-managed.")<br>    dest_range    = string<br>    next_hop_type = string # gateway, instance, ip, vpn_tunnel, ilb<br>    next_hop      = string<br>    priority      = optional(number)<br>    tags          = optional(list(string))<br>  }))</pre> | `{}` | no |
| <a name="input_routing_mode"></a> [routing\_mode](#input\_routing\_mode) | The network routing mode (default 'GLOBAL'). | `string` | `"GLOBAL"` | no |
| <a name="input_shared_vpc_host"></a> [shared\_vpc\_host](#input\_shared\_vpc\_host) | Enable shared VPC for this project. | `bool` | `false` | no |
| <a name="input_shared_vpc_service_projects"></a> [shared\_vpc\_service\_projects](#input\_shared\_vpc\_service\_projects) | Shared VPC service projects to register with this host. | `list(string)` | `[]` | no |
| <a name="input_subnets"></a> [subnets](#input\_subnets) | Subnet configuration. | <pre>list(object({<br>    name                  = string<br>    ip_cidr_range         = string<br>    region                = string<br>    description           = optional(string)<br>    enable_private_access = optional(bool, true)<br>    flow_logs_config = optional(object({<br>      aggregation_interval = optional(string)<br>      filter_expression    = optional(string)<br>      flow_sampling        = optional(number)<br>      metadata             = optional(string)<br>      # only if metadata == "CUSTOM_METADATA"<br>      metadata_fields = optional(list(string))<br>    }))<br>    ipv6 = optional(object({<br>      access_type = optional(string, "INTERNAL")<br>      # this field is marked for internal use in the API documentation<br>      # enable_private_access = optional(string)<br>    }))<br>    secondary_ip_ranges = optional(map(string))<br><br>    iam = optional(map(list(string)), {})<br>    iam_bindings = optional(map(object({<br>      role    = string<br>      members = list(string)<br>      condition = optional(object({<br>        expression  = string<br>        title       = string<br>        description = optional(string)<br>      }))<br>    })), {})<br>    iam_bindings_additive = optional(map(object({<br>      member = string<br>      role   = string<br>      condition = optional(object({<br>        expression  = string<br>        title       = string<br>        description = optional(string)<br>      }))<br>    })), {})<br>  }))</pre> | `[]` | no |
| <a name="input_subnets_proxy_only"></a> [subnets\_proxy\_only](#input\_subnets\_proxy\_only) | List of proxy-only subnets for Regional HTTPS or Internal HTTPS load balancers. Note: Only one proxy-only subnet for each VPC network in each region can be active. | <pre>list(object({<br>    name          = string<br>    ip_cidr_range = string<br>    region        = string<br>    description   = optional(string)<br>    active        = optional(bool, true)<br>    global        = optional(bool, false)<br><br>    iam = optional(map(list(string)), {})<br>    iam_bindings = optional(map(object({<br>      role    = string<br>      members = list(string)<br>      condition = optional(object({<br>        expression  = string<br>        title       = string<br>        description = optional(string)<br>      }))<br>    })), {})<br>    iam_bindings_additive = optional(map(object({<br>      member = string<br>      role   = string<br>      condition = optional(object({<br>        expression  = string<br>        title       = string<br>        description = optional(string)<br>      }))<br>    })), {})<br>  }))</pre> | `[]` | no |
| <a name="input_subnets_psc"></a> [subnets\_psc](#input\_subnets\_psc) | List of subnets for Private Service Connect service producers. | <pre>list(object({<br>    name          = string<br>    ip_cidr_range = string<br>    region        = string<br>    description   = optional(string)<br><br>    iam = optional(map(list(string)), {})<br>    iam_bindings = optional(map(object({<br>      role    = string<br>      members = list(string)<br>      condition = optional(object({<br>        expression  = string<br>        title       = string<br>        description = optional(string)<br>      }))<br>    })), {})<br>    iam_bindings_additive = optional(map(object({<br>      member = string<br>      role   = string<br>      condition = optional(object({<br>        expression  = string<br>        title       = string<br>        description = optional(string)<br>      }))<br>    })), {})<br>  }))</pre> | `[]` | no |
| <a name="input_vpc_create"></a> [vpc\_create](#input\_vpc\_create) | Create VPC. When set to false, uses a data source to reference existing VPC. | `bool` | `true` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_id"></a> [id](#output\_id) | Fully qualified network id. |
| <a name="output_internal_ipv6_range"></a> [internal\_ipv6\_range](#output\_internal\_ipv6\_range) | ULA range. |
| <a name="output_name"></a> [name](#output\_name) | Network name. |
| <a name="output_network"></a> [network](#output\_network) | Network resource. |
| <a name="output_project_id"></a> [project\_id](#output\_project\_id) | Project ID containing the network. Use this when you need to create resources *after* the VPC is fully set up (e.g. subnets created, shared VPC service projects attached, Private Service Networking configured). |
| <a name="output_self_link"></a> [self\_link](#output\_self\_link) | Network self link. |
| <a name="output_subnet_ids"></a> [subnet\_ids](#output\_subnet\_ids) | Map of subnet IDs keyed by name. |
| <a name="output_subnet_ips"></a> [subnet\_ips](#output\_subnet\_ips) | Map of subnet address ranges keyed by name. |
| <a name="output_subnet_ipv6_external_prefixes"></a> [subnet\_ipv6\_external\_prefixes](#output\_subnet\_ipv6\_external\_prefixes) | Map of subnet external IPv6 prefixes keyed by name. |
| <a name="output_subnet_regions"></a> [subnet\_regions](#output\_subnet\_regions) | Map of subnet regions keyed by name. |
| <a name="output_subnet_secondary_ranges"></a> [subnet\_secondary\_ranges](#output\_subnet\_secondary\_ranges) | Map of subnet secondary ranges keyed by name. |
| <a name="output_subnet_self_links"></a> [subnet\_self\_links](#output\_subnet\_self\_links) | Map of subnet self links keyed by name. |
| <a name="output_subnets"></a> [subnets](#output\_subnets) | Subnet resources. |
| <a name="output_subnets_proxy_only"></a> [subnets\_proxy\_only](#output\_subnets\_proxy\_only) | L7 ILB or L7 Regional LB subnet resources. |
| <a name="output_subnets_psa"></a> [subnets\_psa](#output\_subnets\_psa) | Private Service Access range for Service Networking. |
| <a name="output_subnets_psc"></a> [subnets\_psc](#output\_subnets\_psc) | Private Service Connect subnet resources. |
<!-- END_TF_DOCS -->
