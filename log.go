package logutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type LogContextKey string

var mutex sync.Mutex

/*
...
	logutil.AddLog(ctx, "validStatusChange", "n")
...
	log.Print(ctx.Value(logutil.LogContextKey("log")))
...
*/

func AddLog(ctx context.Context, key string, value interface{}) {
	log, ok := ctx.Value(LogContextKey("log")).(Fields)
	if ok {
		mutex.Lock()
		log[key] = value
		mutex.Unlock()
	}
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// get error's trace
func Trace(err error) []string {
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

/*
add file name and function name automatically
...
	if err != nil {
		return errors.Wrap(err, logutil.MarshalFields(logutil.Fields{"userID": userID}))
	}
...
*/
func MarshalFields(fields Fields) string {
	pcs := make([]uintptr, 1)
	runtime.Callers(2, pcs[:])
	frame, _ := runtime.CallersFrames(pcs).Next()
	file := frame.File
	paths := strings.Split(file, "/")
	begin := len(paths) - 4
	if begin < 0 {
		begin = 0
	}
	fields["file"] = strings.Join(paths[begin:], "/") + ":" + strconv.Itoa(frame.Line)
	name := frame.Function
	paths = strings.Split(name, "/")
	begin = len(paths) - 5
	if begin < 0 {
		begin = 0
	}
	fields["name"] = strings.Join(paths[begin:], "/")

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
