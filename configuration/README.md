# Configuration files

This directory serves as a centralized repository for all Terraform configuration files (.tfvars) used across the various stages of your infrastructure deployment. By organizing these configuration files in one place, we maintain a clear and structured approach to managing environment-specific variables and settings.

## File Organization by Stage

- 00-bootstrap stage (Filename : bootstrap.tfvars)
- 01-organisation stage (Filename : organisation.tfvars)
- 02-networking stage (Filename : networking.tfvars)
- 03-security stage 
        - MRC (Filename : mrc-firewall.tfvars)
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

- Sensitive Data: If your configuration files contain securrely handle sensitive values (e.g., API keys) and ensure they are securely stored. We strongly recommend to not store senstive information in plain text and suggest you to carefully manage sensitive information.