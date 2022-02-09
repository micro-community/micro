package main

import (
	"fmt"
	"os"

	"github.com/micro-community/micro/v3/cmd/protoc-gen-openapi/converter"
	"github.com/micro-community/micro/v3/service/logger"
	"google.golang.org/protobuf/proto"
	plugin "google.golang.org/protobuf/types/pluginpb"
)

func main() {

	// Get a converter:
	protoConverter := converter.New()

	// Convert the generator request:
	var ok = true
	logger.Debugf("Processing code generator request")
	res, err := protoConverter.ConvertFrom(os.Stdin)
	if err != nil {
		ok = false
		if res == nil {
			message := fmt.Sprintf("Failed to read input: %v", err)
			res = &plugin.CodeGeneratorResponse{
				Error: &message,
			}
		}
	}

	logger.Debug("Serializing code generator response")
	data, err := proto.Marshal(res)
	if err != nil {
		logger.Fatalf("Cannot marshal response: %v", err)
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		logger.Fatalf("Failed to write response: %v", err)
	}

	if ok {
		logger.Debug("Succeeded to process code generator request")
	} else {
		logger.Warn("Failed to process code generator but successfully sent the error to protoc")
		os.Exit(1)
	}
}
