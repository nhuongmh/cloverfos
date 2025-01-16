package energy

import (
	"context"
	"encoding/base64"

	"github.com/nhuongmh/cvfs/timemn/internal/logger"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	FORM_MINNA_ID_COLUMN    = 0
	FORM_FORMULA_COLUMN     = 1
	FORM_BACKWARD_COLUMN    = 2
	FORM_DESCRIPTION_COLUMN = 3
)

func InitNewGoogleSheetService(googleKeyBase64 string) (*sheets.Service, error) {
	logger.Log.Info().Msg("Initializing google sheet service")
	// create api context
	ctx := context.Background()

	// get bytes from base64 encoded google service accounts key
	credBytes, err := base64.StdEncoding.DecodeString(googleKeyBase64)
	if err != nil {
		return nil, errors.Wrap(err, "Failed decode base64 google key")
	}
	// logger.Log.Debug().Msgf("API: %v", string(credBytes))

	// authenticate and get configuration
	config, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to authenticate to google API")
	}

	// create client with config and context
	client := config.Client(ctx)

	// create new service using client
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, errors.Wrap(err, "Failed creating new google service")
	}
	return srv, nil
}
