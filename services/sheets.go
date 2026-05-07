package services

import (
	"context"
	"time"

	_ "embed"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

//go:embed auto-sheets.json
var bytes []byte

type Sheet struct {
	Srv     *sheets.Service
	SheetID string
}

func NewSheet(ctx context.Context, sheetID string) (sheet *Sheet, err error) {

	srv, err := sheets.NewService(ctx,
		option.WithAuthCredentialsJSON(option.ServiceAccount, bytes),
		option.WithScopes(sheets.SpreadsheetsScope),
	)
	if err != nil {
		return
	}

	return &Sheet{
		Srv:     srv,
		SheetID: sheetID,
	}, nil
}

func (sh *Sheet) ReadRange(ctx context.Context, rng string) ([][]any, error) {
	resp, err := sh.Srv.Spreadsheets.Values.Get(
		sh.SheetID,
		rng,
	).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return resp.Values, nil
}

func (sh *Sheet) Clear(ctx context.Context, rng string) error {
	return retry(func() error {
		_, err := sh.Srv.Spreadsheets.Values.Clear(
			sh.SheetID,
			rng,
			&sheets.ClearValuesRequest{},
		).Context(ctx).Do()
		return err
	})
}

func (sh *Sheet) UpdateRange(ctx context.Context, rng string, values [][]any) error {
	return retry(func() error {
		_, err := sh.Srv.Spreadsheets.Values.Update(
			sh.SheetID,
			rng,
			&sheets.ValueRange{Values: values},
		).ValueInputOption("RAW").Context(ctx).Do()
		return err
	})
}

func (sh *Sheet) ClearAndUpdateRange(ctx context.Context, rng string, values [][]any) error {
	err := sh.Clear(ctx, rng)
	if err != nil {
		return err
	}
	return sh.UpdateRange(ctx, rng, values)
}

func (sh *Sheet) Append(ctx context.Context, rng string, values [][]any) error {
	return retry(func() error {
		_, err := sh.Srv.Spreadsheets.Values.Append(
			sh.SheetID,
			rng,
			&sheets.ValueRange{Values: values},
		).ValueInputOption("RAW").Context(ctx).Do()
		return err
	})
}

func (sh *Sheet) BatchWrite(ctx context.Context, data []*sheets.ValueRange) error {
	return retry(func() error {
		_, err := sh.Srv.Spreadsheets.Values.BatchUpdate(
			sh.SheetID,
			&sheets.BatchUpdateValuesRequest{
				ValueInputOption: "RAW",
				Data:             data,
			},
		).Context(ctx).Do()
		return err
	})
}

func retry(fn func() error) error {
	var err error
	for i := range 3 {
		err = fn()
		if err == nil {
			return nil
		}
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	return err
}
