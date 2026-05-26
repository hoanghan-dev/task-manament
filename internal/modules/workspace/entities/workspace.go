package workspace

import (
	"time"

	"github.com/google/uuid"
)

type Workspace struct {
	WorkspaceId uuid.UUID
	Name        string
	Description string
	OwnerId     uuid.UUID
	CreateAt    time.Time
}
