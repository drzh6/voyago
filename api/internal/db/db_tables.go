package db

import (
	"github.com/google/uuid"
	"time"
)

type TripData struct {
	id          uuid.UUID
	name        string
	description string
	ownerId     uuid.UUID
	startTime   time.Time
	endTime     time.Time
	status      string
	idPublic    bool
	inviteCode  string
	coverImage  string
	createdAt   time.Time
	updatedAt   time.Time
}
