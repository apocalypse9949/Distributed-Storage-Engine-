The Palantir gRPC server has been successfully deployed to your Kubernetes cluster.

Here are the instructions to benchmark its performance:

### 1. Verify Kubernetes Deployment Status

First, ensure your deployment and service are running correctly:

*   **Check Deployment Status:**
    ```bash
    kubectl get deployments
    kubectl describe deployment palantir
    ```
*   **Check Pod Status:**
    ```bash
    kubectl get pods -l app=palantir
    kubectl logs -f <pod-name> # Replace <pod-name> with the actual pod name
    ```
*   **Check Service Status:**
    ```bash
    kubectl get services
    kubectl describe service palantir
    ```

### 2. Get the Service's External IP/Hostname

If you used `type: LoadBalancer`, it might take a few moments for an external IP to be provisioned.

*   **Get External IP:**
    ```bash
    kubectl get service palantir -o jsonpath='{.status.loadBalancer.ingress[0].ip}'
    ```
    If `ip` is empty, your cluster might provision a `hostname` instead:
    ```bash
    kubectl get service palantir -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'
    ```
    If neither an IP nor hostname is available (e.g., in minikube or some on-prem clusters), you might need to use `kubectl port-forward` or configure `type: NodePort` to access the service.

### 3. Choose a Benchmarking Tool

You'll need a tool capable of generating gRPC traffic. Here are a few options:

*   **`ghz` (gRPC benchmarking utility):** A powerful tool specifically designed for gRPC.
    *   **Installation:** `go install github.com/bojand/ghz/cmd/ghz@latest`
    *   **Example Usage (assuming external IP is 192.168.1.100):**
        ```bash
        # For the Get RPC
        ghz --proto api/proto/palantir.proto --call api.Palantir.Get --insecure --host 192.168.1.100:50051 --data '{"key":"some_key"}' -n 1000 -c 50

        # For the Set RPC
        ghz --proto api/proto/palantir.proto --call api.Palantir.Set --insecure --host 192.168.1.100:50051 --data '{"key":"new_key", "value":"new_value"}' -n 1000 -c 50
        ```
        Adjust `--data` to match your `GetRequest`, `SetRequest`, etc., and define your load with `-n` (number of requests) and `-c` (concurrency).

*   **Custom Go Client:** Write a simple Go client that continuously makes gRPC calls. This gives you maximum flexibility to simulate specific load patterns.

*   **`k6` (with gRPC support):** A modern load testing tool.
    *   **Example (JavaScript with gRPC):** You'd write a script to define your gRPC calls.

### 4. Monitor Performance

While benchmarking, you should monitor your Kubernetes cluster and the Palantir application's metrics.

*   **Kubernetes Metrics:** Use `kubectl top pods` and `kubectl top nodes` (if Metrics Server is installed) to observe CPU/memory usage.
*   **Application Logs:** Keep an eye on the pod logs for any errors or performance degradation messages.

### 5. Clean Up Kubernetes Resources

When you're done with benchmarking and testing, remember to clean up the deployed resources:

```bash
kubectl delete -f deploy/kubernetes/palantir.yaml
```

This will delete the `palantir` Deployment and Service from your cluster.
