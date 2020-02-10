package logutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type LogContextKey string

func AddLog(ctx context.Context, key string, value interface{}) {
	log, ok := ctx.Value(LogContextKey("log")).(map[string]interface{})
	if ok {
		log[key] = value
	}
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func SprintError(err error) []string {
	type causer interface {
		Cause() error
	}
	for err != nil {
		if stack, ok := err.(stackTracer); ok {
			results := make([]string, 0)
			for _, f := range stack.StackTrace() {
				v := fmt.Sprintf("%+v", f)
				v = v[(strings.Index(v, "\n\t") + 2):]
				paths := strings.Split(v, "/")
				begin := len(paths) - 4
				if begin < 0 {
					begin = 0
				}
				v = strings.Join(paths[begin:], "/")
				results = append(results, v)
			}
			return results
		}
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return nil
}

type Fields map[string]interface{}

func MarshalFields(fields Fields) string {
	_, file, line, _ := runtime.Caller(1)
	paths := strings.Split(file, "/")
	begin := len(paths) - 4
	if begin < 0 {
		begin = 0
	}
	file = strings.Join(paths[begin:], "/")
	fields["file"] = file + ":" + strconv.Itoa(line)
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "")
	err := encoder.Encode(fields)
	if err != nil {
		return "MarshalFieldsErr: " + err.Error()
	}
	data := buffer.String()
	return data[:len(data)-1]
}
