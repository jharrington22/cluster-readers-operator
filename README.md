# Cluster Readers Operator

The goal of this operator is to monitor a cluster role binding and ensure that only the users in the cluster resource (CR) "ClusterReader" are "User" subjects with with cluster role "cluster-reader" role assigned.

Should a change to either the CR or ClusterRoleBinding change it will be replaced with what is in the CR.