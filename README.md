# Cluster Readers Operator

[![Docker Repository on Quay](https://quay.io/repository/jharrington22/cluster-readers-operator/status "Docker Repository on Quay")](https://quay.io/repository/jharrington22/cluster-readers-operator)

The goal of this operator is to monitor a cluster role binding and ensure that only the users listed in the cluster resource (CR) are "User" subjects within the cluster role binding and have the role "cluster-reader" assigned.

Should either the CR or ClusterRoleBinding change it will be replaced with what is in the CR.
