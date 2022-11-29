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
// Original source: github.com/micro/go-micro/v3/logger/options.go

package logger

import (
	"context"
	"golang.org/x/exp/slog"
	"io"
)

type Option func(*Options)

type Options struct {
	// The logging level the logger should log at. default is `InfoLevel`
	Level Level
	// fields to always be logged
	Fields map[string]interface{}
	// It's common to set this to a file, or leave it default which is `os.Stderr`
	Out io.Writer
	// Caller skip frame count for file:line info
	CallerSkipCount int
	// Alternative options
	Context context.Context
	// Format log print format options
	Format Format
	// A Handler handles log records produced by a Logger..
	//
	// A typical handler may print log records to standard error,
	// or write them to a file or database, or perhaps augment them
	// with additional attributes and pass them on to another handler.
	//
	// Any of the Handler's methods may be called concurrently with itself
	// or with other methods. It is the responsibility of the Handler to
	// manage this concurrency.
	Handler slog.Handler
}

// WithFields set default fields for the logger
func WithFields(fields map[string]interface{}) Option {
	return func(args *Options) {
		args.Fields = fields
	}
}

// WithLevel set default level for the logger
func WithLevel(level Level) Option {
	return func(args *Options) {
		args.Level = level
	}
}

// WithOutput set default output writer for the logger
func WithOutput(out io.Writer) Option {
	return func(args *Options) {
		args.Out = out
	}
}

// WithCallerSkipCount set frame count to skip
func WithCallerSkipCount(c int) Option {
	return func(args *Options) {
		args.CallerSkipCount = c
	}
}

// WithFormat set default output format for the logger
func WithFormat(f Format) Option {
	return func(args *Options) {
		args.Format = f
	}
}

// WithContext set default context for the logger
func WithContext(k, v interface{}) Option {
	return func(args *Options) {
		if args.Context == nil {
			args.Context = context.Background()
		}
		args.Context = context.WithValue(args.Context, k, v)
	}
}

// WithHandler set default handler for the logger
func WithHandler(h slog.Handler) Option {
	return func(args *Options) {
		if h != nil {
			args.Handler = h
			return
		}
		
		args.Handler = NewHandler(args)
	}
}
