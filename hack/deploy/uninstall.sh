#!/bin/bash
set -x

kubectl delete deployment -l app=analytics -n kube-system
kubectl delete service -l app=analytics -n kube-system

# Delete RBAC objects, if --rbac flag was used.
kubectl delete serviceaccount -l app=analytics -n kube-system
kubectl delete clusterrolebindings -l app=analytics -n kube-system
kubectl delete clusterrole -l app=analytics -n kube-system
