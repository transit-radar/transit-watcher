package aggregator

import (
	"context"
	"log/slog"
	"regexp"

	apiv1 "buf.build/gen/go/catou/transit-radar/protocolbuffers/go/api/v1"
	"connectrpc.com/connect"
	"github.com/catouberos/transit-watcher/providers/gobus"
)

func (a *aggregator) processRoute(ctx context.Context, route gobus.Route) (int64, error) {
	ebmsID, err := route.Id.Int64()
	if err != nil {
		return 0, err
	}

	response, err := a.routeService.GetRouteByEbmsID(ctx, &connect.Request[apiv1.GetRouteByEbmsIDRequest]{
		Msg: &apiv1.GetRouteByEbmsIDRequest{EbmsId: ebmsID},
	})
	if err != nil {
		// continue-able error
		slog.Debug("cannot get route by ebms ID", "error", err)
	}

	// known markup characters regex
	re, err := regexp.Compile(`<br\/?>|  +|&nbsp`)
	if err != nil {
		return 0, err
	}

	route.Info.Organization = re.ReplaceAllString(route.Info.Organization, "")
	route.Info.Ticketing = re.ReplaceAllString(route.Info.Ticketing, "")

	// create when upstream don't have data
	if response == nil {
		slog.Debug("creating new route", "ebmsID", route.Id)
		apiRoute, err := a.routeService.CreateRoute(ctx, &connect.Request[apiv1.CreateRouteRequest]{
			Msg: &apiv1.CreateRouteRequest{
				Number:        route.Number,
				Name:          route.Name,
				EbmsId:        &ebmsID,
				OperationTime: &route.Info.OperationTime,
				Organization:  &route.Info.Organization,
				Ticketing:     &route.Info.Ticketing,
				RouteType:     &route.Info.RouteType,
			},
		})
		if err != nil {
			slog.Error("cannot create new route", "error", err)
			return 0, err
		}
		return apiRoute.Msg.Route.Id, nil
	}

	// check if upstream and current data is outdated
	if shouldUpdateRoute(&route, response.Msg.Route) == false {
		return response.Msg.Route.Id, nil
	}

	// update the upstream with current data
	active := true
	slog.Debug("updating existing route", "id", response.Msg.Route.Id, "ebmsID", route.Id)
	apiRoute, err := a.routeService.UpdateRoute(ctx, &connect.Request[apiv1.UpdateRouteRequest]{
		Msg: &apiv1.UpdateRouteRequest{
			Id:            response.Msg.Route.Id,
			Number:        &route.Number,
			Name:          &route.Name,
			Active:        &active,
			OperationTime: &route.Info.OperationTime,
			Organization:  &route.Info.Organization,
			Ticketing:     &route.Info.Ticketing,
			RouteType:     &route.Info.RouteType,
		},
	})
	if err != nil {
		slog.Error("cannot update new route", "routeID", response.Msg.Route.Id, "error", err)
		return 0, err
	}
	return apiRoute.Msg.Route.Id, nil
}

func shouldUpdateRoute(route *gobus.Route, apiRoute *apiv1.Route) bool {
	if apiRoute.Name == route.Name &&
		apiRoute.Number == route.Number &&
		apiRoute.OperationTime != nil && *apiRoute.OperationTime == route.Info.OperationTime &&
		apiRoute.Organization != nil && *apiRoute.Organization == route.Info.Organization &&
		apiRoute.RouteType != nil && *apiRoute.RouteType == route.Info.RouteType &&
		apiRoute.Ticketing != nil && *apiRoute.Ticketing == route.Info.Ticketing &&
		apiRoute.Active == true {

		return false
	}

	return true
}
