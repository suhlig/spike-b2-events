package hmac

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
)

type Authenticator struct {
	sharedKey string
}

func NewAuthenticatorMiddleware(sharedKey string) (*Authenticator, error) {
	if sharedKey == "" {
		return nil, errors.New("shared key must not be empty")
	}

	return &Authenticator{sharedKey: sharedKey}, nil
}

func (hma *Authenticator) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		signatures, isSigned := c.Request().Header["X-Bz-Event-Notification-Signature"]

		if !isSigned || len(signatures) == 0 || signatures[0] == "" {
			return fmt.Errorf("request from %s (UA '%s') does not contain an HMAC signature", c.Request().Header["X-Forwarded-For"], c.Request().Header["User-Agent"])
		}

		signature := signatures[0]
		c.Logger().Debugf("One signature present: %v", signature)
		// TODO verify HMAC

		return next(c)
	}
}
