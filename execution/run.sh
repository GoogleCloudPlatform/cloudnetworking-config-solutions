#!/bin/bash
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

set -eo pipefail  # Exit on error or pipe failure
# Initialize default values for the flags/variables
stage=""
tfcommand="init"

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

# Define valid stages to be accespted by the -s flag
valid_stages="all organization networking security/alloydb security/mrc security/cloudsql security/gce producer/alloydb producer/mrc producer/cloudsql producer/gke producer/vectorsearch producer/onlineendpoint networking-manual consumer/gce consumer/cloudrun/job consumer/cloudrun/service"

# Define valid Terraform commands to be accepted by the -tf or --tfcommand flag
valid_tf_commands="init apply apply-auto-approve destroy destroy-auto-approve init-apply init-apply-auto-approve"

# Define stage to path mapping (excluding "all")
stage_path_map=(
    "organization=01-organization"
    "networking=02-networking"
    "security/alloydb=03-security/AlloyDB"
    "security/mrc=03-security/MRC"
    "security/cloudsql=03-security/CloudSQL"
    "security/gce=03-security/GCE"
    "producer/alloydb=04-producer/AlloyDB"
    "producer/mrc=04-producer/MRC"
    "producer/cloudsql=04-producer/CloudSQL"
    "producer/gke=04-producer/GKE"
    "producer/vectorsearch=04-producer/VectorSearch"
    "producer/onlineendpoint=04-producer/Vertex-AI-Online-Endpoints"
    "networking-manual=05-networking-manual"
    "consumer/gce=06-consumer/GCE"
    "consumer/cloudrun/job=06-consumer/CloudRun/Job"
    "consumer/cloudrun/service=06-consumer/CloudRun/Service"
)

# Define tfvars to stage path mapping (excluding "all")
stagewise_tfvar_path_map=(
    "01-organization=../../configuration/organization.tfvars"
    "02-networking=../../configuration/networking.tfvars"
    "03-security/AlloyDB=../../../configuration/security/alloydb.tfvars"
    "03-security/MRC=../../../configuration/security/mrc.tfvars"
    "03-security/CloudSQL=../../../configuration/security/cloudsql.tfvars"
    "03-security/GCE=../../../configuration/security/gce.tfvars"
    "04-producer/AlloyDB=../../../configuration/producer/AlloyDB/alloydb.tfvars"
    "04-producer/MRC=../../../configuration/producer/MRC/mrc.tfvars"
    "04-producer/CloudSQL=../../../configuration/producer/CloudSQL/cloudsql.tfvars"
    "04-producer/GKE=../../../configuration/producer/GKE/gke.tfvars"
    "04-producer/VectorSearch=../../../configuration/producer/VectorSearch/vectorsearch.tfvars"
    "04-producer/Vertex-AI-Online-Endpoints=../../../configuration/producer/Vertex-AI-Online-Endpoints/vertex-ai-online-endpoints.tfvars"
    "05-networking-manual=../../configuration/networking-manual.tfvars"
    "06-consumer/GCE=../../../configuration/consumer/GCE/gce.tfvars"
    "06-consumer/CloudRun/Job=../../../../configuration/consumer/CloudRun/Job/cloudrunjob.tfvars"
    "06-consumer/CloudRun/Service=../../../../configuration/consumer/CloudRun/Service/cloudrunservice.tfvars"
)

# Define stage to description mapping (excluding "all")
stage_wise_description_map=(
  "all=Progresses through each stage individually."
  "organization=Executes 01-organization stage, manages Google Cloud APIs."
  "networking=Executes 02-networking stage, manages network resources."
  "security/alloydb=Executes 03-security/AlloyDB stage, manages AlloyDB firewall rules."
  "security/mrc=Executes 03-security/MRC stage, manages MRC firewall rules."
  "security/cloudsql=Executes 03-security/CloudSQL stage, manages CloudSQL firewall rules."
  "security/gce=Executes 03-security/GCE stage, manages GCE firewall rules."
  "producer/alloydb=Executes 04-producer/AlloyDB stage, manages AlloyDB instance."
  "producer/mrc=Executes 04-producer/MRC stage, manages MRC instance."
  "producer/cloudsql=Executes 04-producer/CloudSQL stage, manages CloudSQL instance."
  "producer/gke=Executes 04-producer/GKE stage, manages GKE clusters."
  "producer/vectorsearch=Executes 04-producer/VectorSearch stage, manages Vector Search instances."
  "producer/onlineendpoint=Executes 04-producer/Vertex-AI-Online-Endpoints stage, manages Online endpoints."
  "networking-manual=Executes 05-networking-manual stage, manages PSC for supported services."
  "consumer/gce=Executes 06-consumer/GCE stage, manages GCE instance."
  "consumer/cloudrun/job=Executes 06-consumer/CloudRun/Job, manages Cloud Run jobs."
  "consumer/cloudrun/service=Executes 06-consumer/CloudRun/Service, manages Cloud Run services."
  )

