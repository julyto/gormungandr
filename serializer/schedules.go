package serializer

import (
	"strings"

	"github.com/CanalTP/gonavitia"
	"github.com/CanalTP/gonavitia/pbnavitia"
	"github.com/golang/protobuf/proto"
)

const maxUInt64 = 18446744073709551615

func (s *Serializer) NewRouteSchedulesResponse(request *pbnavitia.Request, pb *pbnavitia.Response) *gonavitia.RouteScheduleResponse {
	if pb == nil {
		return nil
	}
	s.Init(pb.Metadatas)
	response := gonavitia.RouteScheduleResponse{
		Error:          s.NewError(pb.Error),
		RouteSchedules: make([]*gonavitia.RouteSchedule, 0, len(pb.RouteSchedules)),
		Pagination:     s.NewPagination(pb.Pagination),
		FeedPublishers: make([]*gonavitia.FeedPublisher, 0, len(pb.FeedPublishers)),
		Links:          make([]gonavitia.Link, 0),
		Context:        s.NewContext(request, pb),
		Exceptions:     make([]struct{}, 0),
		Notes:          make([]struct{}, 0),
		Disruptions:    make([]struct{}, 0),
	}
	for _, r := range pb.RouteSchedules {
		response.RouteSchedules = append(response.RouteSchedules, s.NewRouteSchedule(r))
	}
	for _, f := range pb.FeedPublishers {
		response.FeedPublishers = append(response.FeedPublishers, s.NewFeedPublisher(f))
	}
	return &response
}

func (s *Serializer) NewRouteSchedule(pb *pbnavitia.RouteSchedule) *gonavitia.RouteSchedule {
	if pb == nil {
		return nil
	}
	var additionalInfo *string
	info := pb.ResponseStatus
	if info != nil {
		tmp := strings.ToLower(info.Enum().String())
		additionalInfo = &tmp
	}
	return &gonavitia.RouteSchedule{
		DisplayInfo:    s.NewPtDisplayInfoForRoute(pb.PtDisplayInformations),
		Table:          s.NewTable(pb.Table),
		AdditionalInfo: additionalInfo,
		GeoJson:        s.NewGeoJsonMultistring(pb.Geojson),
		Links:          s.NewLinksFromUris(pb.PtDisplayInformations),
	}
}

func (s *Serializer) NewTable(pb *pbnavitia.Table) *gonavitia.Table {
	if pb == nil {
		return nil
	}
	t := gonavitia.Table{
		Headers: make([]*gonavitia.Header, 0, len(pb.Headers)),
		Rows:    make([]gonavitia.Row, 0, len(pb.Rows)),
	}
	for _, h := range pb.Headers {
		t.Headers = append(t.Headers, s.NewHeader(h))
	}
	for _, r := range pb.Rows {
		t.Rows = append(t.Rows, s.NewRow(r))
	}
	return &t
}

func (s *Serializer) NewHeader(pb *pbnavitia.Header) *gonavitia.Header {
	if pb == nil {
		return nil
	}
	header := gonavitia.Header{
		DisplayInfo:     s.NewPtDisplayInfoForVJ(pb.PtDisplayInformations),
		Links:           s.NewLinksFromUris(pb.PtDisplayInformations),
		AdditionalInfos: s.NewAdditionalInformations(pb.AdditionalInformations),
	}
	return &header
}

func (s *Serializer) NewRow(pb *pbnavitia.RouteScheduleRow) gonavitia.Row {
	if pb == nil {
		return gonavitia.Row{}
	}
	r := gonavitia.Row{
		StopPoint: s.NewStopPoint(pb.StopPoint),
		DateTimes: make([]gonavitia.DateTime, 0, len(pb.DateTimes)),
	}
	for _, d := range pb.DateTimes {
		r.DateTimes = append(r.DateTimes, s.NewDatetime(d))
	}
	return r
}

func (s *Serializer) NewDatetime(pb *pbnavitia.ScheduleStopTime) gonavitia.DateTime {
	if pb == nil {
		return gonavitia.DateTime{}
	}
	if pb.GetTime() == maxUInt64 {
		//This is an "empty" datetime cell in the response
		return gonavitia.DateTime{
			AdditionalInfo: s.NewAdditionalInformation(pb.Properties.AdditionalInformations),
			Links:          s.NewLinksFromProperties(pb.Properties),
		}
	}

	rtLevel := strings.ToLower(pb.GetRealtimeLevel().Enum().String())
	baseDateTime := s.NewNavitiaDatetime(int64(pb.GetBaseDateTime()))
	return gonavitia.DateTime{
		DateTime:       s.NewNavitiaDatetime(int64(pb.GetDate() + pb.GetTime())),
		BaseDateTime:   &baseDateTime,
		AdditionalInfo: s.NewAdditionalInformation(pb.Properties.AdditionalInformations),
		DataFreshness:  &rtLevel,
		Links:          s.NewLinksFromProperties(pb.Properties),
	}
}

func (s *Serializer) NewLinksFromProperties(pb *pbnavitia.Properties) []gonavitia.Link {
	result := make([]gonavitia.Link, 0, 1)
	if pb == nil {
		return result
	}
	if pb.VehicleJourneyId != nil {
		result = append(result, gonavitia.Link{
			Id:    pb.VehicleJourneyId,
			Value: pb.VehicleJourneyId,
			Type:  proto.String("vehicle_journey"),
			Rel:   proto.String("vehicle_journeys"),
		})
	}
	return result
}

func (s *Serializer) NewAdditionalInformation(pb []pbnavitia.Properties_AdditionalInformation) []string {
	infos := make([]string, 0, len(pb))
	for _, v := range pb {
		additionalInfo := strings.ToLower(v.Enum().String())
		infos = append(infos, additionalInfo)
	}
	return infos
}
