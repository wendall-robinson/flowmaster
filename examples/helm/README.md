# OpenTelemetry Collector Kubernetes Deployment Example

This repository provides a basic example of how to deploy the OpenTelemetry Collector in a Kubernetes environment. This setup is intended for developers who want to get started with the OpenTelemetry Collector in Kubernetes without needing in-depth knowledge of distributed tracing or telemetry systems.

## Introduction

This example demonstrates how to deploy an OpenTelemetry Collector in a Kubernetes cluster using Kubernetes manifests. The Collector is configured to receive telemetry data over OTLP (gRPC and HTTP), process it using a batch processor, and export it to stdout for simplicity.

## Prerequisites
* Kubernetes Cluster: A running Kubernetes cluster. You can use [Minikube](https://minikube.sigs.k8s.io/docs/) for local testing.
* [kubectl](https://kubernetes.io/docs/reference/kubectl/): Command-line tool for Kubernetes. Ensure it's configured to communicate with your cluster.

## Setup Instructions
1. Clone the Repository

    Clone this repository to your local machine:
    ```bash
    git clone git@github.com:wendall-robinson/flowmaster.git
    cd flowmaster
    ```

2. Start MiniKube
    ```
    minikube start
    ```

3. Deploy the OpenTelemetry Collector

    Apply the Kubernetes manifests to deploy the Collector:
    ```bash
    kubectl apply -f examples/helm/templates/otel/otel-config.yaml
    kubectl apply -f examples/helm/templates/otel/deployment.yaml
    kubectl apply -f examples/helm/templates/otel/service.yaml
    ```
    This will create:

    A ConfigMap containing the Collector's configuration.
    A Deployment that runs the Collector.
    A Service to expose the Collector within the cluster.

## Testing the Deployment
### Verify the Pods

Check that the Collector pod is running:
```bash
kubectl get pods
```

You should see output similar to:
```bash
NAME                              READY   STATUS    RESTARTS   AGE
otel-collector-xxxxxxxxxx-xxxxx   1/1     Running   0          1m
```

### View Logs

Since the Collector exports data to stdout, you can view the logs to see incoming telemetry data:
```bash
kubectl logs deployment/otel-collector
```

### Cleanup

To remove all the resources created by this example:
```bash
kubectl delete -f examples/helm/templates/otel/service.yaml
kubectl delete -f examples/helm/templates/otel/deployment.yaml
kubectl delete -f examples/helm/templates/otel/otel-config.yaml
```
