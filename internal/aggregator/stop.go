package aggregator

import (
	"context"
	"log/slog"

	apiv1 "buf.build/gen/go/catou/transit-radar/protocolbuffers/go/api/v1"
	"connectrpc.com/connect"
	"github.com/catouberos/transit-watcher/providers/gobus"
)

func (a *aggregator) processStop(ctx context.Context, variantID int64, stop gobus.Stop) (int64, error) {
	ebmsID, err := stop.Property.Id.Int64()
	if err != nil {
		return 0, err
	}

	response, err := a.stopService.GetStopByEbmsID(ctx, &connect.Request[apiv1.GetStopByEbmsIDRequest]{
		Msg: &apiv1.GetStopByEbmsIDRequest{EbmsId: ebmsID},
	})
	if err != nil {
		// continue-able error
		slog.InfoContext(ctx, "cannot get stop by ebms ID", "error", err)
	}

	if response == nil {
		slog.InfoContext(ctx, "no existing stop, attempt to create", "ebmsID", ebmsID)
		_, _ = a.stopService.CreateStop(ctx, &connect.Request[apiv1.CreateStopRequest]{
			Msg: &apiv1.CreateStopRequest{
				Code: stop.Property.Code,
				Name: stop.Property.Name,
				// TODO: get type
				TypeId:    0,
				EbmsId:    ebmsID,
				Active:    true,
				Latitude:  stop.Geometry.Coordinates[1],
				Longitude: stop.Geometry.Coordinates[0],
			},
		})
	}

	return 0, nil
}
