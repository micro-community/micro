// Package proxy is a proxy for grpc/http/mucp
package proxy

import (
	"os"
	"strings"

	mCli "github.com/micro-community/micro/v3/client"
	"github.com/micro-community/micro/v3/service"
	"github.com/micro-community/micro/v3/service/client"
	"github.com/micro-community/micro/v3/service/logger"
	"github.com/micro-community/micro/v3/service/proxy"
	"github.com/micro-community/micro/v3/service/router"
	"github.com/micro-community/micro/v3/service/server"
	"github.com/micro-community/micro/v3/service/store"

	//Platform support
	"github.com/micro-community/micro/v3/platform/api/server/acme"
	"github.com/micro-community/micro/v3/platform/api/server/acme/autocert"
	"github.com/micro-community/micro/v3/platform/api/server/acme/certmagic"
	"github.com/micro-community/micro/v3/platform/helper"
	"github.com/micro-community/micro/v3/platform/muxer"
	"github.com/micro-community/micro/v3/platform/sync/memory"

	mProxyGrpc "github.com/micro-community/micro/v3/service/proxy/grpc"
	mProxyHttp "github.com/micro-community/micro/v3/service/proxy/http"
	mProxyMucp "github.com/micro-community/micro/v3/service/proxy/mucp"

	mBrokerMemory "github.com/micro-community/micro/v3/service/broker/memory"
	mRegistryNoop "github.com/micro-community/micro/v3/service/registry/noop"
	mServerGrpc "github.com/micro-community/micro/v3/service/server/grpc"

	"github.com/urfave/cli/v2"

	"github.com/go-acme/lego/v3/providers/dns/cloudflare"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
)