# Define tfcommand to description mapping.
tfcommand_wise_description_map=(
    "init=Prepare your working directory for other commands."
    "apply=Create or update infrastructure."
    "apply-auto-approve=Create or Update infrastructure, skips user input."
    "destroy=Destroy previously-created infrastructure."
    "destroy-auto-approve=Destroy previously-created infrastructure, skips user input."
    "init-apply=Prepares working directory and creates/updates infrastructure."
    "init-apply-auto-approve=Prepares working directory and creates/updates infrastructure, skips user input."
)

# Function to get the value associated with a key present in the *_map variables created
function get_value {
  local key="$1"
  local map_name="$2"    # Name of the map array
  local map_ref
  eval "map_ref=(\"\${$map_name[@]}\")"

  # Iterate directly over the elements of the array
  for pair in "${map_ref[@]}"; do
    key_from_map="${pair%%=*}"       # Extract key (part before '=')

    if [[ "$key_from_map" == "$key" ]]; then
      value="${pair#*=}"
      echo "${value}"
    fi
  done
}

# Displays the table formatting.
tableprint() {
    printf "\t\t "
    printf "~%.0s" {1..109}
    printf "\n"
}

# Describing the usage of the run.sh shell script.
usage() {
  printf "Usage: $0 [\033[1m-s|--stage\033[0m <stage>] [[\033[1m-t|--tfcommand\033[0m <command>] [\033[1m-h|--help\033[0m]\n"
  printf " \033[1m-h, --help\033[0m              Displays the detailed help.\n"
  printf " \033[1m-s, --stage\033[0m             STAGENAME to be executed (STAGENAME is case insensitive). e.g. '-s all'  \n\t Valid options are: \n"
  tableprint
  printf "\t\t |%-25s| %-80s|\n" "STAGENAME" "Description"
  tableprint
  for stage_name in $valid_stages; do
    value=$(get_value $stage_name "stage_wise_description_map")
    printf "\t\t |%-25s| %-80s|\n" "$stage_name"  "$value"
  done
  tableprint
  printf " \033[1m-t, --tfcommand\033[0m         TFCOMMAND to be executed (TFCOMMAND is case insensitive). e.g. '-t init' \n\t Valid options are: \n"
  tableprint
  printf "\t\t |%-25s| %-80s|\n" "TFCOMMAND" "Description"
  tableprint
  for tfcommand_value in $valid_tf_commands; do
    value=$(get_value $tfcommand_value "tfcommand_wise_description_map")
    printf "\t\t |%-25s| %-80s|\n" "$tfcommand_value"  "$value"
  done
  tableprint
}

# This function asks for a confirmation before a user provides a auto-approve functionality
confirm() {
    while true; do
        echo -e "${RED} [WARNING] : This action modifies existing resources on all stages without further confirmation. Proceed with caution..${NC}"
        read -p "Do you want to continue. Please answer y or n. $1 (y/n) " confirmation_input
        case $confirmation_input in
            [Yy]* ) break;;
            [Nn]* ) exit 1;;
            * ) echo "Please answer yes or no.";;
        esac
    done
}

# Handle arguments
while [[ $# -gt 0 ]]; do
    case "$1" in
        -s | --stage)
            stage="$2"
            if [[ ! " $valid_stages " =~ " $stage " ]]; then
                printf "${RED}Error: Invalid stage '$stage'. Valid options are: '${valid_stages// /\',\'}' ${NC}" >&2
                exit 1
            fi
            shift 2 ;;
        -t | --tfcommand)
            tfcommand="$2"
            if [[ ! " $valid_tf_commands " =~ " $tfcommand " ]]; then
                printf "${RED}Error: Invalid Terraform command '$tfcommand'. Valid options are: '${valid_tf_commands// /\',\'}' ${NC}" >&2
                exit 1
            fi
            shift 2 ;;
        -h | --help)
            usage
            exit 0 ;;
        *)
            echo "Invalid option: $1" >&2
            usage
            exit 1 ;;
    esac
