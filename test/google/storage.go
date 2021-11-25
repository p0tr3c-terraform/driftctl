package google

import (
	"context"
	"net"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

type FakeStorageServer struct {
	routes map[string]http.HandlerFunc
}

func (s *FakeStorageServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for path, handler := range s.routes {
		if path == r.RequestURI {
			handler(w, r)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func NewFakeStorageServer(routes map[string]http.HandlerFunc) (*storage.Client, error) {
	return newStorageClient(&FakeStorageServer{routes: routes})
}

func newStorageClient(fakeServer *FakeStorageServer) (*storage.Client, error) {
	ctx := context.Background()
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}
	fakeServerAddr := l.Addr().String()
	go func() {
		if err := http.Serve(l, fakeServer); err != nil {
			panic(err)
		}
	}()

	// This prevents the client to use HTTPS
	_ = os.Setenv("STORAGE_EMULATOR_HOST", fakeServerAddr)

	// Create a client.
	client, err := storage.NewClient(ctx,
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithInsecure()),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}
