package b2

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type HmacAuthenticator struct {
	sharedKey string
}

func NewHmacAuthenticator(sharedKey string) (*HmacAuthenticator, error) {
	if sharedKey == "" {
		return nil, errors.New("shared key must not be empty")
	}

	return &HmacAuthenticator{sharedKey: sharedKey}, nil
}

func (hma *HmacAuthenticator) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method != http.MethodPost {
			c.Logger().Debugf("skipping HMAC authentication for method %s", c.Request().Method)
			return next(c)
		}

		receivedSignature, err := parseReceivedSignature(c.Request())

		if err != nil {
			c.Logger().Errorf("parsing HMAC signature header failed: %w", err)
			return echo.NewHTTPError(http.StatusBadRequest, "parsing HMAC signature header failed")
		}

		calculatedSig, err := hma.calculateSignature(c.Request())

		if err != nil {
			c.Logger().Errorf("calculating HMAC signature failed: %w", err)
			return echo.NewHTTPError(http.StatusBadRequest, "calculating HMAC signature failed")
		}

		if calculatedSig != receivedSignature {
			c.Logger().Errorf("HMAC signatures do not match: %s != %s", calculatedSig, receivedSignature)
			return echo.NewHTTPError(http.StatusBadRequest, "HMAC signatures do not match")
		}

		return next(c)
	}
}

func parseReceivedSignature(request *http.Request) (string, error) {
	signatures, isSigned := request.Header["X-Bz-Event-Notification-Signature"]

	if !isSigned || len(signatures) == 0 || signatures[0] == "" {
		return "", fmt.Errorf("request from %s (UA '%s') does not contain an HMAC signature", request.Header["X-Forwarded-For"], request.Header["User-Agent"])
	}

	parts := strings.Split(signatures[0], "=")

	if len(parts) != 2 {
		return "", errors.New("invalid signature format")
	}

	if parts[0] != "v1" {
		return "", errors.New("invalid signature version")
	}

	return parts[1], nil
}

func (hma *HmacAuthenticator) calculateSignature(request *http.Request) (string, error) {
	body, err := io.ReadAll(request.Body)

	if err != nil {
		return "", fmt.Errorf("reading request body: %w", err)
	}

	request.Body = io.NopCloser(bytes.NewBuffer(body)) // restore body for the next handler

	hash := hmac.New(sha256.New, []byte(hma.sharedKey))

	_, err = hash.Write(body)

	if err != nil {
		return "", fmt.Errorf("writing to hash: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
