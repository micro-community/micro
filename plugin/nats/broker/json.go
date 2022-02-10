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
// Original source: github.com/micro/go-micro/v3/codec/json/marshaler.go

package nats

import (
	"encoding/json"

	jsonpb "google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Marshaler struct{}

func (j Marshaler) Marshal(v interface{}) ([]byte, error) {
	if pb, ok := v.(proto.Message); ok {
		return jsonpb.Marshal(pb)
	}
	return json.Marshal(v)
}

func (j Marshaler) Unmarshal(d []byte, v interface{}) error {
	if pb, ok := v.(proto.Message); ok {
		return jsonpb.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(d, pb)
	}
	return json.Unmarshal(d, v)
}

func (j Marshaler) String() string {
	return "json"
}
