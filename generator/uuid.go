package generator

import "github.com/google/uuid"

type GoogleUUID interface {
	New() uuid.UUID
}

type GoogleUUIDCtx struct{}

func ProvideGoogleUUID() GoogleUUID {
	return &GoogleUUIDCtx{}
}

func (guc *GoogleUUIDCtx) New() uuid.UUID {
	return uuid.New()
}
