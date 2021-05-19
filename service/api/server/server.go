// Package server is an API Gateway
package server

import (
	"fmt"

	"net/http"
	"os"

	"github.com/micro-community/micro/v3/plugin"
	"github.com/micro-community/micro/v3/service"
	"github.com/micro-community/micro/v3/service/api/auth"
	"github.com/micro-community/micro/v3/service/logger"
	"github.com/micro-community/micro/v3/service/registry"
	"github.com/micro-community/micro/v3/service/store"

	inApiHandler "github.com/micro-community/micro/v3/service/api/handler"
	inApiResolver "github.com/micro-community/micro/v3/service/api/resolver"
	inApiRouter "github.com/micro-community/micro/v3/service/api/router"
	inApiServer "github.com/micro-community/micro/v3/service/api/server/server"

	inApiHandlerApi "github.com/micro-community/micro/v3/service/api/handler/api"
	inApiHandlerEvent "github.com/micro-community/micro/v3/service/api/handler/event"
	inApiHandlerHttp "github.com/micro-community/micro/v3/service/api/handler/http"
	inApiHandlerRpc "github.com/micro-community/micro/v3/service/api/handler/rpc"
	inApiHandlerWeb "github.com/micro-community/micro/v3/service/api/handler/web"

	inApiResolverGrpc "github.com/micro-community/micro/v3/service/api/resolver/grpc"
	inApiResolverHost "github.com/micro-community/micro/v3/service/api/resolver/host"
	inApiResolverPath "github.com/micro-community/micro/v3/service/api/resolver/path"
	inApiResolverSubDomain "github.com/micro-community/micro/v3/service/api/resolver/subdomain"

	inApiRouterReg "github.com/micro-community/micro/v3/service/api/router/registry"
	inApiServerHttp "github.com/micro-community/micro/v3/service/api/server/http"

	"github.com/micro-community/micro/v3/util/acme"
	"github.com/micro-community/micro/v3/util/acme/autocert"
	"github.com/micro-community/micro/v3/util/acme/certmagic"

	inHelper "github.com/micro-community/micro/v3/util/helper"
	inSyncMemory "github.com/micro-community/micro/v3/util/sync/memory"

	"github.com/go-acme/lego/v3/providers/dns/cloudflare"
	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
)

//export default fields
var (
	Name                  = "api"
	Address               = ":8080"
	Handler               = "meta"
	Resolver              = "micro"
	APIPath               = "/"
	ProxyPath             = "/{service:[a-zA-Z0-9]+}"
	Namespace             = ""
	ACMEProvider          = "autocert"
	ACMEChallengeProvider = "cloudflare"
	ACMECA                = acme.LetsEncryptProductionCA
)

var (
	Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "address",
			Usage:   "Set the api address e.g 0.0.0.0:8080",
			EnvVars: []string{"MICRO_API_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "handler",
			Usage:   "Specify the request handler to be used for mapping HTTP requests to services; {api, event, http, rpc}",
			EnvVars: []string{"MICRO_API_HANDLER"},
		},
		&cli.StringFlag{
			Name:    "namespace",
			Usage:   "Set the namespace used by the API e.g. com.example",
			EnvVars: []string{"MICRO_API_NAMESPACE"},
		},
		&cli.StringFlag{
			Name:    "resolver",
			Usage:   "Set the hostname resolver used by the API {host, path, grpc}",
			EnvVars: []string{"MICRO_API_RESOLVER"},
		},
		&cli.BoolFlag{
			Name:    "enable_cors",
			Usage:   "Enable CORS, allowing the API to be called by frontend applications",
			EnvVars: []string{"MICRO_API_ENABLE_CORS"},
			Value:   true,
		},
		&cli.BoolFlag{
			Name:    "enable_acme",
			Usage:   "Enables ACME support via Let's Encrypt. ACME hosts should also be specified.",
			EnvVars: []string{"MICRO_API_ENABLE_ACME"},
		},
		&cli.StringFlag{
			Name:    "acme_hosts",
			Usage:   "Comma separated list of hostnames to manage ACME certs for",
			EnvVars: []string{"MICRO_API_ACME_HOSTS"},
		},
		&cli.StringFlag{
			Name:    "acme_provider",
			Usage:   "The provider that will be used to communicate with Let's Encrypt. Valid options: autocert, certmagic",
			EnvVars: []string{"MICRO_API_ACME_PROVIDER"},
		},
		&cli.BoolFlag{
			Name:    "enable_tls",
			Usage:   "Enable TLS support. Expects cert and key file to be specified",
			EnvVars: []string{"MICRO_API_ENABLE_TLS"},
		},
		&cli.StringFlag{
			Name:    "tls_cert_file",
			Usage:   "Path to the TLS Certificate file",
			EnvVars: []string{"MICRO_API_TLS_CERT_FILE"},
		},
		&cli.StringFlag{
			Name:    "tls_key_file",
			Usage:   "Path to the TLS Key file",
			EnvVars: []string{"MICRO_API_TLS_KEY_FILE"},
		},
		&cli.StringFlag{
			Name:    "tls_client_ca_file",
			Usage:   "Path to the TLS CA file to verify clients against",
			EnvVars: []string{"MICRO_API_TLS_CLIENT_CA_FILE"},
		},
	}
)

