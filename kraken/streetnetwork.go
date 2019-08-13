package kraken

import (
	"github.com/CanalTP/gonavitia/pbnavitia"
	"github.com/golang/protobuf/proto"
)

type DirectPathBuilder struct {
	Kraken Kraken
	From   string
	To     string
}

func (b DirectPathBuilder) Get() (*pbnavitia.Response, error) {
	mode := "walking"
	speed := 1.11

	request := &pbnavitia.Request{
		RequestedApi: pbnavitia.API_direct_path.Enum(),
		DirectPath: &pbnavitia.DirectPathRequest{
			Origin:      &pbnavitia.LocationContext{Place: proto.String(b.From), AccessDuration: proto.Int32(0)},
			Destination: &pbnavitia.LocationContext{Place: proto.String(b.To), AccessDuration: proto.Int32(0)},
			StreetnetworkParams: &pbnavitia.StreetNetworkParams{
				OriginMode:   proto.String(mode),
				WalkingSpeed: proto.Float64(speed),
			},
			Clockwise: proto.Bool(true),
		},
	}
	return b.Kraken.Call(request)
}

type StreetNetworkMatrixBuilder struct {
	Kraken      Kraken
	From        []string
	To          []string
	MaxDuration int32
}

func (b StreetNetworkMatrixBuilder) Get() (*pbnavitia.Response, error) {
	mode := "walking"
	var speed float32 = 1.11
	var maxDuration int32 = 86400
	if b.MaxDuration > 0 {
		maxDuration = b.MaxDuration
	}

	snMatrix := &pbnavitia.StreetNetworkRoutingMatrixRequest{
		Mode:        proto.String(mode),
		Speed:       proto.Float32(speed),
		MaxDuration: proto.Int32(maxDuration),
	}
	for _, item := range b.From {
		snMatrix.Origins = append(snMatrix.Origins, &pbnavitia.LocationContext{Place: proto.String(item), AccessDuration: proto.Int32(0)})
	}
	for _, item := range b.To {
		snMatrix.Destinations = append(snMatrix.Destinations, &pbnavitia.LocationContext{Place: proto.String(item), AccessDuration: proto.Int32(0)})
	}

	request := &pbnavitia.Request{
		RequestedApi:    pbnavitia.API_street_network_routing_matrix.Enum(),
		SnRoutingMatrix: snMatrix,
	}
	return b.Kraken.Call(request)
}
