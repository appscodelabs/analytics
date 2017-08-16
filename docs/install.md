> New to Analytics? Please start [here](/docs/tutorial.md).

# Installation Guide

## Using YAML
Analytics can be installed using YAML files includes in the [/hack/deploy](/hack/deploy) folder.

```console
$ curl https://raw.githubusercontent.com/appscode/analytics/0.1.0/hack/deploy/analytics.yaml \
  | kubectl apply -f -
```

## Verify installation
To check if Analytics operator pods have started, run the following command:
```console
$ kubectl get pods --all-namespaces -l app=analytics --watch
```

Once the operator pods are running, you can cancel the above command by typing `Ctrl+C`.
