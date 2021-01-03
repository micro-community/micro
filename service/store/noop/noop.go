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
// Original source: github.com/micro/go-micro/v3/store/noop.go

package noop

import (
	"github.com/micro-community/micro/v3/service/store"
)

type noopStore struct{}

func (n *noopStore) Init(opts ...store.Option) error {
	return nil
}

func (n *noopStore) Options() store.Options {
	return store.Options{}
}

func (n *noopStore) String() string {
	return "noop"
}

func (n *noopStore) Read(key string, opts ...store.ReadOption) ([]*store.Record, error) {
	return []*store.Record{}, nil
}

func (n *noopStore) Write(r *store.Record, opts ...store.WriteOption) error {
	return nil
}

func (n *noopStore) Delete(key string, opts ...store.DeleteOption) error {
	return nil
}

func (n *noopStore) List(opts ...store.ListOption) ([]string, error) {
	return []string{}, nil
}

func (n *noopStore) Close() error {
	return nil
}
