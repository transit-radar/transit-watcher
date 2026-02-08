package utils

import (
	"errors"
	"regexp"
	"strings"

	radarv1 "buf.build/gen/go/transit-radar/apis/protocolbuffers/go/transit/radar/v1"
	"github.com/catouberos/transit-watcher/providers/gobus"
)

/*
func (a *aggregator) processRoute(ctx context.Context, route *gobus.Route) (*apiv1.Route, error) {
	ebmsID, err := route.Id.Int64()
	if err != nil {
		return nil, err
	}

	response, err := a.routeService.ListRoutes(ctx, connect.NewRequest(&apiv1.ListRoutesRequest{
		Filter: fmt.Sprintf("attributes.ebmsID = %d", ebmsID),
	}))
	if err != nil {
		slog.DebugContext(ctx, "cannot get route by ebms ID", "error", err)
		return nil, err
	}

	// known markup characters regex
	re, err := regexp.Compile(`<br\/?>|  +|&nbsp`)
	if err != nil {
		return nil, err
	}

	route.Info.Organization = re.ReplaceAllString(route.Info.Organization, "")
	route.Info.Ticketing = re.ReplaceAllString(route.Info.Ticketing, "")

	routeType, err := routeType(route)
	if err != nil {
		return nil, err
	}

	// create when upstream don't have data
	if len(response.Msg.Routes) == 0 {
		slog.Debug("creating new route", "ebmsID", route.Id)
		return a.createRoute(ctx, route, routeType)
	}

	existing := response.Msg.Routes[0] // get first route

	// check if upstream and current data is outdated
	if shouldUpdateRoute(route, existing) {
		slog.Debug("updating existing route", "id", existing.Id, "ebmsID", route.Id)
		return a.updateRoute(ctx, existing, route, routeType)
	}

	return existing, nil
}
*/

func RouteType(route *gobus.Route) (radarv1.RouteType, error) {
	// matches HCMC Metro
	if strings.HasPrefix(route.Number, "MRT") {
		return radarv1.RouteType_ROUTE_TYPE_METRO, nil
	}

	// matches Saigon Waterbus
	if strings.HasPrefix(route.Number, "SWB") {
		return radarv1.RouteType_ROUTE_TYPE_FERRY, nil
	}

	// matches Public Tour Bus
	if strings.HasPrefix(route.Number, "DL") {
		return radarv1.RouteType_ROUTE_TYPE_BUS, nil
	}

	// exclusively match Vinbus "D-" Bus
	if strings.HasPrefix(route.Number, "D") {
		return radarv1.RouteType_ROUTE_TYPE_BUS, nil
	}

	// matches Public Bus
	re, err := regexp.Compile("[0-9]+(-[0-9]+)?[A-Z]?")
	if err != nil {
		return radarv1.RouteType_ROUTE_TYPE_UNSPECIFIED, err
	}
	if re.MatchString(route.Number) {
		return radarv1.RouteType_ROUTE_TYPE_BUS, nil
	}

	return radarv1.RouteType_ROUTE_TYPE_UNSPECIFIED, errors.New("unhandled route type")
}
