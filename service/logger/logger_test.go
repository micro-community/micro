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
// Original source: github.com/micro/go-micro/v3/logger/logger_test.go

package logger

import (
	"bufio"
	"strings"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	l := NewLogger(WithLevel(TraceLevel), WithFormat(JsonFormat))
	h1 := NewHelper(l).WithFields(map[string]interface{}{"key1": "val1"})
	t1 := time.Now()

	for i := 0; i < 10000; i++ {
		h1.Info("hello")
	}

	h1.Infof("%v", time.Since(t1).String())
	h1.Trace("trace_msg1")
	h1.Warn("warn_msg1")

	h2 := NewHelper(l).WithFields(map[string]interface{}{"key2": "val2"})
	h2.Trace("trace_msg2")
	h2.Warn("warn_msg2")

	l.Fields(map[string]interface{}{"key3": "val4"}).Log(InfoLevel, "test_msg")
}

func TestLoggerRedirection(t *testing.T) {
	b := strings.Builder{}
	wr := bufio.NewWriter(&b)
	NewLogger(WithOutput(wr)).Logf(InfoLevel, "test message")
	wr.Flush()
	if !strings.Contains(b.String(), "level=info msg=test message") {
		t.Fatalf("Redirection failed, received '%s'", b.String())
	}
}
