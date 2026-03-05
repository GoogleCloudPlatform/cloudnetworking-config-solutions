# Copyright 2025 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

if [ -z "$GOOGLE_CLOUD_PROJECT" ]
then
   echo Project not set!
   echo What project do you want to deploy the solution to?
   read var_project_id
   gcloud config set project $var_project_id
   export GOOGLE_CLOUD_PROJECT=$var_project_id
fi

echo Running prerequisites on project $GOOGLE_CLOUD_PROJECT
BUCKET_NAME=$GOOGLE_CLOUD_PROJECT-tf-state-1
if gcloud storage ls gs://$BUCKET_NAME; then
    echo Terraform bucket already created!
else
    echo Creating Terraform state bucket...
    #gcloud storage buckets create gs://$BUCKET_NAME
    gcloud storage buckets create gs://$BUCKET_NAME --project=$GOOGLE_CLOUD_PROJECT --uniform-bucket-level-access
fi

cat > execution/01-organization/providers.tf << EOF
terraform {
  backend "gcs" {
    bucket                      = "$BUCKET_NAME"
    prefix                      = "organization_stage"
  }
}
EOF

cat > execution/02-networking/providers.tf << EOF
terraform {
  backend "gcs" {
    bucket                      = "$BUCKET_NAME"
    prefix                      = "networking_stage"
  }
}
EOF

cat > execution/04-producer/AlloyDB/providers.tf << EOF
terraform {
  backend "gcs" {
    bucket                      = "$BUCKET_NAME"
    prefix                      = "producer/alloydb_stage"
  }
}
EOF

echo Enabling required APIs...
gcloud services enable cloudbuild.googleapis.com \
    cloudresourcemanager.googleapis.com \
    iam.googleapis.com \
    logging.googleapis.com \
    storage.googleapis.com

echo "Granting Cloud Build's Service Account IAM roles to deploy the resources..."
PROJECT_NUMBER=$(gcloud projects describe $GOOGLE_CLOUD_PROJECT --format='value(projectNumber)')
MEMBER=serviceAccount:$PROJECT_NUMBER-compute@developer.gserviceaccount.com
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT --member=$MEMBER --role=roles/logging.logWriter
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT --member=$MEMBER --role=roles/editor
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT --member=$MEMBER --role=roles/iam.securityAdmin
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT --member=$MEMBER --role=roles/compute.networkAdmin
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT --member=$MEMBER --role=roles/secretmanager.secretAccessor
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT --member=$MEMBER --role=roles/iam.serviceAccountUser
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT --member=$MEMBER --role=roles/serviceusage.serviceUsageAdmin
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT --member=$MEMBER --role=roles/storage.objectAdmin
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT --member=$MEMBER --role=roles/alloydb.admin

echo Script completed successfully!
