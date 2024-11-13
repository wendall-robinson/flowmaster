# OpenTelemetry Collector Kubernetes Deployment Example

This repository provides a basic example of how to deploy the OpenTelemetry Collector and Zipkin in a Kubernetes environment. This setup is intended for developers who want to get started with the OpenTelemetry Collector and visualization in Kubernetes without needing in-depth knowledge of distributed tracing or telemetry systems.

## Introduction

This example demonstrates how to deploy an OpenTelemetry Collector and Zipkin in a Kubernetes cluster using Kubernetes manifests. The Collector is configured to receive telemetry data over OTLP (gRPC and HTTP), process it using a batch processor, and export it to stdout for simplicity.

## Prerequisites
* Kubernetes Cluster: A running Kubernetes cluster. You can use [Minikube](https://minikube.sigs.k8s.io/docs/) for local testing.
* [kubectl](https://kubernetes.io/docs/reference/kubectl/): Command-line tool for Kubernetes. Ensure it's configured to communicate with your cluster.
* [Helm](https://helm.sh/) to deploy the helm charts

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

3. Deploy the OpenTelemetry Collector and Zipkin with Helm

    Apply the Kubernetes manifests to deploy the Collector and Zipkin:
    ```bash
    helm install otel ./examples/helm
    ```


## Testing the Deployment
### Verify the Pods

Check that the Collector pod is running:
```bash
kubectl get pods
```

You should see output similar to:
```bash
NAME                              READY   STATUS              RESTARTS   AGE
otel-collector-858f6dbd5c-vc9tl   1/1     ContainerCreating   0          30s
zipkin-7b8f845b86-mvt5t           1/1     ContainerCreating   0          30s
```

### View Logs

Since the Collector exports data to stdout, you can view the logs to see incoming telemetry data:
```bash
kubectl logs otel-collector-858f6dbd5c-vc9tl
```

### Connect to Zipkin
You will first need to port-forward to the Zipkin container
```bash
kubectl port-forward szipkin-7b8f845b86-mvt5t 9411:9411
```

#### Navigate to Zipkin Endpoint
```
http://localhost:9411/zipkin/
```

### Cleanup

To remove all the resources created by this example:
```bash
helm uninstall otel ./examples/helm
```
