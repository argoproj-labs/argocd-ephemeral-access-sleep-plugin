# Argo CD EphemeralAccess Sleep Plugin

A simple [Ephemeral Access][1] plugin for tests and serve as an
implementation reference. It demostrates how to implement the
[AccessRequester][2] interface and make it available in a docker image
so it installs itself using the init-container pattern.

## How it works

When users request to have their access elevated it will:

1. Initially reply the request as Pending
1. After 15 seconds it will reply:
  - If the current minute is an even number it retuns Approved
  - If the current minute is on odd number it returns Denied

## Using it

The steps below will deploy the EphemeralAccess in a Kubernetes
cluster with the Sleep Plugin configured.

Run the command:

    docker build -t argocd-sleep-plugin:latest .

Apply the manifests in the cluster:

    kustomize build ./manifests | kubectl apply -f -

[1]: https://github.com/argoproj-labs/argocd-ephemeral-access/tree/main
[2]: https://github.com/argoproj-labs/argocd-ephemeral-access/blob/532ff45b5a6824d941101c21c59433b7a37ec9d8/pkg/plugin/rpc.go#L44-L48
