package debug

import (
	"github.com/micro-community/micro/v3/platform/debug/log"
	memLog "github.com/micro-community/micro/v3/platform/debug/log/memory"
	"github.com/micro-community/micro/v3/platform/debug/profile"
	"github.com/micro-community/micro/v3/platform/debug/stats"
	memStats "github.com/micro-community/micro/v3/platform/debug/stats/memory"
	"github.com/micro-community/micro/v3/platform/debug/trace"
	memTrace "github.com/micro-community/micro/v3/platform/debug/trace/memory"
)

var (
	DefaultLog      log.Log         = memLog.NewLog()
	DefaultTracer   trace.Tracer    = memTrace.NewTracer()
	DefaultStats    stats.Stats     = memStats.NewStats()
	DefaultProfiler profile.Profile = nil
)
