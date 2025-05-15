package client

import (
	"context"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	intrv1 "github.com/rwpp/RzWeLook/api/proto/gen/intr/v1"
	"google.golang.org/grpc"
	"math/rand"
)

type GreyScaleInteractiveServiceClient struct {
	remote    intrv1.InteractiveServiceClient
	local     intrv1.InteractiveServiceClient
	threshold *atomicx.Value[int32]
}

func NewGreyScaleInteractiveServiceClient(remote intrv1.InteractiveServiceClient, local intrv1.InteractiveServiceClient) *GreyScaleInteractiveServiceClient {
	return &GreyScaleInteractiveServiceClient{
		remote:    remote,
		local:     local,
		threshold: atomicx.NewValue[int32](),
	}
}

func (g *GreyScaleInteractiveServiceClient) IncrReadCnt(ctx context.Context, in *intrv1.IncrReadCntRequest, opts ...grpc.CallOption) (*intrv1.IncrReadCntResponse, error) {
	return g.client().IncrReadCnt(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) Like(ctx context.Context, in *intrv1.LikeRequest, opts ...grpc.CallOption) (*intrv1.LikeResponse, error) {
	return g.client().Like(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) CancelLike(ctx context.Context, in *intrv1.CancelLikeRequest, opts ...grpc.CallOption) (*intrv1.CancelLikeResponse, error) {
	return g.client().CancelLike(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) Collect(ctx context.Context, in *intrv1.CollectRequest, opts ...grpc.CallOption) (*intrv1.CollectResponse, error) {
	return g.client().Collect(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) Get(ctx context.Context, in *intrv1.GetRequest, opts ...grpc.CallOption) (*intrv1.GetResponse, error) {
	return g.client().Get(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) GetByIds(ctx context.Context, in *intrv1.GetByIdRequest, opts ...grpc.CallOption) (*intrv1.GetByIdResponse, error) {
	return g.client().GetByIds(ctx, in, opts...)
}
func (g *GreyScaleInteractiveServiceClient) UpdateThreshold(newthreshold int32) {
	g.threshold.Store(newthreshold)
}
func (g *GreyScaleInteractiveServiceClient) client() intrv1.InteractiveServiceClient {
	threshold := g.threshold.Load()
	num := rand.Int31n(100)
	if num < threshold {
		return g.remote
	} else {
		return g.local
	}
}
