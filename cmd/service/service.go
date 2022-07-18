// Package service provides micro service commands
package service

import (
	"fmt"
	"os"
	"strings"

	"github.com/micro-community/micro/v3/cmd"
	"github.com/micro-community/micro/v3/plugin"
	"github.com/micro-community/micro/v3/service"
	"github.com/micro-community/micro/v3/service/logger"
	"github.com/micro-community/micro/v3/service/proxy"
	"github.com/micro-community/micro/v3/service/runtime"
	"github.com/micro-community/micro/v3/service/server"
	"github.com/micro-community/micro/v3/service/web"
	ccli "github.com/urfave/cli/v2"

	"github.com/micro-community/micro/v3/service/proxy/grpc"
	"github.com/micro-community/micro/v3/service/proxy/http"
	"github.com/micro-community/micro/v3/service/proxy/mucp"

	// services: server implementation
	mApiServer "github.com/micro-community/micro/v3/service/api/server"
	mAuthServer "github.com/micro-community/micro/v3/service/auth/server"
	mBrokerServer "github.com/micro-community/micro/v3/service/broker/server"
	mConfigServer "github.com/micro-community/micro/v3/service/config/server"
	mEventsServer "github.com/micro-community/micro/v3/service/events/server"
	mNetworkServer "github.com/micro-community/micro/v3/service/network/server"
	mProxyServer "github.com/micro-community/micro/v3/service/proxy/server"
	mRegistryServer "github.com/micro-community/micro/v3/service/registry/server"
	mRuntimeServer "github.com/micro-community/micro/v3/service/runtime/server"
	mStoreServer "github.com/micro-community/micro/v3/service/store/server"

	// misc commands
	"github.com/micro-community/micro/v3/cmd/service/handler/exec"
	"github.com/micro-community/micro/v3/cmd/service/handler/file"
)

// Run starts a micro service sidecar to encapsulate any app
func Run(ctx *ccli.Context) {
	name := ctx.String("name")
	address := ctx.String("address")
	endpoint := ctx.String("endpoint")

	metadata := make(map[string]string)
	for _, md := range ctx.StringSlice("metadata") {
		parts := strings.Split(md, "=")
		if len(parts) < 2 {
			continue
		}
		key := parts[0]
		val := strings.Join(parts[1:], "=")

		// set the key/val
		metadata[key] = val
	}

	var opts []service.Option
	if len(metadata) > 0 {
		opts = append(opts, service.Metadata(metadata))
	}
	if len(name) > 0 {
		opts = append(opts, service.Name(name))
	}
	if len(address) > 0 {
		opts = append(opts, service.Address(address))
	}
	if len(endpoint) == 0 {
		endpoint = proxy.DefaultEndpoint
	}

	var p proxy.Proxy

	switch {
	case strings.HasPrefix(endpoint, "grpc"):
		endpoint = strings.TrimPrefix(endpoint, "grpc://")
		p = grpc.NewProxy(proxy.WithEndpoint(endpoint))
	case strings.HasPrefix(endpoint, "http"):
		p = http.NewProxy(proxy.WithEndpoint(endpoint))
	case strings.HasPrefix(endpoint, "file"):
		endpoint = strings.TrimPrefix(endpoint, "file://")
		p = file.NewProxy(proxy.WithEndpoint(endpoint))
	case strings.HasPrefix(endpoint, "exec"):
		endpoint = strings.TrimPrefix(endpoint, "exec://")
		p = exec.NewProxy(proxy.WithEndpoint(endpoint))
	default:
		p = mucp.NewProxy(proxy.WithEndpoint(endpoint))
	}

	// run the service if asked to
	if ctx.Args().Len() > 0 {
		args := []runtime.CreateOption{
			runtime.WithCommand(ctx.Args().Slice()...),
			runtime.WithOutput(os.Stdout),
		}

		// create new local runtime
		r := runtime.DefaultRuntime

		// start the runtime
		r.Start()

		// register the service
		r.Create(&runtime.Service{Name: name}, args...)

		// stop the runtime
		defer func() {
			r.Delete(&runtime.Service{Name: name})
			r.Stop()
		}()
	}

	logger.Infof("Service [%s] Serving %s at endpoint %s\n", p.String(), name, endpoint)

	// new service
	srv := service.New(opts...)

	// create new muxer
	//	muxer := mux.New(name, p)

	// set the router
	srv.Server().Init(server.WithRouter(p))

	// run service
	srv.Run()
}