//service for proxy
var (
	// Name of the proxy
	Name = "proxy"
	// The address of the proxy
	Address = ":8081"
	// Is gRPCWeb enabled
	GRPCWebEnabled = false
	// The address of the proxy
	GRPCWebAddress = ":8082"
	// the proxy protocol
	Protocol = "grpc"
	// The endpoint host to route to
	Endpoint string
	// ACME (Cert management)
	ACMEProvider          = "autocert"
	ACMEChallengeProvider = "cloudflare"
	ACMECA                = acme.LetsEncryptProductionCA

	Flags = append(mCli.Flags,
		&cli.StringFlag{
			Name:    "address",
			Usage:   "Set the proxy http address e.g 0.0.0.0:8081",
			EnvVars: []string{"MICRO_PROXY_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "protocol",
			Usage:   "Set the protocol used for proxying e.g mucp, grpc, http",
			EnvVars: []string{"MICRO_PROXY_PROTOCOL"},
		},
		&cli.StringFlag{
			Name:    "endpoint",
			Usage:   "Set the endpoint to route to e.g greeter or localhost:9090",
			EnvVars: []string{"MICRO_PROXY_ENDPOINT"},
		},
		&cli.BoolFlag{
			Name:    "grpc_web",
			Usage:   "Enable the gRPCWeb server",
			EnvVars: []string{"MICRO_PROXY_GRPC_WEB"},
		},
		&cli.StringFlag{
			Name:    "grpc_web_addr",
			Usage:   "Set the gRPC web addr on the proxy",
			EnvVars: []string{"MICRO_PROXY_GRPC_WEB_ADDRESS"},
		},
	)
)

//Run proxy
func Run(ctx *cli.Context) error {
	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if ctx.Bool("grpc_web") {
		GRPCWebEnabled = ctx.Bool("grpc_web")
	}
	if len(ctx.String("grpc_web_address")) > 0 {
		GRPCWebAddress = ctx.String("grpc_web_address")
	}
	if len(ctx.String("endpoint")) > 0 {
		Endpoint = ctx.String("endpoint")
	}
	if len(ctx.String("protocol")) > 0 {
		Protocol = ctx.String("protocol")
	}
	if len(ctx.String("acme_provider")) > 0 {
		ACMEProvider = ctx.String("acme_provider")
	}

	// set the context
	pOpts := []proxy.Option{
		proxy.WithRouter(router.DefaultRouter),
		proxy.WithClient(client.DefaultClient),
	}

	// set endpoint
	if len(Endpoint) > 0 {
		ep := Endpoint

		switch {
		case strings.HasPrefix(Endpoint, "grpc://"):
			ep = strings.TrimPrefix(Endpoint, "grpc://")
			Protocol = "grpc"
		case strings.HasPrefix(Endpoint, "http://"):
			Protocol = "http"
		case strings.HasPrefix(Endpoint, "mucp://"):
			ep = strings.TrimPrefix(Endpoint, "mucp://")
			Protocol = "mucp"
		}

		pOpts = append(pOpts, proxy.WithEndpoint(ep))
	}

	serverOpts := []server.Option{
		server.Name(Name),
		server.Address(Address),
		server.Registry(mRegistryNoop.NewRegistry()),
		server.Broker(mBrokerMemory.NewBroker()),
	}

	// enable acme will create a net.Listener which
	if ctx.Bool("enable_acme") {
		var ap acme.Provider

		switch ACMEProvider {
		case "autocert":
			ap = autocert.NewProvider()
		case "certmagic":
			if ACMEChallengeProvider != "cloudflare" {
				logger.Fatal("The only implemented DNS challenge provider is cloudflare")
			}

			apiToken := os.Getenv("CF_API_TOKEN")
			if len(apiToken) == 0 {
				logger.Fatal("env variables CF_API_TOKEN and CF_ACCOUNT_ID must be set")
			}

			storage := certmagic.NewStorage(memory.NewSync(), store.DefaultStore)

			config := cloudflare.NewDefaultConfig()
			config.AuthToken = apiToken
			config.ZoneToken = apiToken
			challengeProvider, err := cloudflare.NewDNSProviderConfig(config)
			if err != nil {
				logger.Fatal(err.Error())
			}

			// define the provider
			ap = certmagic.NewProvider(
				acme.AcceptToS(true),
				acme.CA(ACMECA),
				acme.Cache(storage),
				acme.ChallengeProvider(challengeProvider),
				acme.OnDemand(false),
			)
		default:
			logger.Fatalf("Unsupported acme provider: %s\n", ACMEProvider)
		}

		// generate the tls config
		config, err := ap.TLSConfig(helper.ACMEHosts(ctx)...)
		if err != nil {
			logger.Fatalf("Failed to generate acme tls config: %v", err)
		}

		// set the tls config
		serverOpts = append(serverOpts, server.TLSConfig(config))
		// enable tls will leverage tls certs and generate a tls.Config
	} else if ctx.Bool("enable_tls") {
		// get certificates from the context
		config, err := helper.TLSConfig(ctx)
		if err != nil {
			logger.Fatal(err)
			return err
		}
		serverOpts = append(serverOpts, server.TLSConfig(config))
	}

	// new proxy
	var p proxy.Proxy

	// set proxy
	switch Protocol {
	case "http":
		p = mProxyHttp.NewProxy(pOpts...)
		// TODO: http server
	case "mucp":
		p = mProxyMucp.NewProxy(pOpts...)
	default:
		// default to the grpc proxy
		p = mProxyGrpc.NewProxy(pOpts...)
	}

	// wrap the proxy using the proxy's authHandler
	authOpt := server.WrapHandler(authHandler())
	serverOpts = append(serverOpts, authOpt)
	serverOpts = append(serverOpts, server.WithRouter(p))

	if len(Endpoint) > 0 {
		logger.Infof("Proxy [%s] serving endpoint: %s", p.String(), Endpoint)
	} else {
		logger.Infof("Proxy [%s] serving protocol: %s", p.String(), Protocol)
	}

	if GRPCWebEnabled {
		serverOpts = append(serverOpts, mServerGrpc.GRPCWebPort(GRPCWebAddress))
		serverOpts = append(serverOpts, mServerGrpc.GRPCWebOptions(
			grpcweb.WithCorsForRegisteredEndpointsOnly(false),
			grpcweb.WithOriginFunc(func(origin string) bool { return true })))

		logger.Infof("Proxy [%s] serving gRPC-Web on %s", p.String(), GRPCWebAddress)
	}

	// create a new grpc server
	srv := mServerGrpc.NewServer(serverOpts...)

	// Start the proxy server
	if err := srv.Start(); err != nil {
		logger.Fatal(err)
	}
	// create a new proxy muxer which includes the debug handler
	muxer := muxer.New(Name, p)

	inSrvOpts := []server.Option{
		server.Registry(mRegistryNoop.NewRegistry()),
		server.WithRouter(muxer),
	}

	// new internal service
	inSrv := service.New(service.Name(Name))

	// set the router
	inSrv.Server().Init(inSrvOpts...)

	// Run internal service
	if err := inSrv.Run(); err != nil {
		logger.Fatal(err)
	}

	// Stop the server
	if err := srv.Stop(); err != nil {
		logger.Fatal(err)
	}

	return nil
}
