package schedules

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/CanalTP/gonavitia"
	"github.com/CanalTP/gonavitia/pbnavitia"
	"github.com/CanalTP/gormungandr"
	"github.com/CanalTP/gormungandr/serializer"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type RouteScheduleRequest struct {
	FromDatetime     time.Time `form:"from_datetime" time_format:"20060102T150405"`
	DisableGeojson   bool      `form:"disable_geojson"`
	StartPage        int32     `form:"start_page"`
	Count            int32     `form:"count"`
	Duration         int32     `form:"duration"`
	ForbiddenUris    []string  //mapping with Binding doesn't work
	Depth            int32     `form:"depth"`
	CurrentDatetime  time.Time `form:"_current_datetime"`
	ItemsPerSchedule int32     `form:"items_per_schedule"`
	DataFreshness    string    `form:"data_freshness"`
	Filters          []string
	User             gormungandr.User
	Coverage         string //requested coverage
	ID               uuid.UUID
}

func NewRouteScheduleRequest() RouteScheduleRequest {
	return RouteScheduleRequest{
		StartPage:        0,
		Count:            10,
		Duration:         86400,
		CurrentDatetime:  time.Now(),
		Depth:            2,
		ItemsPerSchedule: 10000,
		DataFreshness:    "base_schedudle",
		FromDatetime:     time.Now(),
	}
}

func RouteSchedule(c *gin.Context, kraken *gormungandr.Kraken, request *RouteScheduleRequest, publisher Publisher, logger *logrus.Entry) {
	pbReq := BuildRequestRouteSchedule(*request)
	resp, err := kraken.Call(pbReq)
	logger.Debug("calling kraken")
	if err != nil {
		logger.Errorf("Error while calling kraken: %+v\n", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err})
		return
	}
	logger.Debug("building response")
	r := serializer.NewRouteSchedulesResponse(pbReq, resp)
	fillPaginationLinks(getUrl(c), r)
	status := http.StatusOK
	if r.Error != nil {
		status = r.Error.Code.HTTPCode()
	}
	c.JSON(status, r)
	logger.Debug("handling stats")

	//the original context must not be used in another goroutine, we have to copy it
	readonlyContext := c.Copy()
	go func() {
		err = publisher.PublishRouteSchedule(*request, *r, *readonlyContext)
		if err != nil {
			logger.Errorf("stat not sent %+v", err)
		} else {
			logger.Debug("stat sent")
		}
	}()
}

func BuildRequestRouteSchedule(req RouteScheduleRequest) *pbnavitia.Request {
	departureFilter := strings.Join(req.Filters, "and ")
	//TODO handle Realtime level from request
	pbReq := &pbnavitia.Request{
		RequestedApi: pbnavitia.API_ROUTE_SCHEDULES.Enum(),
		NextStopTimes: &pbnavitia.NextStopTimeRequest{
			DepartureFilter:  proto.String(departureFilter),
			ArrivalFilter:    proto.String(""),
			FromDatetime:     proto.Uint64(uint64(req.FromDatetime.Unix())),
			Duration:         proto.Int32(req.Duration),
			Depth:            proto.Int32(req.Depth),
			NbStoptimes:      proto.Int32(req.Count),
			Count:            proto.Int32(req.Count),
			StartPage:        proto.Int32(req.StartPage),
			DisableGeojson:   proto.Bool(req.DisableGeojson),
			ItemsPerSchedule: proto.Int32(req.ItemsPerSchedule),
			RealtimeLevel:    pbnavitia.RTLevel_BASE_SCHEDULE.Enum(),
		},
		XCurrentDatetime: proto.Uint64(uint64(req.CurrentDatetime.Unix())),
	}
	pbReq.NextStopTimes.ForbiddenUri = append(pbReq.NextStopTimes.ForbiddenUri, req.ForbiddenUris...)

	return pbReq
}

func getUrl(c *gin.Context) *url.URL {
	u := location.Get(c)
	if u == nil {
		//if locatation doesn't give us an url, we use the one from gin
		return c.Request.URL
	}
	u.RawQuery = c.Request.URL.RawQuery
	u.Path = c.Request.URL.Path
	return u
}

func fillPaginationLinks(url *url.URL, response *gonavitia.RouteScheduleResponse) {
	if response == nil || response.Pagination == nil {
		return
	}
	pagination := *response.Pagination
	values := url.Query()
	if pagination.StartPage > 0 {
		values.Set("start_page", strconv.Itoa(int(pagination.StartPage-1)))
		url.RawQuery = values.Encode()
		response.Links = append(response.Links, gonavitia.Link{
			Href:      proto.String(url.String()),
			Type:      proto.String("previous"),
			Templated: proto.Bool(false),
		})
	}

	if pagination.TotalResult > (pagination.StartPage+1)*pagination.ItemsPerPage {
		values.Set("start_page", strconv.Itoa(int(pagination.StartPage+1)))
		url.RawQuery = values.Encode()
		response.Links = append(response.Links, gonavitia.Link{
			Href:      proto.String(url.String()),
			Type:      proto.String("next"),
			Templated: proto.Bool(false),
		})
	}

	if pagination.ItemsPerPage > 0 && pagination.TotalResult > 0 {
		lastPage := (pagination.TotalResult - 1) / pagination.ItemsPerPage
		values.Set("start_page", strconv.Itoa(int(lastPage)))
		url.RawQuery = values.Encode()
		response.Links = append(response.Links, gonavitia.Link{
			Href:      proto.String(url.String()),
			Type:      proto.String("last"),
			Templated: proto.Bool(false),
		})
	}
	values.Del("start_page")
	url.RawQuery = values.Encode()
	response.Links = append(response.Links, gonavitia.Link{
		Href:      proto.String(url.String()),
		Type:      proto.String("first"),
		Templated: proto.Bool(false),
	})
}
