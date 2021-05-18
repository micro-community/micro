package debug

import (
	"github.com/micro-community/micro/v3/service/debug/log"
	memLog "github.com/micro-community/micro/v3/service/debug/log/memory"
	"github.com/micro-community/micro/v3/service/debug/profile"
	"github.com/micro-community/micro/v3/service/debug/stats"
	memStats "github.com/micro-community/micro/v3/service/debug/stats/memory"
	"github.com/micro-community/micro/v3/service/debug/trace"
	memTrace "github.com/micro-community/micro/v3/service/debug/trace/memory"
)

// export default
var (
	DefaultLog      log.Log         = memLog.NewLog()
	DefaultTracer   trace.Tracer    = memTrace.NewTracer()
	DefaultStats    stats.Stats     = memStats.NewStats()
	DefaultProfiler profile.Profile = nil
)
