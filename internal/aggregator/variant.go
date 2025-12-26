package aggregator

import (
	"context"
	"log/slog"

	apiv1 "buf.build/gen/go/catou/transit-radar/protocolbuffers/go/api/v1"
	"connectrpc.com/connect"
	"github.com/catouberos/transit-watcher/providers/gobus"
	"google.golang.org/protobuf/types/known/structpb"
)

func (a *aggregator) processVariant(ctx context.Context, routeID string, variant gobus.RouteVariant, description string) (*apiv1.Variant, error) {
	response, err := a.variantService.ListVariants(ctx, connect.NewRequest(&apiv1.ListVariantsRequest{
		Filter: "",
	}))
	if err != nil {
		slog.InfoContext(ctx, "cannot get variant by ebms ID", "error", err)
		return nil, err
	}

	if len(response.Msg.Variants) == 0 {
		// create
		slog.InfoContext(ctx, "no existing variant, creating one", "ebmsID", variant.Id)
		return a.createVariant(ctx, routeID, variant)
	}

	apiVariant := response.Msg.Variants[0]
	slog.InfoContext(ctx, "updating existing variant", "id", apiVariant.Id)
	return a.updateVariant(ctx, routeID, apiVariant, variant)
}

func (a *aggregator) createVariant(ctx context.Context, routeID string, variant gobus.RouteVariant) (*apiv1.Variant, error) {
	attributes, err := structpb.NewStruct(map[string]any{
		"ebmsID": variant.Id,
	})
	if err != nil {
		return nil, err
	}

	apiVariant, err := a.variantService.CreateVariant(ctx, connect.NewRequest(&apiv1.CreateVariantRequest{
		RouteId:    routeID,
		Name:       variant.Name,
		ShortName:  &variant.ShortName,
		Distance:   &variant.Distance,
		Direction:  0,
		Duration:   &variant.Duration,
		Attributes: attributes,
	}))
	if err != nil {
		return nil, err
	}

	return apiVariant.Msg.Variant, nil
}

func (a *aggregator) updateVariant(ctx context.Context, routeID string, existing *apiv1.Variant, variant gobus.RouteVariant) (*apiv1.Variant, error) {
	apiVariant, err := a.variantService.UpdateVariant(ctx, &connect.Request[apiv1.UpdateVariantRequest]{
		Msg: &apiv1.UpdateVariantRequest{
			Id:        existing.Id,
			Name:      &variant.Name,
			RouteId:   &routeID,
			ShortName: &variant.ShortName,
			Distance:  &variant.Distance,
			Duration:  &variant.Duration,
		},
	})
	if err != nil {
		return nil, err
	}

	return apiVariant.Msg.Variant, nil
}
