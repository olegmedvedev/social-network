#!/bin/sh
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install social-postgres bitnami/postgresql --namespace social-network -f k8s/postgres-values.yaml 