done

# Shift to remove processed options from positional arguments
shift $((OPTIND-1))

# Error handling: Check if both flags are provided
if [ -z "$stage" ] || [ -z "$tfcommand" ]; then
  usage
  exit 1
fi

# Execute Terraform commands based on the stage and tfcommand
if [[ $stage == "all" ]]; then
  # Handles the execution of all the stages one by one if the stage="all" value is passed
  # Create an array of stage paths in the correct sequential order ensuring the incremental order is maintained
  stage_path_array=()
  for stage_name in $valid_stages; do
    if [[ $stage_name != "all" ]]; then
      stage_path_value=$(get_value $stage_name "stage_path_map")
      stage_path_array+=("${stage_path_value}")
    fi
  done
  # Determine execution order based on tfcommand, reverse the order of execution if the -tf/--tfcommand contain destroy/deletion instructions
  if [[ $tfcommand == destroy || $tfcommand == destroy-auto-approve ]]; then
     for (( i=${#stage_path_array[@]}-1; i>=0; i-- )); do
        reversed_array+=("${stage_path_array[i]}")
    done
    stage_path_array=("${reversed_array[@]}")
  fi

  # Present a warning if a user uses auto-approve flag
  if [[ $tfcommand =~ "auto-approve" ]]; then
    confirm
  fi

  # Iterate over stages in the determined order
  for stage_path in "${stage_path_array[@]}"; do
    echo -e "Executing Terraform command(s) in ${GREEN}$stage_path${NC}..."
    tfvar_file_path=$(get_value $stage_path  "stagewise_tfvar_path_map")
    echo "tfvars file path : ${tfvar_file_path}"
    (cd "$stage_path" &&
      case "$tfcommand" in
          init) terraform init -var-file=$tfvar_file_path ;;
          apply) terraform apply -var-file=$tfvar_file_path ;;
          apply-auto-approve) terraform apply --auto-approve -var-file=$tfvar_file_path ;;
          destroy) terraform destroy -var-file=$tfvar_file_path ;;
          destroy-auto-approve) terraform destroy -var-file=$tfvar_file_path --auto-approve ;;
          init-apply) terraform init && terraform apply -var-file=$tfvar_file_path ;;
          init-apply-auto-approve) terraform init && terraform apply -var-file=$tfvar_file_path --auto-approve ;;
          *) echo "${RED}Error: Invalid tfcommand '$tfcommand'${NC}" >&2; exit 1 ;;
      esac
    )
  done
else
  # Otherwise, get the path for the specified stage. logic for single stage execution
  stage_path=$(get_value $stage  "stage_path_map")
  tfvar_file_path=$(get_value $stage_path  "stagewise_tfvar_path_map")
  echo "tfvars file path : ${tfvar_file_path}"
  if [[ -z "$stage_path" ]]; then  # Check if a path was found
      echo "${RED}: Unexpected error finding path for stage '$stage'${NC}" >&2
      exit 1
  else
    echo "Executing Terraform command(s) in $stage_path..."
    (cd "$stage_path" &&
      case "$tfcommand" in
          init) terraform init -var-file=$tfvar_file_path;;
          apply) terraform apply -var-file=$tfvar_file_path ;;
          apply-auto-approve) terraform apply -var-file=$tfvar_file_path --auto-approve ;;
          destroy) terraform destroy -var-file=$tfvar_file_path ;;
          destroy-auto-approve) terraform destroy -var-file=$tfvar_file_path --auto-approve ;;
          init-apply) terraform init && terraform apply -var-file=$tfvar_file_path ;;
          init-apply-auto-approve) terraform init && terraform apply -var-file=$tfvar_file_path --auto-approve ;;
          *) echo "${RED}Error: Invalid tfcommand '$tfcommand'${NC}" >&2; exit 1 ;;
      esac
    )
  fi
fi
