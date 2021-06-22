package main

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	logutil "github.com/ruzulinjun/logutil"
)

func main() {
	log := logutil.Fields{}
	ctx := context.WithValue(context.Background(), logutil.LogContextKey("log"), log)
	err := func1(ctx)
	slog := logutil.Fields{
		"trace": logutil.Trace(err),
		"err":   err.Error(),
		"log":   ctx.Value(logutil.LogContextKey("log")),
	}
	fmt.Println(slog)
}

/* print log
map[err:{"add_message_to_err":"call_func2_gen_err","file":"ruzulinjun/logutil/example/example.go:28","name":"main.func1"}: {"err":"an error generated here","file":"ruzulinjun/logutil/example/example.go:35","name":"main.func2"} log:map[func1:enter func2:enter] trace:[ruzulinjun/logutil/example/example.go:35 ruzulinjun/logutil/example/example.go:25 ruzulinjun/logutil/example/example.go:14 go/src/runtime/proc.go:225 go/src/runtime/asm_amd64.s:1371]]
*/

func func1(ctx context.Context) error {
	logutil.AddLog(ctx, "func1", "enter")
	err := func2(ctx)
	if err != nil {
		//err = errors.Wrap(err, logutil.MarshalFields(logutil.Fields{"add_message_to_err": "call_func2_gen_err"}))
		err = errors.WithMessage(err, logutil.MarshalFields(logutil.Fields{"add_message_to_err": "call_func2_gen_err"}))
	}
	return err
}

func func2(ctx context.Context) error {
	logutil.AddLog(ctx, "func2", "enter")
	return errors.New(logutil.MarshalFields(logutil.Fields{"err": "an error generated here"}))
}
