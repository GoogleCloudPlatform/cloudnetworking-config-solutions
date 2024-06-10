# Configuration files

This directory serves as a centralized repository for all Terraform configuration files (.tfvars) used across the various stages of your infrastructure deployment. By organizing these configuration files in one place, we maintain a clear and structured approach to managing environment-specific variables and settings.

## File Organization by Stage

- 00-bootstrap stage (Filename : bootstrap.tfvars)
- 01-organisation stage (Filename : organisation.tfvars)
- 02-networking stage (Filename : networking.tfvars)
- 03-security stage 
        - AlloyDB (alloydb-firewall.tfvars)
        - MRC (mrc-firewall.tfvars)
        - Cloud SQL (sql-firewall.tfvars)
        - GCE (gce-firewall.tfvars)
- 06-networking-manual stage (Filename : psc-manual.tfvars)

# Usage

## Specifying Variable Files

When executing a Terraform stage (e.g., plan, apply, destroy), you must explicitly instruct Terraform to use the corresponding configuration file. This is achieved using the `-var-file` flag followed by the relative path to the .tfvars file.

## Relative Paths

Relative paths are essential for maintaining flexibility and ensuring your Terraform configuration works across different environments. While running any of the stages, use the [-var-file flag](https://developer.hashicorp.com/terraform/language/values/variables#variable-definitions-tfvars-files) to give relative path of the .tfvars file. Let's assume you're within the networking directory and want to execute terraform plan using the networking.tfvars configuration file:

```none
terraform plan -var-file=../config-files/networking.tfvars
```

This would run the terraform plan based on the vars provided in the networking.tfvars file in the `config-files` folder. In this example:

- `-var-file` instructs Terraform to load variables from the specified file.
- `../` moves up one directory level from networking.
- `config-files/networking.tfvars points` to the exact location of the configuration file.

## Benefits of Centralized Configuration

- Improved Readability: A dedicated directory makes it easy to locate and manage configuration files.
- Enhanced Maintainability: Changes to environment-specific variables can be made in one place, minimizing the risk of errors.
- Streamlined Collaboration: Team members can easily access and understand the configuration structure.
- Simplified Automation: Terraform workflows can automatically reference the appropriate configuration file based on the stage being executed.

## Considerations

<<<<<<< PATCH SET (6a4c22 Adding the PSC-Manual stage and stage wise description to th)
- Sensitive Data: If your configuration files contain securrely handle sensitive values (e.g., API keys) and ensure they are securely stored. We strong recommend to not store senstive information in plain text and suggest you to carefully manage sensitive information.

# Stage wise details

## 03-Security

1. `project_id` : this variable identifies the GCP project where the firewall rule will be created. 

2. `network` : this variable specifies the name of the Virtual Private Cloud (VPC) network to which the firewall rule will be applied. Firewall rules control traffic flow in and out of your VPC network.

3. `egress_rules` : this variable defines a set of egress (outbound) firewall rules. These rules determine what kind of outgoing traffic is permitted from your VPC network to destinations outside the network.

***Example Usage** 

```
project_id              = "my-project"
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

## 06-Networking-Manual (psc-manual.tfvars)

Defined using `psc_endpoints` which is a list of PSC endpoint configurations consisting of:

1. `project_id` : Consumer project ID (where the forwarding rule is created).

2. `producer_project_id` : Project where the producer service such as Cloud SQL is created.

2. `producer_instance_name` : Name of the producer service instance.

3. `subnetwork_name` : this variable names the specific subnetwork within your Virtual Private Cloud (VPC) network from which the internal IP address for the PSC connection will be allocated.

4. `network_name` : VPC network for the forwarding rule which hosts the subnetwork mentioned above.

5. `ip_address` : **(Optional)** Specific internal IP, or leave null for automatic allocation.

**Example Usage**

```
psc_endpoints = [
  {
    project_id             = "project-for-endpoint"
    producer_project_id    = "project-producer"
    producer_instance_name = "sql-1"
    subnetwork_name        = "subnetwork-1"
    network_name           = "network-1"
    ip_address             = "10.128.0.5"
  },
  # Add more endpoint objects as needed
]
```
=======
- Sensitive Data: If your configuration files contain securrely handle sensitive values (e.g., API keys) and ensure they are securely stored. We strongly recommend to not store senstive information in plain text and suggest you to carefully manage sensitive information.
>>>>>>> BASE      (5798f1 Merge "Implementing the security stage for MRC and config-fi)
