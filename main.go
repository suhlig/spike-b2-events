package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/suhlig/spike-b2-events/b2"
)

func main() {
	err := mainE()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %s\n", err)
		os.Exit(1)
	}
}

func mainE() error {
	e := echo.New()
	e.HideBanner = true
	e.Logger.SetLevel(log.DEBUG)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &b2.EventNotificationValidator{Validator: validator.New()}

	hmacMw, err := b2.NewHmacAuthenticator(os.Getenv("B2_EVENT_NOTIFICATIONS_SHARED_SECRET"))

	if err != nil {
		e.Logger.Warnf("%s, if you want to verify the request payload set the environment variable B2_EVENT_NOTIFICATIONS_SHARED_SECRET", err)
	} else {
		e.Use(hmacMw.Process)
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "B2 Event Notifications Server; see https://github.com/suhlig/spike-b2-events\n")
	})

	e.POST("/", func(c echo.Context) error {
		notification := new(b2.EventNotification)

		if err := c.Bind(notification); err != nil {
			e.Logger.Errorf("parsing the request failed: %w", err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := c.Validate(notification); err != nil {
			e.Logger.Errorf("validating the request failed: %w", err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		for _, event := range notification.Events {
			e.Logger.Infof("Object %s was uploaded to %s at %s", event.ObjectName, event.BucketName, event.EventTimestamp)
		}

		return c.String(http.StatusOK, fmt.Sprintf("processed %d events\n", len(notification.Events)))
	})

	return e.Start("127.0.0.1:62057")
}
