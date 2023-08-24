// Package usage tracks micro usage
package usage

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	pb "github.com/micro-community/micro/v3/cmd/usage/proto"
	"github.com/micro-community/micro/v3/util/version"
	"google.golang.org/protobuf/proto"
)

var (
	// usage url
	u = "https://micro.dev/usage"
	// usage agent
	a = "micro/usage"
	// usage version
	v = version.V
	// 24 hour window
	w = 8.64e13
)

// New generates a new usage report to be filled in
func New(service string) *pb.Usage {
	id := fmt.Sprintf("micro.%s.%s.%s", service, version.V, uuid.New().String())
	srv := "micro." + service

	if len(service) == 0 {
		id = fmt.Sprintf("micro.%s.%s", version.V, uuid.New().String())
		srv = "micro"
	}

	sum := sha256.Sum256([]byte(id))

	return &pb.Usage{
		Service:   srv,
		Version:   v,
		Id:        fmt.Sprintf("%x", sum),
		Timestamp: uint64(time.Now().UnixNano()),
		Window:    uint64(w),
		Metrics: &pb.Metrics{
			Count: make(map[string]uint64),
		},
	}
}

// Report reports the current usage
func Report(ug *pb.Usage) error {
	if v := os.Getenv("MICRO_REPORT_USAGE"); v == "false" {
		return nil
	}

	// update timestamp/window
	now := uint64(time.Now().UnixNano())
	ug.Window = now - ug.Timestamp
	ug.Timestamp = now

	p, err := proto.Marshal(ug)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", u, bytes.NewReader(p))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/protobuf")
	req.Header.Set("User-Agent", a)
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	io.Copy(io.Discard, rsp.Body)
	return nil
}