func Run(ctx *cli.Context) error {
	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(ctx.String("handler")) > 0 {
		Handler = ctx.String("handler")
	}
	if len(ctx.String("resolver")) > 0 {
		Resolver = ctx.String("resolver")
	}
	if len(ctx.String("acme_provider")) > 0 {
		ACMEProvider = ctx.String("acme_provider")
	}
	if len(ctx.String("namespace")) > 0 {
		Namespace = ctx.String("namespace")
	}
	if len(ctx.String("api_handler")) > 0 {
		Handler = ctx.String("api_handler")
	}
	if len(ctx.String("api_address")) > 0 {
		Address = ctx.String("api_address")
	}
	// initialize internal service
	inSrv := service.New(service.Name(Name))

	// Init API
	var opts []inApiServer.Option

	if ctx.Bool("enable_acme") {
		hosts := inHelper.ACMEHosts(ctx)
		opts = append(opts, inApiServer.EnableACME(true))
		opts = append(opts, inApiServer.ACMEHosts(hosts...))
		switch ACMEProvider {
		case "autocert":
			opts = append(opts, inApiServer.ACMEProvider(autocert.NewProvider()))
		case "certmagic":
			if ACMEChallengeProvider != "cloudflare" {
				logger.Fatal("The only implemented DNS challenge provider is cloudflare")
			}

			apiToken := os.Getenv("CF_API_TOKEN")
			if len(apiToken) == 0 {
				logger.Fatal("env variables CF_API_TOKEN and CF_ACCOUNT_ID must be set")
			}

			storage := certmagic.NewStorage(
				inSyncMemory.NewSync(),
				store.DefaultStore,
			)

			config := cloudflare.NewDefaultConfig()
			config.AuthToken = apiToken
			config.ZoneToken = apiToken
			challengeProvider, err := cloudflare.NewDNSProviderConfig(config)
			if err != nil {
				logger.Fatal(err.Error())
			}

			opts = append(opts,
				inApiServer.ACMEProvider(
					certmagic.NewProvider(
						acme.AcceptToS(true),
						acme.CA(ACMECA),
						acme.Cache(storage),
						acme.ChallengeProvider(challengeProvider),
						acme.OnDemand(false),
					),
				),
			)
		default:
			logger.Fatalf("%s is not a valid ACME provider\n", ACMEProvider)
		}
	} else if ctx.Bool("enable_tls") {
		config, err := inHelper.TLSConfig(ctx)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		opts = append(opts, inApiServer.EnableTLS(true))
		opts = append(opts, inApiServer.TLSConfig(config))
	}

	if ctx.Bool("enable_cors") {
		opts = append(opts, inApiServer.EnableCORS(true))
	}

	// create the router
	var h http.Handler
	r := mux.NewRouter()
	h = r

	// return version and list of services
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			return
		}

		response := fmt.Sprintf(`{"version": "%s"}`, ctx.App.Version)
		w.Write([]byte(response))
	})

	// strip favicon.ico
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	// resolver options
	resolverOpts := []inApiResolver.Option{
		inApiResolver.WithServicePrefix(Namespace),
		inApiResolver.WithHandler(Handler),
	}

	// default resolver
	rr := NewResolver(resolverOpts...)

	switch Resolver {
	case "subdomain":
		rr = inApiResolverSubDomain.NewResolver(rr)
	case "host":
		rr = inApiResolverHost.NewResolver(resolverOpts...)
	case "path":
		rr = inApiResolverPath.NewResolver(resolverOpts...)
	case "grpc":
		rr = inApiResolverGrpc.NewResolver(resolverOpts...)
	}

	switch Handler {
	case "rpc":
		logger.Infof("Registering API RPC Handler at %s", APIPath)
		rt := inApiRouterReg.NewRouter(
			inApiRouter.WithHandler(inApiHandlerRpc.Handler),
			inApiRouter.WithResolver(rr),
			inApiRouter.WithRegistry(registry.DefaultRegistry),
		)
		rp := inApiHandlerRpc.NewHandler(
			inApiHandler.WithNamespace(Namespace),
			inApiHandler.WithRouter(rt),
			inApiHandler.WithClient(inSrv.Client()),
		)
		r.PathPrefix(APIPath).Handler(rp)
	case "api":
		logger.Infof("Registering API Request Handler at %s", APIPath)
		rt := inApiRouterReg.NewRouter(
			inApiRouter.WithHandler(inApiHandlerApi.Handler),
			inApiRouter.WithResolver(rr),
			inApiRouter.WithRegistry(registry.DefaultRegistry),
		)
		ap := inApiHandlerApi.NewHandler(
			inApiHandler.WithNamespace(Namespace),
			inApiHandler.WithRouter(rt),
			inApiHandler.WithClient(inSrv.Client()),
		)
		r.PathPrefix(APIPath).Handler(ap)
	case "event":
		logger.Infof("Registering API Event Handler at %s", APIPath)
		rt := inApiRouterReg.NewRouter(
			inApiRouter.WithHandler(inApiHandlerEvent.Handler),
			inApiRouter.WithResolver(rr),
			inApiRouter.WithRegistry(registry.DefaultRegistry),
		)
		ev := inApiHandlerEvent.NewHandler(
			inApiHandler.WithNamespace(Namespace),
			inApiHandler.WithRouter(rt),
			inApiHandler.WithClient(inSrv.Client()),
		)
		r.PathPrefix(APIPath).Handler(ev)
	case "http":
		logger.Infof("Registering API HTTP Handler at %s", ProxyPath)
		rt := inApiRouterReg.NewRouter(
			inApiRouter.WithHandler(inApiHandlerHttp.Handler),
			inApiRouter.WithResolver(rr),
			inApiRouter.WithRegistry(registry.DefaultRegistry),
		)
		ht := inApiHandlerHttp.NewHandler(
			inApiHandler.WithNamespace(Namespace),
			inApiHandler.WithRouter(rt),
			inApiHandler.WithClient(inSrv.Client()),
		)
		r.PathPrefix(ProxyPath).Handler(ht)
	case "web":
		logger.Infof("Registering API Web Handler at %s", APIPath)
		rt := inApiRouterReg.NewRouter(
			inApiRouter.WithHandler(inApiHandlerWeb.Handler),
			inApiRouter.WithResolver(rr),
			inApiRouter.WithRegistry(registry.DefaultRegistry),
		)
		w := inApiHandlerWeb.NewHandler(
			inApiHandler.WithNamespace(Namespace),
			inApiHandler.WithRouter(rt),
			inApiHandler.WithClient(inSrv.Client()),
		)
		r.PathPrefix(APIPath).Handler(w)
	default:
		logger.Infof("Registering API Default Handler at %s", APIPath)
		rt := inApiRouterReg.NewRouter(
			inApiRouter.WithResolver(rr),
			inApiRouter.WithRegistry(registry.DefaultRegistry),
		)
		r.PathPrefix(APIPath).Handler(Meta(inSrv, rt, Namespace))
	}

	// register all the http handler plugins
	for _, p := range plugin.Plugins() {
		if v := p.Handler(); v != nil {
			h = v(h)
		}
	}

	// append the auth wrapper
	h = auth.Wrapper(rr, Namespace)(h)

	// create a new api server with wrappers
	api := inApiServerHttp.NewServer(Address)
	// initialize
	api.Init(opts...)
	// register the handler
	api.Handle("/", h)

	// Start API
	if err := api.Start(); err != nil {
		logger.Fatal(err)
	}

	// Run internal server
	if err := inSrv.Run(); err != nil {
		logger.Fatal(err)
	}

	// Stop API
	if err := api.Stop(); err != nil {
		logger.Fatal(err)
	}

	return nil
}
