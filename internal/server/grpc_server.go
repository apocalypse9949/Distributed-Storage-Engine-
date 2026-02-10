package server

import (
	"context"

	"github.com/gemini-cli/distributed-storage-engine/api"
	"github.com/gemini-cli/distributed-storage-engine/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PalantirServer implements the api.PalantirServer interface
type PalantirServer struct {
	api.UnimplementedPalantirServer
	store *storage.BadgerStorage
}

// NewPalantirServer creates a new PalantirServer
func NewPalantirServer(s *storage.BadgerStorage) *PalantirServer {
	return &PalantirServer{
		store: s,
	}
}

// Get retrieves a value for a given key
func (s *PalantirServer) Get(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	if req.Key == nil || len(req.Key) == 0 {
		return nil, status.Error(codes.InvalidArgument, "key cannot be empty")
	}

	value, err := s.store.Get(req.Key)
	if err != nil {
		if err == storage.ErrKeyNotFound {
			return nil, status.Error(codes.NotFound, "key not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get key: %v", err)
	}

	return &api.GetResponse{Value: value}, nil
}

// Set sets a value for a given key
func (s *PalantirServer) Set(ctx context.Context, req *api.SetRequest) (*api.SetResponse, error) {
	if req.Key == nil || len(req.Key) == 0 {
		return nil, status.Error(codes.InvalidArgument, "key cannot be empty")
	}
	if req.Value == nil {
		return nil, status.Error(codes.InvalidArgument, "value cannot be nil")
	}

	// Assuming a Set operation does not return the previous value by default in this context.
	// If it needs to return something, the proto should be updated and logic here.
	return &api.SetResponse{}, nil
}

// Delete deletes a key-value pair
func (s *PalantirServer) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	if req.Key == nil || len(req.Key) == 0 {
		return nil, status.Error(codes.InvalidArgument, "key cannot be empty")
	}

	return &api.DeleteResponse{}, nil
}
