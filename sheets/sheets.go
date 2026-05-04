package sheets

import (
	"context"
	"log"
	"time"

	_ "embed"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var srv *sheets.Service

//go:embed auto-sheets.json
var bytes []byte

func init() {
	ctx := context.Background()

	var err error

	srv, err = sheets.NewService(ctx,
		option.WithAuthCredentialsJSON(option.ServiceAccount, bytes),
		option.WithScopes(sheets.SpreadsheetsScope),
	)
	if err != nil {
		log.Fatal(err)
	}
}

// ReadRange reads values from a range
func ReadRange(ctx context.Context, spreadsheetID, rng string) ([][]any, error) {
	resp, err := srv.Spreadsheets.Values.Get(
		spreadsheetID,
		rng,
	).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return resp.Values, nil
}

// WriteRange writes values to a range (overwrite)
func WriteRange(ctx context.Context, spreadsheetID, rng string, values [][]any) error {
	return retry(func() error {
		_, err := srv.Spreadsheets.Values.Update(
			spreadsheetID,
			rng,
			&sheets.ValueRange{Values: values},
		).ValueInputOption("RAW").Context(ctx).Do()
		return err
	})
}

// Append appends rows to a sheet
func Append(ctx context.Context, spreadsheetID, rng string, values [][]any) error {
	return retry(func() error {
		_, err := srv.Spreadsheets.Values.Append(
			spreadsheetID,
			rng,
			&sheets.ValueRange{Values: values},
		).ValueInputOption("RAW").Context(ctx).Do()
		return err
	})
}

// BatchWrite performs multiple updates in one request (better performance)
func BatchWrite(ctx context.Context, spreadsheetID string, data []*sheets.ValueRange) error {
	return retry(func() error {
		_, err := srv.Spreadsheets.Values.BatchUpdate(
			spreadsheetID,
			&sheets.BatchUpdateValuesRequest{
				ValueInputOption: "RAW",
				Data:             data,
			},
		).Context(ctx).Do()
		return err
	})
}

// retry helper (basic exponential backoff)
func retry(fn func() error) error {
	var err error
	for i := 0; i < 3; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	return err
}
