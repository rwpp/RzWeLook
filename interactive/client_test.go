package main

import (
	"context"
	intrv1 "github.com/rwpp/RzWeLook/api/proto/gen/intr/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

func TestGRPCClient(t *testing.T) {
	// Initialize the gRPC client
	cc, err := grpc.Dial("localhost:8090",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	client := intrv1.NewInteractiveServiceClient(cc)
	resp, err := client.Get(context.Background(), &intrv1.GetRequest{
		Biz:   "test",
		BizId: 123,
		Uid:   123,
	})
	require.NoError(t, err)
	t.Log(resp.Intr)
}
