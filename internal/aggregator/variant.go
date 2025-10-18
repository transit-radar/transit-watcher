package aggregator

import (
	"context"
	"log/slog"

	apiv1 "buf.build/gen/go/catou/transit-radar/protocolbuffers/go/api/v1"
	"connectrpc.com/connect"
	"github.com/catouberos/transit-watcher/providers/gobus"
)

func (a *aggregator) processVariant(ctx context.Context, routeID int64, variant gobus.RouteVariant, description string) (int64, error) {
	ebmsID, err := variant.Id.Int64()
	if err != nil {
		return 0, err
	}

	response, err := a.variantService.GetVariantByEbmsID(ctx, &connect.Request[apiv1.GetVariantByEbmsIDRequest]{
		Msg: &apiv1.GetVariantByEbmsIDRequest{EbmsId: ebmsID, RouteId: routeID},
	})
	if err != nil {
		// continue-able error
		slog.InfoContext(ctx, "cannot get variant by ebms ID", "error", err)
	}

	if response == nil {
		slog.InfoContext(ctx, "no existing variant, creating one", "ebmsID", variant.Id)
		apiVariant, err := a.variantService.CreateVariant(ctx, &connect.Request[apiv1.CreateVariantRequest]{
			Msg: &apiv1.CreateVariantRequest{
				Name:          variant.Name,
				EbmsId:        &ebmsID,
				IsOutbound:    variant.IsOutbound,
				RouteId:       routeID,
				Description:   &description,
				ShortName:     &variant.ShortName,
				Distance:      &variant.Distance,
				Duration:      &variant.Duration,
				StartStopName: &variant.StartStop,
				EndStopName:   &variant.EndStop,
			},
		})
		if err != nil {
			slog.ErrorContext(ctx, "cannot create new variant", "error", err)
			return 0, err
		}
		return apiVariant.Msg.Variant.Id, nil
	}

	slog.InfoContext(ctx, "updating existing variant", "id", response.Msg.Variant.Id, "ebmsID", variant.Id)
	apiVariant, err := a.variantService.UpdateVariant(ctx, &connect.Request[apiv1.UpdateVariantRequest]{
		Msg: &apiv1.UpdateVariantRequest{
			Id:            response.Msg.Variant.Id,
			Name:          &variant.Name,
			EbmsId:        &ebmsID,
			IsOutbound:    &variant.IsOutbound,
			RouteId:       &routeID,
			Description:   &description,
			ShortName:     &variant.ShortName,
			Distance:      &variant.Distance,
			Duration:      &variant.Duration,
			StartStopName: &variant.StartStop,
			EndStopName:   &variant.EndStop,
		},
	})
	if err != nil {
		slog.ErrorContext(ctx, "cannot update new variant", "error", err)
		return 0, err
	}
	return apiVariant.Msg.Variant.Id, nil
}
