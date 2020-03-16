package main

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	logutil "github.com/ruzulinjun/logutil"
)

func main() {
	log := make(map[string]interface{})
	ctx := context.WithValue(context.Background(), logutil.LogContextKey("log"), log)
	err := func1(ctx)
	slog := logutil.Fields{
		"trace": logutil.Trace(err),
		"err":   err.Error(),
		"log":   ctx.Value(logutil.LogContextKey("log")),
	}
	fmt.Println(slog)
}

func func1(ctx context.Context) error {
	logutil.AddLog(ctx, "func1", "enter")
	err := func2(ctx)
	if err != nil {
		err = errors.Wrap(err, logutil.MarshalFields(logutil.Fields{"value": "wrap"}))
	}
	return err
}

func func2(ctx context.Context) error {
	logutil.AddLog(ctx, "func2", "enter")
	return errors.New("this is an error")
}
