package main

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	util "github.com/ruzulinjun/golang-errors"
)

func main() {
	log := make(map[string]interface{})
	ctx := context.WithValue(context.Background(), util.LogContextKey("log"), log)
	err := func1(ctx)
	slog := util.Fields{
		"trace": util.SprintError(err),
		"err":   err.Error(),
		"log":   ctx.Value(util.LogContextKey("log")),
	}
	fmt.Println(slog)
}

func func1(ctx context.Context) error {
	util.AddLog(ctx, "func1", "enter")
	err := func2(ctx)
	if err != nil {
		err = errors.Wrap(err, util.MarshalFields(util.Fields{"value": "wrap"}))
	}
	return err
}

func func2(ctx context.Context) error {
	util.AddLog(ctx, "func2", "enter")
	return errors.New("this is an error")
}
