package gowaze

import "encoding/json"

type Route struct {
	Alternatives []Alternatives `json:"alternatives"`
}

type Path struct {
	SegmentID int     `json:"segmentId"`
	NodeID    int     `json:"nodeId"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Direction bool    `json:"direction"`
}

type Instruction struct {
	Opcode          string      `json:"opcode"`
	Arg             int         `json:"arg"`
	InstructionText interface{} `json:"instructionText"`
	LaneGuidance    interface{} `json:"laneGuidance"`
	Name            interface{} `json:"name"`
	Tts             interface{} `json:"tts"`
}

type Results struct {
	Path                     Path          `json:"path"`
	Street                   int           `json:"street"`
	AltStreets               interface{}   `json:"altStreets"`
	Distance                 int           `json:"distance"`
	Length                   int           `json:"length"`
	CrossTime                int           `json:"crossTime"`
	CrossTimeWithoutRealTime int           `json:"crossTimeWithoutRealTime"`
	Tiles                    interface{}   `json:"tiles"`
	ClientIds                interface{}   `json:"clientIds"`
	KnownDirection           bool          `json:"knownDirection"`
	Penalty                  int           `json:"penalty"`
	RoadType                 int           `json:"roadType"`
	IsToll                   bool          `json:"isToll"`
	NaiveRoute               interface{}   `json:"naiveRoute"`
	DetourSavings            int           `json:"detourSavings"`
	DetourSavingsNoRT        int           `json:"detourSavingsNoRT"`
	UseHovLane               bool          `json:"useHovLane"`
	Attributes               int           `json:"attributes"`
	Lane                     string        `json:"lane"`
	LaneType                 interface{}   `json:"laneType"`
	Areas                    []interface{} `json:"areas"`
	RequiredPermits          []interface{} `json:"requiredPermits"`
	DetourRoute              interface{}   `json:"detourRoute"`
	NaiveRouteFullResult     interface{}   `json:"naiveRouteFullResult"`
	DetourRouteFullResult    interface{}   `json:"detourRouteFullResult"`
	MergeOffset              int           `json:"mergeOffset"`
	AvoidStatus              string        `json:"avoidStatus"`
	ClientLaneSet            interface{}   `json:"clientLaneSet"`
	AdditionalInstruction    interface{}   `json:"additionalInstruction"`
	Instruction              Instruction   `json:"instruction"`
}

type Response struct {
	Results                 []Results     `json:"results"`
	StreetNames             []interface{} `json:"streetNames"`
	TileIds                 []interface{} `json:"tileIds"`
	TileUpdateTimes         []interface{} `json:"tileUpdateTimes"`
	Geom                    interface{}   `json:"geom"`
	FromFraction            float64       `json:"fromFraction"`
	ToFraction              float64       `json:"toFraction"`
	SameFromSegment         bool          `json:"sameFromSegment"`
	SameToSegment           bool          `json:"sameToSegment"`
	AstarPoints             interface{}   `json:"astarPoints"`
	WayPointIndexes         interface{}   `json:"wayPointIndexes"`
	WayPointFractions       interface{}   `json:"wayPointFractions"`
	TollMeters              int           `json:"tollMeters"`
	PreferedRouteID         int           `json:"preferedRouteId"`
	IsInvalid               bool          `json:"isInvalid"`
	IsBlocked               bool          `json:"isBlocked"`
	ServerUniqueID          string        `json:"serverUniqueId"`
	DisplayRoute            bool          `json:"displayRoute"`
	AstarVisited            int           `json:"astarVisited"`
	AstarResult             string        `json:"astarResult"`
	AstarData               interface{}   `json:"astarData"`
	IsRestricted            bool          `json:"isRestricted"`
	AvoidStatus             string        `json:"avoidStatus"`
	DueToOverride           interface{}   `json:"dueToOverride"`
	PassesThroughDangerArea bool          `json:"passesThroughDangerArea"`
	DistanceFromSource      int           `json:"distanceFromSource"`
	DistanceFromTarget      int           `json:"distanceFromTarget"`
	MinPassengers           int           `json:"minPassengers"`
	HovIndex                int           `json:"hovIndex"`
	TimeZone                interface{}   `json:"timeZone"`
	RouteType               []string      `json:"routeType"`
	RouteAttr               []interface{} `json:"routeAttr"`
	AstarCost               int           `json:"astarCost"`
	ReorderChoice           interface{}   `json:"reorderChoice"`
	TotalRouteTime          int           `json:"totalRouteTime"`
	LaneTypes               []interface{} `json:"laneTypes"`
	PreferredStoppingPoints interface{}   `json:"preferredStoppingPoints"`
	Areas                   []interface{} `json:"areas"`
	RequiredPermits         []interface{} `json:"requiredPermits"`
	EtaHistograms           []interface{} `json:"etaHistograms"`
	EntryPoint              interface{}   `json:"entryPoint"`
	ShortRouteName          string        `json:"shortRouteName"`
	TollPrice               interface{}   `json:"tollPrice"`
	SegGeoms                interface{}   `json:"segGeoms"`
	RouteName               string        `json:"routeName"`
	RouteNameStreetIds      []int         `json:"routeNameStreetIds"`
	Open                    bool          `json:"open"`
}

type Coords struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z string  `json:"z"`
}

type Alternatives struct {
	Response  Response    `json:"response"`
	Coords    []Coords    `json:"coords"`
	SegCoords interface{} `json:"segCoords"`
}

func (h *Handler) GetRoute(fromX, fromY, toX, toY float64) (Route, error) {
	var i Route
	data, err := h.Get("GET", urlGetRoute, apiURL, fromX, fromY, toX, toY)

	err = json.Unmarshal(data, &i)
	if err != nil {
		return Route{}, err
	}

	return i, nil
}

func (r *Route) TravelTimes() []int {
	var times []int
	for _, a := range r.Alternatives {
		times = append(times, a.Response.TotalRouteTime)
	}
	return times
}

func (r *Route) DistanceFromTarget() []int {
	var times []int
	for _, a := range r.Alternatives {
		times = append(times, a.Response.DistanceFromTarget)
	}
	return times
}

func (r *Route) Distance() []int {
	var times []int
	for _, a := range r.Alternatives {
		t := 0
		for _, r := range a.Response.Results {
			t += r.Length
		}
		times = append(times, t)
	}
	return times
}

func (r *Response) Distance() int {
	t := 0
	for _, r := range r.Results {
		t += r.Length
	}
	return t
}