type srvCommand struct {
	Name    string
	Command ccli.ActionFunc
	Flags   []ccli.Flag
}

var srvCommands = []srvCommand{
	{
		Name:    "api",
		Command: mApiServer.Run,
		Flags:   mApiServer.Flags,
	},
	{
		Name:    "auth",
		Command: mAuthServer.Run,
		Flags:   mAuthServer.Flags,
	},
	{
		Name:    "broker",
		Command: mBrokerServer.Run,
	},
	{
		Name:    "config",
		Command: mConfigServer.Run,
		Flags:   mConfigServer.Flags,
	},
	{
		Name:    "events",
		Command: mEventsServer.Run,
	},
	{
		Name:    "network",
		Command: mNetworkServer.Run,
		Flags:   mNetworkServer.Flags,
	},
	{
		Name:    "proxy",
		Command: mProxyServer.Run,
		Flags:   mProxyServer.Flags,
	},
	{
		Name:    "registry",
		Command: mRegistryServer.Run,
	},
	{
		Name:    "runtime",
		Command: mRuntimeServer.Run,
		Flags:   mRuntimeServer.Flags,
	},
	{
		Name:    "store",
		Command: mStoreServer.Run,
	},
	{
		Name:    "web",
		Command: web.Run,
		Flags:   web.Flags,
	},
}

func init() {
	// move newAction outside the loop and pass c as an arg to
	// set the scope of the variable
	newAction := func(c srvCommand) func(ctx *ccli.Context) error {
		return func(ctx *ccli.Context) error {
			// configure the loggerger
			logger.DefaultLogger.Init(logger.WithFields(map[string]interface{}{"service": c.Name}))

			// run the service
			c.Command(ctx)
			return nil
		}
	}

	subcommands := make([]*ccli.Command, len(srvCommands))
	for i, c := range srvCommands {
		// construct the command
		command := &ccli.Command{
			Name:   c.Name,
			Flags:  c.Flags,
			Usage:  fmt.Sprintf("Run micro %v", c.Name),
			Action: newAction(c),
		}

		// setup the plugins
		for _, p := range plugin.Plugins(plugin.Module(c.Name)) {
			if cmds := p.Commands(); len(cmds) > 0 {
				command.Subcommands = append(command.Subcommands, cmds...)
			}

			if flags := p.Flags(); len(flags) > 0 {
				command.Flags = append(command.Flags, flags...)
			}
		}

		// set the command
		subcommands[i] = command
	}

	command := &ccli.Command{
		Name:  "service",
		Usage: "Run a micro service",
		Action: func(ctx *ccli.Context) error {
			Run(ctx)
			return nil
		},
		Flags: []ccli.Flag{
			&ccli.StringFlag{
				Name:    "name",
				Usage:   "Name of the service",
				EnvVars: []string{"MICRO_SERVICE_NAME"},
				Value:   "service",
			},
			&ccli.StringFlag{
				Name:    "address",
				Usage:   "Address of the service",
				EnvVars: []string{"MICRO_SERVICE_ADDRESS"},
			},
			&ccli.StringFlag{
				Name:    "endpoint",
				Usage:   "The local service endpoint (Defaults to localhost:9090); {http, grpc, file, exec}://path-or-address e.g http://localhost:9090",
				EnvVars: []string{"MICRO_SERVICE_ENDPOINT"},
			},
			&ccli.StringSliceFlag{
				Name:    "metadata",
				Usage:   "Add metadata as key-value pairs describing the service e.g owner=john@example.com",
				EnvVars: []string{"MICRO_SERVICE_METADATA"},
			},
		},
		Subcommands: subcommands,
	}

	// register global plugins and flags
	for _, p := range plugin.Plugins() {
		if cmds := p.Commands(); len(cmds) > 0 {
			command.Subcommands = append(command.Subcommands, cmds...)
		}

		if flags := p.Flags(); len(flags) > 0 {
			command.Flags = append(command.Flags, flags...)
		}
	}

	cmd.Register(command)
}
