package models

import (
	entities "mvp-2-spms/domain-aggregate"
)

type PlannerData struct {
	Id   string
	Name string
}

type PlannerMeeting struct {
	entities.Meeting
	MeetingPlannerId string
}

type PlannerSlot struct {
	entities.Slot
	MeetingPlannerId string
}
