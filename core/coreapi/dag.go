package coreapi

import (
	"context"

	cid "github.com/ipfs/go-cid"
	pin "github.com/ipfs/go-ipfs-pinner"
	ipld "github.com/ipfs/go-ipld-format"
	dag "github.com/ipfs/go-merkledag"
	"github.com/xbradylee/ipfs-kubo/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type dagAPI struct {
	ipld.DAGService

	core *CoreAPI
}

type pinningAdder CoreAPI

func (adder *pinningAdder) Add(ctx context.Context, nd ipld.Node) error {
	ctx, span := tracing.Span(ctx, "CoreAPI.PinningAdder", "Add", trace.WithAttributes(attribute.String("node", nd.String())))
	defer span.End()
	defer adder.blockstore.PinLock(ctx).Unlock(ctx)

	if err := adder.dag.Add(ctx, nd); err != nil {
		return err
	}

	adder.pinning.PinWithMode(nd.Cid(), pin.Recursive)

	return adder.pinning.Flush(ctx)
}

func (adder *pinningAdder) AddMany(ctx context.Context, nds []ipld.Node) error {
	ctx, span := tracing.Span(ctx, "CoreAPI.PinningAdder", "AddMany", trace.WithAttributes(attribute.Int("nodes.count", len(nds))))
	defer span.End()
	defer adder.blockstore.PinLock(ctx).Unlock(ctx)

	if err := adder.dag.AddMany(ctx, nds); err != nil {
		return err
	}

	cids := cid.NewSet()

	for _, nd := range nds {
		c := nd.Cid()
		if cids.Visit(c) {
			adder.pinning.PinWithMode(c, pin.Recursive)
		}
	}

	return adder.pinning.Flush(ctx)
}

func (api *dagAPI) Pinning() ipld.NodeAdder {
	return (*pinningAdder)(api.core)
}

func (api *dagAPI) Session(ctx context.Context) ipld.NodeGetter {
	return dag.NewSession(ctx, api.DAGService)
}

var (
	_ ipld.DAGService  = (*dagAPI)(nil)
	_ dag.SessionMaker = (*dagAPI)(nil)
)
