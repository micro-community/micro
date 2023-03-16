// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/go-micro/v3/logger/default.go

package logger

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/bytedance/sonic"
	"golang.org/x/exp/slog"

	dlog "github.com/micro-community/micro/v3/service/debug/log"
)

var DefaultLogger Logger

func init() {
	lvl, err := GetLevel(os.Getenv("MICRO_LOG_LEVEL"))
	if err != nil {
		lvl = InfoLevel
	}

	SetDefault(NewHelper(NewLogger(WithLevel(lvl))))
}

var defaultLoggerCache atomic.Value

// GetDefault returns the default Logger.
func GetDefault() Logger {
	return defaultLoggerCache.Load().(*Helper)
}

// SetDefault sets the default Logger and returns default.
func SetDefault(l Logger) *Helper {
	defaultLoggerCache.Store(l)
	DefaultLogger = GetDefault()
	return DefaultLogger.(*Helper)
}

type defaultLogger struct {
	sync.RWMutex
	opts Options
}

// Init(opts...) should only overwrite provided options
func (l *defaultLogger) Init(opts ...Option) error {
	for _, o := range opts {
		o(&l.opts)
	}
	return nil
}

func (l *defaultLogger) String() string {
	return "default"
}

func (l *defaultLogger) Fields(fields map[string]interface{}) Logger {
	l.Lock()
	l.opts.Fields = copyFields(fields)
	l.Unlock()
	return l
}

func copyFields(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (l *defaultLogger) Log(level Level, v ...interface{}) {
	// TODO decide does we need to write message if log level not used?
	if !l.opts.Level.Enabled(level) {
		return
	}

	record := slog.Record{}
	record.Time = time.Now().Local()
	record.Level = slog.Level(l.opts.Level * 4)

	for k, v := range l.opts.Fields {
		record.AddAttrs(slog.Any(k, v))
	}

	record.AddAttrs(slog.String("level", level.String()))
	record.AddAttrs(slog.String("file", getLineByRuntimeCache(l.opts.CallerSkipCount)))

	switch {
	case len(v) == 2:
		if object, ok := v[1].([]byte); ok {
			attrs := make(map[string]interface{})
			_ = sonic.Unmarshal(object, &attrs)
			for k, v := range attrs {
				record.AddAttrs(slog.Any(k, v))
			}
		} else {
			refType := reflect.TypeOf(v[1])
			refValue := reflect.ValueOf(v[1])
			for i := 0; i < refType.NumField(); i++ {
				record.AddAttrs(slog.Any(refType.Field(i).Name, refValue.Field(i).Interface()))
			}
		}

	case len(v) > 2:
		for i := 1; i < len(v); {
			if i+2 > len(v) {
				break
			}
			key, val := v[i], v[i+1]
			if k, ok := key.(string); ok {
				record.AddAttrs(slog.Any(k, val))
			}
			i += 2
		}
	default:
	}

	if len(v) > 0 {
		record.Message, _ = v[0].(string)
	}

	if l.opts.Handler == nil {
		l.opts.Handler = NewHandler(&l.opts)
	}
	l.opts.Handler.Handle(l.opts.Context, record)

}

func (l *defaultLogger) Logf(level Level, format string, v ...interface{}) {
	//	 TODO decide does we need to write message if log level not used?
	if level < l.opts.Level {
		return
	}

	l.RLock()
	fields := copyFields(l.opts.Fields)
	l.RUnlock()

	fields["level"] = level.String()

	if _, file, line, ok := runtime.Caller(l.opts.CallerSkipCount); ok {
		fields["file"] = fmt.Sprintf("%s:%d", logCallerfilePath(file), line)
	}

	rec := dlog.Record{
		Timestamp: time.Now().Local(),
		Message:   strings.ReplaceAll(fmt.Sprintf(format, v...), "\n", ""),
		Metadata:  make(map[string]string, len(fields)),
	}

	keys := make([]string, 0, len(fields))
	for k, v := range fields {
		keys = append(keys, k)
		rec.Metadata[k] = fmt.Sprintf("%v", v)
	}

	sort.Strings(keys)
	metadata := strings.Builder{}

	for _, k := range keys {
		metadata.WriteString(fmt.Sprintf(" %s=%v", k, fields[k]))
	}

	t := rec.Timestamp.Format(NanoTimeFieldFormat)
	fmt.Fprintf(l.opts.Out, "ts=%s %s msg=%v\n", t, metadata.String(), rec.Message)
}

func (l *defaultLogger) Options() Options {
	// not guard against options Context values
	l.RLock()
	opts := l.opts
	opts.Fields = copyFields(l.opts.Fields)
	l.RUnlock()
	return opts
}

var (
	mapRuntimeCache unsafe.Pointer = func() unsafe.Pointer {
		m := make(map[uintptr]string, 1024)
		return unsafe.Pointer(&m)
	}()
)

func getLineByRuntimeCache(skip int) (line string) {
	var pcs [1]uintptr
	runtime.Callers(skip+2, pcs[:])
	pc := pcs[0]
	mPCs := *(*map[uintptr]string)(atomic.LoadPointer(&mapRuntimeCache))
	line, ok := mPCs[pc]
	if !ok {
		fs := runtime.CallersFrames([]uintptr{pc})
		// TODO: error-checking?
		f, _ := fs.Next()
		line = logCallerfilePath(f.File) + ":" + strconv.Itoa(f.Line)
		mPCs2 := make(map[uintptr]string, len(mPCs)+10)
		mPCs2[pc] = line
		for {
			p := atomic.LoadPointer(&mapRuntimeCache)
			mPCs = *(*map[uintptr]string)(p)
			for k, v := range mPCs {
				mPCs2[k] = v
			}
			swapped := atomic.CompareAndSwapPointer(&mapRuntimeCache, p, unsafe.Pointer(&mPCs2))
			if swapped {
				break
			}
		}
	}
	return
}

// logCallerfilePath returns a package/file:line description of the caller,
// preserving only the leaf directory name and file name.
func logCallerfilePath(loggingFilePath string) string {
	// To make sure we trim the path correctly on Windows too, we
	// counter-intuitively need to use '/' and *not* os.PathSeparator here,
	// because the path given originates from Go stdlib, specifically
	// runtime.Caller() which (as of Mar/17) returns forward slashes even on
	// Windows.
	//
	// See https://github.com/golang/go/issues/3335
	// and https://github.com/golang/go/issues/18151
	//
	// for discussion on the issue on Go side.
	idx := strings.LastIndexByte(loggingFilePath, '/')
	if idx == -1 {
		return loggingFilePath
	}
	idx = strings.LastIndexByte(loggingFilePath[:idx], '/')
	if idx == -1 {
		return loggingFilePath
	}
	return loggingFilePath[idx+1:]
}

// NewLogger builds a new logger based on options
func NewLogger(opts ...Option) Logger {
	// Default options
	options := Options{
		Level:           InfoLevel,
		Fields:          make(map[string]interface{}),
		Out:             os.Stderr,
		CallerSkipCount: 2,
		Context:         context.Background(),
	}

	l := &defaultLogger{opts: options}
	if err := l.Init(opts...); err != nil {
		l.Log(FatalLevel, err)
	}

	return l
}
