# Configuration files

This directory serves as a centralized repository for all Terraform configuration files (.tfvars) used across the various stages of your infrastructure deployment. By organizing these configuration files in one place, we maintain a clear and structured approach to managing environment-specific variables and settings.

## File Organization by Stage

- 00-bootstrap stage (bootstrap.tfvars)
- 01-organisation stage (organisation.tfvars)
- 02-networking stage (networking.tfvars)
- 03-security stage
    - AlloyDB (alloydb.tfvars)
    - MRC (mrc.tfvars)
    - Cloud SQL (sql.tfvars)
    - GCE (gce.tfvars)
- 05-networking-manual stage (networking-manual.tfvars)

# Usage

## Specifying Variable Files

When executing a Terraform stage (e.g. terraform plan, terraform apply, terraform destroy), you must explicitly instruct Terraform to use the corresponding configuration file. This is achieved using the `-var-file` flag followed by the relative path to the .tfvars file.

## Relative Paths

Relative paths are essential for maintaining flexibility and ensuring your Terraform configuration works across different environments. While running any of the stages, use the [-var-file flag](https://developer.hashicorp.com/terraform/language/values/variables#variable-definitions-tfvars-files) to give relative path of the .tfvars file. Let's assume you're within the networking directory and want to execute terraform plan using the networking.tfvars configuration file:

```none
terraform plan -var-file=../configuration/networking.tfvars
```

This would run the terraform plan based on the values for the variables declared in the networking.tfvars file in the `configuration` folder. In this example:

- `-var-file` : instructs Terraform to load variables from the specified file.
- `../` : moves up one directory level from networking.
- `configuration/networking.tfvars` : points to the configuration folder containing the networking.tfvars file.

## Benefits of Centralized Configuration

- Improved Readability: A dedicated directory makes it easy to locate and manage configuration files.
- Enhanced Maintainability: Changes to environment-specific variables can be made in one place, minimizing the risk of errors.
- Streamlined Collaboration: Team members can easily access and understand the configuration structure.
- Simplified Automation: Terraform workflows can automatically reference the appropriate configuration file based on the stage being executed.


# Stage wise details

## 00-bootstrap

  - This tfvars file provides project IDs and administrator email addresses for different stages of the infrastructure setup. These values are used by Terraform to configure resources and permissions in the respective Google Cloud projects.
  - bootstrap project ID (`bootstrap_project_id`): project used to create resources such as service accounts or grant permissions to users to run the stages.
  - networking projects (`network_hostproject_id`/`network_serviceproject_id`) : host/service project IDs
  - Administrators : in stage wise administrator variables, you can set user accounts/groups to delegate permissions.


**Example usage**

```
bootstrap_project_id                  = "test-bootstrap-project"
network_hostproject_id                = "host-project-id"
network_serviceproject_id             = "consumer-project-id"
organization_stage_administrator      = ["example@example.com"]
networking_stage_administrator        = ["example@example.com"]
security_stage_administrator          = ["example@example.com"]
producer_stage_administrator          = ["example@example.com"]
networking_manual_stage_administrator = ["example@example.com"]
consumer_stage_administrator          = ["example@example.com"]
```

## 01-organization

**Example usage**

```
  activate_api_identities = {
    "project-01" = {
      project_id = "test-project",
      activate_apis = [
        "servicenetworking.googleapis.com",
        "alloydb.googleapis.com",
        "sqladmin.googleapis.com",
        "iam.googleapis.com",
        "compute.googleapis.com",
        "redis.googleapis.com",
        "aiplatform.googleapis.com",
        "container.googleapis.com",
        "run.googleapis.com",
      ],
    },
  }
```

## 02-networking

**Example usage**

```
project_id = "test-project"
region     = "us-central1"

## VPC input variables

network_name = "network-test"
subnets = [
  {
    ip_cidr_range = "10.0.0.0/24"
    name          = "subnet-test"
    region        = "us-central1"
  }
]


create_scp_policy      = true
subnets_for_scp_policy = ["subnet-test"]

create_nat = true

create_havpn = false
peer_gateways = {
  default = {
    gcp = "" # e.g. projects/<google-cloud-peer-projectid>/regions/<google-cloud-region>/vpnGateways/<peer-vpn-name>
  }
}

tunnel_1_router_bgp_session_range = ""
tunnel_1_bgp_peer_asn             = 64514
tunnel_1_bgp_peer_ip_address      = ""
tunnel_1_shared_secret            = ""

tunnel_2_router_bgp_session_range = ""
tunnel_2_bgp_peer_asn             = 64514
tunnel_2_bgp_peer_ip_address      = ""
tunnel_2_shared_secret            = ""
```

## 03-security

  - `project_id` : this variable identifies the GCP project where the firewall rule will be created.
  - `network` : this variable specifies the name of the Virtual Private Cloud (VPC) network to which the firewall rule will be applied. Firewall rules control traffic flow in and out of your VPC network.
  - `egress_rules/ingress_rules` : this variable defines a set of egress (outbound) firewall rules. These rules determine what kind of outgoing traffic is permitted from your VPC network to destinations outside the network.

***Example Usage**

```
project_id              = "project-id"
network                 = "network-name"
egress_rules = {
  allow-egress = {
    deny = false
    rules = [{
      protocol = "tcp"
      ports    = ["6379"]
    }]
  }
}
```

## 04-producer

Producer specific configuration examples can be found under the `/config` folder of that specific producer. Such as for AlloyDB, the example would be in the folder `04-producer/AlloyDB/config`.

## 05-networking-manual (networking-manual.tfvars)

Defined using `psc_endpoints` which is a list of PSC endpoint configurations consisting of:

1. `endpoint_project_id` : Consumer project ID (where the forwarding rule is created).

2. `producer_instance_project_id` : Project where the producer service such as Cloud SQL is created.

2. `producer_instance_name` : Name of the producer service instance.

3. `subnetwork_name` : this variable names the specific subnetwork within your Virtual Private Cloud (VPC) network from which the internal IP address for the PSC connection will be allocated.

4. `network_name` : VPC network for the forwarding rule which hosts the subnetwork mentioned above.

5. `ip_address_literal` : **(Optional)** Specific internal IP, or leave null for automatic allocation.

**Example Usage**

```
psc_endpoints = [
  {
    endpoint_project_id          = "endpoint-project-id"
    producer_instance_project_id = "instance-project-id"
    producer_instance_name       = "sql-1"
    subnetwork_name              = "subnetwork-1"
    network_name                 = "network-1"
    ip_address_literal           = "10.128.0.50"
  },
  # Add more endpoint objects as needed
]
```

## 06-consumer

Consumer specific configuration examples can be found under the `/config` folder of that specific consumer. Such as for GCE, the example would be in the folder `06-consumer/GCE/config`.

## Considerations

- Sensitive Data: If your configuration files contain securrely handle sensitive values (e.g., API keys) and ensure they are securely stored. We strong recommend to not store senstive information in plain text and suggest you to carefully manage sensitive information.
