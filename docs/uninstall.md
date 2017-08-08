# Uninstall Analytics
Please follow the steps below to uninstall Analytics:

1. Delete the various objects created for Analytics operator.
```console
$ ./hack/deploy/uninstall.sh
+ kubectl delete deployment -l app=analytics -n kube-system
deployment "analytics" deleted
+ kubectl delete service -l app=analytics -n kube-system
service "analytics" deleted
+ kubectl delete serviceaccount -l app=analytics -n kube-system
No resources found
+ kubectl delete clusterrolebindings -l app=analytics -n kube-system
No resources found
+ kubectl delete clusterrole -l app=analytics -n kube-system
No resources found
```

2. Now, wait several seconds for Analytics to stop running. To confirm that Analytics operator pod(s) have stopped running, run:
```console
$ kubectl get pods --all-namespaces -l app=analytics
```
