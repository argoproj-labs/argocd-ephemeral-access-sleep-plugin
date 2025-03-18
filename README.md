# Argo CD EphemeralAccess Sleep Plugin

A simple Ephemeral Access plugin for tests and serve as an implementation reference.

## How it works

When users request to have their access elevated it will:

1. Initially reply the request as Pending
1. After 30 seconds it will reply:
  - If the current minute is an even number it retuns Approved
  - If the current minute is on odd number it returns Denied

## Using it

The steps below will deploy the EphemeralAccess in a Kubernetes
cluster with the Sleep Plugin configured.

Run the command:

    docker build -t argocd-sleep-plugin:latest .

Apply the manifests in the cluster:

    kustomize build ./manifests | kubectl apply -f -
