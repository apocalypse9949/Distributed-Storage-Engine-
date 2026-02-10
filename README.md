# Distributed Storage Engine

## Introduction

This project implements a distributed key-value storage engine designed for reliable and scalable data management. It provides a gRPC interface for fundamental storage operations (Get, Set, Delete) and is packaged for deployment within a Kubernetes environment. The core focus is on demonstrating robust data persistence with secondary indexing capabilities.

## Architecture

The Distributed Storage Engine is structured around several key internal components, exposing its functionality primarily through a gRPC API.

### Core Components

*   **`internal/storage`**: Manages the underlying data persistence. It utilizes BadgerDB as its storage backend, offering key-value storage with support for secondary indexes.
*   **`internal/server`**: Implements the gRPC server logic. This component handles incoming gRPC requests, translating them into operations on the `internal/storage` layer.
*   **`internal/cluster` (Placeholder)**: Intended for managing cluster membership and coordination among multiple nodes of the storage engine.
*   **`internal/replication` (Placeholder)**: Designed to handle data redundancy and fault tolerance across the cluster.
*   **`internal/routing` (Placeholder)**: Responsible for directing client requests to the appropriate node within a distributed setup.
*   **`cmd/distributed-storage-engine`**: Contains the main application entry point, responsible for initializing the storage engine and starting the gRPC server.

### Data Model

The engine primarily operates on a key-value data model, where both keys and values are byte arrays. It enhances this with a secondary indexing mechanism, allowing for efficient retrieval of primary keys based on indexed values.

### gRPC Interface

The primary interaction point with the storage engine is a gRPC (Google Remote Procedure Call) API. This interface is defined using Protocol Buffers and provides the following remote procedures:

*   `Get(GetRequest)`: Retrieves a value associated with a given key.
*   `Set(SetRequest)`: Stores a key-value pair, with an option to create a secondary index.
*   `Delete(DeleteRequest)`: Removes a key-value pair.

### Protocol Buffers

Protocol Buffers are used for defining the service interface (`.proto` files) and the structure of messages exchanged between clients and the server. This ensures efficient, language-agnostic data serialization and communication.

### Kubernetes Deployment

The application is containerized using Docker and configured for deployment on Kubernetes. It utilizes a Deployment to manage instances of the storage engine and a Service to expose the gRPC API to other services within or outside the cluster.

## How It Works

Upon startup, the `cmd/distributed-storage-engine` application initializes an instance of `internal/storage.BadgerStorage`. This storage instance is then passed to `internal/server.PalantirServer`, which registers itself with a new gRPC server. The server then begins listening for incoming gRPC requests on a specified port (default `50051`).

*   **Data Storage**: Data is persistently stored using BadgerDB, an embedded key-value store optimized for fast read/write operations.
*   **Secondary Indexes**: When a `Set` operation includes index details, the storage layer creates a composite index key (`_idx:indexName:indexValue:primaryKey`) alongside the primary data entry. The `GetByIndex` operation efficiently retrieves primary keys by querying these composite index keys.
*   **gRPC Communication**: Clients interact with the engine by sending gRPC requests, which are handled by the `internal/server.PalantirServer` methods. These methods internally call the corresponding operations on the `internal/storage` layer.
*   **Distributed Aspects**: While the current implementation focuses on the core storage and gRPC interface, the `internal/cluster`, `internal/replication`, and `internal/routing` packages are structured to facilitate future expansion into a truly distributed, fault-tolerant system.

## Getting Started

To set up and run the Distributed Storage Engine, follow these steps:

### Prerequisites

Ensure you have the following tools installed and configured on your system:

*   **Go (1.24.5 or later)**: The programming language used for the project.
*   **Docker**: For building the application container image.
*   **`protoc` (Protocol Buffers Compiler)**: For generating gRPC code from `.proto` definitions.
*   **`protoc-gen-go` and `protoc-gen-go-grpc`**: Go plugins for `protoc`.
    *   Install with: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
    *   Install with: `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`
*   **`kubectl`**: The Kubernetes command-line tool, configured to connect to your Kubernetes cluster.

### Building the Project

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/apocalypse9949/Distributed-Storage-Engine-.git
    cd Distributed-Storage-Engine-
    ```
2.  **Update Go modules:**
    ```bash
    go mod tidy
    ```
3.  **Re-generate gRPC code (if `.proto` files are modified):**
    ```bash
    protoc --proto_path=api/proto --go_out=api --go-grpc_out=api api/proto/distributed-storage-engine.proto
    ```
4.  **Build the Docker image:**
    ```bash
    docker build -t distributed-storage-engine-image:latest .
    ```

### Deployment to Kubernetes

1.  **Apply Kubernetes manifests:**
    ```bash
    kubectl apply -f deploy/kubernetes/distributed-storage-engine.yaml
    ```
2.  **Verify Deployment:**
    ```bash
    kubectl get deployments
    kubectl get pods -l app=distributed-storage-engine
    kubectl get services
    ```
3.  **Accessing the Service:**
    *   If using `LoadBalancer`, retrieve the external IP: `kubectl get service distributed-storage-engine -o jsonpath='{.status.loadBalancer.ingress[0].ip}'`
    *   For other environments (e.g., Minikube), you might need `kubectl port-forward` to access the service locally.

## Benchmarking

Detailed instructions for benchmarking the deployed Distributed Storage Engine can be found in the `benchmark_instructions.md` file within this repository. This guide includes steps for verifying deployment, obtaining service access points, and suggestions for gRPC benchmarking tools like `ghz`.

