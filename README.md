# Micro [![License](https://img.shields.io/badge/license-apache-blue)](https://opensource.org/licenses/Apache-2.0) [![Go Report Card](https://goreportcard.com/badge/micro-community/micro)](https://goreportcard.com/report/github.com/micro-community/micro)

<kbd><img src="https://raw.githubusercontent.com/micro-community/micro/master/docs/images/banner.png" /></kbd>
Micro is a cloud platform for API development.

## Overview

**A community fork and extension of [micro](https://github.com/micro/micro) with great hornor.**

Micro addresses the key requirements for building services in the cloud. It leverages the microservices architecture
pattern and provides a set of services which act as the building blocks of a platform. Micro deals with the complexity
of distributed systems and provides simpler programmable abstractions to build on.

## Contents

- [Introduction](https://micro.arch.wiki/introduction) - A high level introduction to Micro
- [Getting Started](https://micro.arch.wiki/getting-started) - The hello-world quick-start guide
- [Upgrade Guide](https://micro.arch.wiki/upgrade-guide) - Update your go-micro project to use micro v3.
- [Architecture](https://micro.arch.wiki/architecture) - Describes the architecture, design and tradeoffs
- [Reference](https://micro.arch.wiki/reference) - In-depth reference for Micro CLI and services
- [Resources](https://micro.arch.wiki/resources) - External resources and contributions
- [Roadmap](https://micro.arch.wiki/roadmap) - Stuff on our agenda over the long haul
- [Users](https://micro.arch.wiki/users) - Developers and companies using Micro in production
- [FAQ](https://micro.arch.wiki/faq) - Frequently asked questions


## Getting Started

Find the cloud hosted services at [m3o.com](https://m3o.com)
Below are the core components that make up Micro.

**Server**

Micro is built as a microservices architecture and abstracts away the complexity of the underlying infrastructure. We compose 
this as a single logical server to the user but decompose that into the various building block primitives that can be plugged 
into any underlying system. 

The server is composed of the following services.

- **API** - HTTP Gateway which dynamically maps http/json requests to RPC using path based resolution
- **Auth** - Authentication and authorization out of the box using jwt tokens and rule based access control.
- **Broker** - Ephemeral pubsub messaging for asynchronous communication and distributing notifications
- **Config** - Dynamic configuration and secrets management for service level config without the need to restart
- **Events** - Event streaming with ordered messaging, replay from offsets and persistent storage
- **Network** - Inter-service networking, isolation and routing plane for all internal request traffic
- **Proxy** - An identity aware proxy used for remote access and any external grpc request traffic
- **Runtime** - Service lifecycle and process management with support for source to running auto build
- **Registry** - Centralized service discovery and API endpoint explorer with feature rich metadata
- **Store** - Key-Value storage with TTL expiry and persistent crud to keep microservices stateless

**Framework**

Micro additionally now contains the incredibly popular Go Micro framework built in for service development. 
The Go framework makes it drop dead simple to write your services without having to piece together lines and lines of boilerplate. Auto 
configured and initialized by default, just import and get started quickly.

**Command Line**

Micro brings not only a rich architectural model but a command line experience tailored for that need. The command line interface includes 
dynamic command mapping for all services running on the platform. Turns any service instantly into a CLI command along with flag parsing 
for inputs. Includes support for multiple environments and namespaces, automatic refreshing of auth credentials, creating and running 
services, status info and log streaming, plus much, much more.

**Environments**

Finally Micro bakes in the concept of `Environments` and multi-tenancy through `Namespaces`. Run your server locally for 
development and in the cloud for staging and production, seamlessly switch between them using the CLI commands `micro env set [environment]` 
and `micro user set [namespace]`.

## Install

**From Source**

```sh
go get github.com/micro-community/micro/v3
```

**Using Docker**

*strong recommended on windows [details check](https://github.com/micro-community/micro/discussions/1650)*

```sh
# install
docker pull crazybber/micro

# run it
docker run -p 8080-8081:8080-8081/tcp crazybber/micro server
```

**Helm Chart**

```
helm repo add micro https://micro.github.io/helm
helm install micro micro/micro
```

**Release binaries**

```sh
# MacOS
curl -fsSL https://raw.githubusercontent.com/micro-community/micro/master/scripts/install.sh | /bin/bash

# Linux
wget -q  https://raw.githubusercontent.com/micro-community/micro/master/scripts/install.sh -O - | /bin/bash

# Windows
powershell -Command "iwr -useb https://raw.githubusercontent.com/micro-community/micro/master/scripts/install.ps1 | iex"
```

## Getting Started

Run the server locally(Recommended on Linux&Mac)

```
micro server
```

Set the environment to local (127.0.0.1:8081)

```
micro env set local
```

Login to the server

```
# user: admin pass: micro
micro login
```

Create a service

```sh
# generate a service (follow instructions in output)
micro new helloworld

# run the service
micro run helloworld

# check the status
micro status

# list running services
micro services

# call the service
micro helloworld --name=Alice

# curl via the api
curl -d '{"name": "Alice"}' http://localhost:8080/helloworld
```

## Example Service

Micro includes a Go framework for writing services wrapping gRPC for the core IDL and transport. 

Define services in proto:

```proto
syntax = "proto3";

package helloworld;

service Helloworld {
	rpc Call(Request) returns (Response) {}
}

message Request {
	string name = 1;
}

message Response {
	string msg = 1;
}
```

Write them using Go:

=======
Install micro

```sh
go install github.com/micro-community/micro/v3@latest
```

Run the server 

```sh
micro server
```

Login with the username 'admin' and password 'micro':

```sh
$ micro login
Enter username: admin
Enter password:
Successfully logged in.
```

See what's running:

```sh
$ micro services
api
auth
broker
config
events
network
proxy
registry
runtime
server
store
```

Run a service

```sh
micro run github.com/micro/services/helloworld
```

Now check the status of the running service

```sh
$ micro status
NAME		VERSION	SOURCE					STATUS	BUILD	UPDATED	METADATA
helloworld	latest	github.com/micro/services/helloworld	running	n/a	4s ago	owner=admin, group=micro
```

We can also have a look at logs of the service to verify it's running.

```sh
$ micro logs helloworld
2020-10-06 17:52:21  file=service/service.go:195 level=info Starting [service] helloworld
2020-10-06 17:52:21  file=grpc/grpc.go:902 level=info Server [grpc] Listening on [::]:33975
2020-10-06 17:52:21  file=grpc/grpc.go:732 level=info Registry [service] Registering node: helloworld-67627b23-3336-4b92-a032-09d8d13ecf95
```

Call the service

```sh
$ micro helloworld call --name=Jane
{
	"msg": "Hello Jane"
}
```

Curl it

```
curl "http://localhost:8080/helloworld?name=John"
```

Write a client


```go
package main

import (
	"context"
  
	"github.com/micro-community/micro/v3/service"
	"github.com/micro-community/micro/v3/service/logger"
	pb "github.com/micro/services/helloworld/proto"
)

type Helloworld struct{}

// Call is a single request handler called via client.Call or the generated client code
func (h *Helloworld) Call(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logger.Info("Received Helloworld.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

func main() {
	// Create service
	srv := service.New(
		service.Name("helloworld"),
	)

	// Register Handler
	srv.Handle(new(Helloworld))

	// Run the service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
```

Call with the client:

```go
import (
	"context"
  
	"github.com/micro-community/micro/v3/service/client"
	pb "github.com/micro-community/services/helloworld/proto"
)

// create a new helloworld service client
helloworld := pb.NewHelloworldService("helloworld", client.DefaultClient) 

// call the endpoint Helloworld.Call
rsp, err := helloworld.Call(context.Background(), &pb.Request{Name: "Alice"})
```

Curl it via the API

```
curl http://localhost:8080/helloworld?name=Alice
```

Hello world

```go
import (
	"fmt"
	"time"

	"github.com/micro-community/micro/v3/service"
	proto "github.com/micro-community/services/helloworld/proto"
)

func main() {
	// create and initialise a new service
	srv := service.New()

	// create the proto client for helloworld
	client := proto.NewHelloworldService("helloworld", srv.Client())

	// call an endpoint on the service
	rsp, err := client.Call(context.Background(), &proto.CallRequest{
		Name: "John",
	})
	if err != nil {
		fmt.Println("Error calling helloworld: ", err)
		return
	}

	// print the response
	fmt.Println("Response: ", rsp.Message)
	
	// let's delay the process for exiting for reasons you'll see below
	time.Sleep(time.Second * 5)
}
```

Run it

```
micro run .
```

For more see the [getting started](https://micro.dev/getting-started) guide.


## Usage

See the [docs](https://micro.dev/docs) for detailed information on the architecture, installation and use.

## License

See [LICENSE](LICENSE) which makes use of [Apache 2.0](https://opensource.org/licenses/Apache-2.0)

Join us on GitHub [Discussions](https://github.com/micro-community/micro/discussions).

## Repo Clone for CN

following cmd:

```bash
git clone https://hub.fastgit.org/micro-community/micro.git
cd micro
git remote remove orign
git remote add origin https://github.com/micro-community/micro.git
```
