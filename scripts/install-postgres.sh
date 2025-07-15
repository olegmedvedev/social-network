#!/bin/sh
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install social-postgres bitnami/postgresql --namespace social-network --set auth.username=socialuser --set auth.password=socialpass --set auth.database=socialdb 