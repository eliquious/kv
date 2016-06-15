#!/bin/bash

[ -z "$PROJECT_NAME" ] && echo "error: missing PROJECT_NAME env variable" && exit 1;
[ -z "$CLUSTER_NAME" ] && echo "error: missing CLUSTER_NAME env variable" && exit 1;
[ -z "$CLUSTER_ZONE" ] && echo "error: missing CLUSTER_ZONE env variable" && exit 1;

# Set project
gcloud config set project ${PROJECT_NAME}

# Set availability zone
gcloud config set compute/zone ${CLUSTER_ZONE}

# Set cluster name
gcloud config set container/cluster ${CLUSTER_NAME}

# Get cluster credentials
gcloud container clusters get-credentials ${CLUSTER_NAME}
