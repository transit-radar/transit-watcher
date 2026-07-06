package ebms

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

const (
	DateTimeLayout      = "02/01/2006"
	OperationTimeLayout = "15:04"
)

type RouteListItem struct {
	ID     json.Number `json:"RouteId"`
	Number string      `json:"RouteNo" conform:"trim"`
	Name   string      `json:"RouteName" conform:"trim"`
}

type Route struct {
	ID                  json.Number   `json:"RouteId"`
	Number              string        `json:"RouteNo" conform:"trim"`
	Name                string        `json:"RouteName" conform:"trim"`
	Color               string        `json:"Color" conform:"trim"`
	Type                string        `json:"Type" conform:"trim"`
	Distance            json.Number   `json:"Distance"`
	Orgs                string        `json:"Orgs" conform:"trim"`
	TimeOfTrip          string        `json:"TimeOfTrip" conform:"trim"`
	Headway             string        `json:"Headway" conform:"trim"`
	OperationTime       OperationTime `json:"OperationTime" conform:"trim"`
	NumOfSeats          string        `json:"NumOfSeats" conform:"trim"`
	OutBoundName        string        `json:"OutBoundName" conform:"trim"`
	InBoundName         string        `json:"InBoundName" conform:"trim"`
	OutBoundDescription string        `json:"OutBoundDescription" conform:"trim"`
	InBoundDescription  string        `json:"InBoundDescription" conform:"trim"`
	TotalTrip           string        `json:"TotalTrip" conform:"trim"`
	Tickets             string        `json:"Tickets" conform:"trim"`
}

type Variant struct {
	ID               json.Number `json:"RouteVarId"`
	RouteID          json.Number `json:"RouteId"`
	VariantName      string      `json:"RouteVarName" conform:"trim"`
	VariantShortName string      `json:"RouteVarShortName" conform:"trim"`
	RouteNumber      string      `json:"RouteNo" conform:"trim"`
	StartStop        string      `json:"StartStop" conform:"trim"`
	EndStop          string      `json:"EndStop" conform:"trim"`
	Distance         json.Number `json:"Distance"`
	Outbound         bool        `json:"Outbound"`
	RunningTime      int         `json:"RunningTime"`
}

type Stop struct {
	ID                json.Number `json:"StopId"`
	Code              string      `json:"Code" conform:"trim"`
	Name              string      `json:"Name" conform:"trim"`
	StopType          string      `json:"StopType" conform:"trim"`
	Zone              string      `json:"Zone" conform:"trim"`
	Ward              *string     `json:"Ward" conform:"trim"`
	AddressNo         string      `json:"AddressNo" conform:"trim"`
	Street            string      `json:"Street" conform:"trim"`
	SupportDisability string      `json:"SupportDisability" conform:"trim"`
	Status            string      `json:"Status" conform:"trim"`
	Longitude         json.Number `json:"Lng"`
	Latitude          json.Number `json:"Lat"`
	Search            string      `json:"Search" conform:"trim"`
	Routes            string      `json:"Routes" conform:"trim"`
}

type Path struct {
	Latitude  []json.Number `json:"lat"`
	Longitude []json.Number `json:"lng"`
}

type Timetable struct {
	RouteID          json.Number   `json:"RouteId"`
	VariantID        json.Number   `json:"RouteVarId"`
	TimetableID      json.Number   `json:"TimeTableId"`
	VariantShortName string        `json:"RouteVarShortName" conform:"trim"`
	StartDate        Date          `json:"StartDate" conform:"trim"`
	EndDate          Date          `json:"EndDate" conform:"trim"`
	StartStop        string        `json:"StartStop" conform:"trim"`
	EndStop          string        `json:"EndStop" conform:"trim"`
	IsCurrent        bool          `json:"IsCurrent"`
	RunningTime      string        `json:"RunningTime" conform:"trim"`
	Headway          string        `json:"Headway" conform:"trim"`
	TotalTrip        int           `json:"TotalTrip"`
	OperationTime    OperationTime `json:"OperationTime" conform:"trim"`
	ApplyDates       AppliedDates  `json:"ApplyDates" conform:"trim"`
}

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if s == "" {
		return nil
	}

	var err error
	d.Time, err = time.Parse(DateTimeLayout, s)

	if err != nil {
		return err
	}

	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(d.Format(DateTimeLayout)), nil
}

type AppliedDate struct {
	time.Weekday
}

func (ad *AppliedDate) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch s {
	case "T2":
		ad.Weekday = time.Monday
	case "T3":
		ad.Weekday = time.Tuesday
	case "T4":
		ad.Weekday = time.Wednesday
	case "T5":
		ad.Weekday = time.Thursday
	case "T6":
		ad.Weekday = time.Friday
	case "T7":
		ad.Weekday = time.Saturday
	case "CN":
		ad.Weekday = time.Sunday
	default:
		return errors.New("invalid weekday")
	}

	return nil
}

func (ad AppliedDate) MarshalJSON() ([]byte, error) {
	switch ad.Weekday {
	case time.Monday:
		return []byte("T2"), nil
	case time.Tuesday:
		return []byte("T3"), nil
	case time.Wednesday:
		return []byte("T4"), nil
	case time.Thursday:
		return []byte("T5"), nil
	case time.Friday:
		return []byte("T6"), nil
	case time.Saturday:
		return []byte("T7"), nil
	case time.Sunday:
		return []byte("CN"), nil
	default:
		return nil, errors.New("invalid weekday")
	}
}

type AppliedDates struct {
	Dates []AppliedDate
}

func (ad *AppliedDates) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	for part := range strings.SplitSeq(s, ",") {
		var appliedDate AppliedDate
		if err := json.Unmarshal([]byte(`"`+strings.TrimSpace(part)+`"`), &appliedDate); err != nil {
			return err
		}

		ad.Dates = append(ad.Dates, appliedDate)
	}

	return nil
}

func (ad AppliedDates) MarshalJSON() ([]byte, error) {
	var dates []string
	for _, date := range ad.Dates {
		b, err := json.Marshal(date)
		if err != nil {
			return nil, err
		}
		dates = append(dates, string(b))
	}

	return []byte(strings.Join(dates, ", ")), nil
}

type Trip struct {
	RouteId     int    `json:"RouteId"`
	TripId      int    `json:"TripId"`
	TimeTableId int    `json:"TimeTableId"`
	StartTime   string `json:"StartTime" conform:"trim"`
	EndTime     string `json:"EndTime" conform:"trim"`
}

type OperationTime struct {
	From, To time.Time
}

func (op *OperationTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	parts := strings.Split(s, " - ")

	if len(parts) != 2 {
		return errors.New("invalid time span")
	}

	var err error

	op.From, err = time.Parse(OperationTimeLayout, parts[0])
	if err != nil {
		return err
	}

	op.To, err = time.Parse(OperationTimeLayout, parts[1])
	if err != nil {
		return err
	}

	return nil
}

func (op OperationTime) MarshalJSON() ([]byte, error) {
	from := op.From.Format(OperationTimeLayout)
	to := op.To.Format(OperationTimeLayout)

	return json.Marshal(from + " - " + to)
}
