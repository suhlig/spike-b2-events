package b2

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type EventNotification struct {
	Events []Events `json:"events"`
}

type Events struct {
	AccountID       string    `json:"accountId" validate:"required"`
	BucketID        string    `json:"bucketId" validate:"required"`
	BucketName      string    `json:"bucketName" validate:"required"`
	EventTimestamp  Timestamp `json:"eventTimestamp" validate:"required"`
	EventType       string    `json:"eventType" validate:"required"`
	EventVersion    int       `json:"eventVersion" validate:"required"`
	MatchedRuleName string    `json:"matchedRuleName" validate:"required"`
	ObjectName      string    `json:"objectName" validate:"required"`
	ObjectSize      int       `json:"objectSize" validate:"required"`
	ObjectVersionID string    `json:"objectVersionId" validate:"required"`
}

type EventNotificationValidator struct {
	Validator *validator.Validate
}

func (v *EventNotificationValidator) Validate(i interface{}) error {
	if err := v.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
