// Package profile is for specific profiles
// @todo this package is the definition of cruft and
// should be rewritten in a more elegant way
package profile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/micro-community/micro/v3/service/auth"
	"github.com/micro-community/micro/v3/service/broker"
	"github.com/micro-community/micro/v3/service/build"
	"github.com/micro-community/micro/v3/service/client"
	"github.com/micro-community/micro/v3/service/config"
	"github.com/micro-community/micro/v3/service/events"
	"github.com/micro-community/micro/v3/service/logger"
	"github.com/micro-community/micro/v3/service/model"
	"github.com/micro-community/micro/v3/service/registry"
	"github.com/micro-community/micro/v3/service/router"
	"github.com/micro-community/micro/v3/service/runtime"
	"github.com/micro-community/micro/v3/service/server"
	"github.com/micro-community/micro/v3/service/store"

	mAuthJwt "github.com/micro-community/micro/v3/service/auth/jwt"
	mAuthNoop "github.com/micro-community/micro/v3/service/auth/noop"
	mBrokerMemory "github.com/micro-community/micro/v3/service/broker/memory"
	mBuildGolang "github.com/micro-community/micro/v3/service/build/golang"
	mConfigStore "github.com/micro-community/micro/v3/service/config/store"
	mEventStore "github.com/micro-community/micro/v3/service/events/store"
	mEventStream "github.com/micro-community/micro/v3/service/events/stream/memory"
	mRegMemory "github.com/micro-community/micro/v3/service/registry/memory"
	mRouterK8s "github.com/micro-community/micro/v3/service/router/kubernetes"
	mRouterReg "github.com/micro-community/micro/v3/service/router/registry"
	mRuntimeK8s "github.com/micro-community/micro/v3/service/runtime/kubernetes"
	mRuntimeLocal "github.com/micro-community/micro/v3/service/runtime/local"
	mStoreFile "github.com/micro-community/micro/v3/service/store/file"
	mStoreMemory "github.com/micro-community/micro/v3/service/store/memory"

	inAuth "github.com/micro-community/micro/v3/util/auth"
	inUser "github.com/micro-community/micro/v3/util/user"

	"github.com/urfave/cli/v2"
)

// profiles which when called will configure micro to run in that environment
var profiles = map[string]*Profile{
	// built in profiles
	"client":     Client,
	"service":    Service,
	"test":       Test,
	"local":      Local,
	"kubernetes": Kubernetes,
	"simple":     Simple,
	"dev":        Dev,
}

// Profile configures an environment
type Profile struct {
	// name of the profile
	Name string
	// function used for setup
	Setup func(*cli.Context) error
	// TODO: presetup dependencies
	// e.g start resources
}

// Register a profile
func Register(name string, p *Profile) error {
	if _, ok := profiles[name]; ok {
		return fmt.Errorf("profile %s already exists", name)
	}
	profiles[name] = p
	return nil
}

// Load a profile
func Load(name string) (*Profile, error) {
	v, ok := profiles[name]
	if !ok {
		return nil, fmt.Errorf("profile %s does not exist", name)
	}
	return v, nil
}

// Client profile is for any entrypoint that behaves as a client
var Client = &Profile{
	Name:  "client",
	Setup: func(ctx *cli.Context) error { return nil },
}

// Local profile to run locally
var Local = &Profile{
	Name: "local",
	Setup: func(ctx *cli.Context) error {
		auth.DefaultAuth = mAuthJwt.NewAuth()
		store.DefaultStore = mStoreFile.NewStore(mStoreFile.WithDir(filepath.Join(inUser.Dir, "server", "store")))
		SetupConfigSecretKey(ctx)
		config.DefaultConfig, _ = mConfigStore.NewConfig(store.DefaultStore, "")
		// SetupBroker(mBrokerMemory.NewBroker())
		// SetupRegistry(mRegMdns.NewRegistry())
		SetupJWT(ctx)

		// the registry service uses the memory registry, the other core services will use the default
		// rpc client and call the registry service
		if ctx.Args().Get(1) == "registry" {
			SetupRegistry(mRegMemory.NewRegistry())
		} else {
			// set the registry address
			registry.DefaultRegistry.Init(
				registry.Addrs("localhost:8000"),
			)

			SetupRegistry(registry.DefaultRegistry)
		}

		// the broker service uses the memory broker, the other core services will use the default
		// rpc client and call the broker service
		if ctx.Args().Get(1) == "broker" {
			SetupBroker(mBrokerMemory.NewBroker())
		} else {
			broker.DefaultBroker.Init(
				broker.Addrs("localhost:8003"),
			)
			SetupBroker(broker.DefaultBroker)
		}

		// set the store in the model
		model.DefaultModel = model.NewModel(
			model.WithStore(store.DefaultStore),
		)

		// use the local runtime, note: the local runtime is designed to run source code directly so
		// the runtime builder should NOT be set when using this implementation
		runtime.DefaultRuntime = mRuntimeLocal.NewRuntime()

		var err error
		events.DefaultStream, err = mEventStream.NewStream()
		if err != nil {
			logger.Fatalf("Error configuring stream: %v", err)
		}
		events.DefaultStore = mEventStore.NewStore(
			mEventStore.WithStore(store.DefaultStore),
		)

		store.DefaultBlobStore, err = mStoreFile.NewBlobStore()
		if err != nil {
			logger.Fatalf("Error configuring file blob store: %v", err)
		}

		return nil
	},
}

// Kubernetes profile to run on kubernetes with zero deps. Designed for use with the micro helm chart
var Kubernetes = &Profile{
	Name: "kubernetes",
	Setup: func(ctx *cli.Context) (err error) {
		auth.DefaultAuth = mAuthJwt.NewAuth()
		SetupJWT(ctx)

		runtime.DefaultRuntime = mRuntimeK8s.NewRuntime()
		build.DefaultBuilder, err = mBuildGolang.NewBuilder()
		if err != nil {
			logger.Fatalf("Error configuring golang builder: %v", err)
		}

		events.DefaultStream, err = mEventStream.NewStream()
		if err != nil {
			logger.Fatalf("Error configuring stream: %v", err)
		}

		store.DefaultStore = mStoreFile.NewStore(mStoreFile.WithDir("/store"))
		store.DefaultBlobStore, err = mStoreFile.NewBlobStore(mStoreFile.WithDir("/store/blob"))
		if err != nil {
			logger.Fatalf("Error configuring file blob store: %v", err)
		}

		// set the store in the model
		model.DefaultModel = model.NewModel(
			model.WithStore(store.DefaultStore),
		)

		// the registry service uses the memory registry, the other core services will use the default
		// rpc client and call the registry service
		if ctx.Args().Get(1) == "registry" {
			SetupRegistry(mRegMemory.NewRegistry())
		}

		// the broker service uses the memory broker, the other core services will use the default
		// rpc client and call the broker service
		if ctx.Args().Get(1) == "broker" {
			SetupBroker(mBrokerMemory.NewBroker())
		}

		config.DefaultConfig, err = mConfigStore.NewConfig(store.DefaultStore, "")
		if err != nil {
			logger.Fatalf("Error configuring config: %v", err)
		}
		SetupConfigSecretKey(ctx)

		// Use k8s routing which is DNS based
		router.DefaultRouter = mRouterK8s.NewRouter()
		client.DefaultClient.Init(client.Router(router.DefaultRouter))
		return nil
	},
}

// Service is the default for any services run
var Service = &Profile{
	Name:  "service",
	Setup: func(ctx *cli.Context) error { return nil },
}

// Test profile is used for the go test suite
var Test = &Profile{
	Name: "test",
	Setup: func(ctx *cli.Context) error {
		auth.DefaultAuth = mAuthNoop.NewAuth()
		store.DefaultStore = mStoreMemory.NewStore()
		store.DefaultBlobStore, _ = mStoreFile.NewBlobStore()
		config.DefaultConfig, _ = mConfigStore.NewConfig(store.DefaultStore, "")
		SetupRegistry(mRegMemory.NewRegistry())
		// set the store in the model
		model.DefaultModel = model.NewModel(
			model.WithStore(store.DefaultStore),
		)
		return nil
	},
}

// SetupRegistry configures the registry
func SetupRegistry(reg registry.Registry) {
	registry.DefaultRegistry = reg
	router.DefaultRouter = mRouterReg.NewRouter(router.Registry(reg))
	client.DefaultClient.Init(client.Registry(reg), client.Router(router.DefaultRouter))
	server.DefaultServer.Init(server.Registry(reg))
}

// SetupBroker configures the broker
func SetupBroker(b broker.Broker) {
	broker.DefaultBroker = b
	client.DefaultClient.Init(client.Broker(b))
	server.DefaultServer.Init(server.Broker(b))
}

// SetupJWT configures the default internal system rules
func SetupJWT(ctx *cli.Context) {
	for _, rule := range inAuth.SystemRules {
		if err := auth.DefaultAuth.Grant(rule); err != nil {
			logger.Fatal("Error creating default rule: %v", err)
		}
	}
}

//SetupConfigSecretKey config default SecretKey
func SetupConfigSecretKey(ctx *cli.Context) {
	key := ctx.String("config_secret_key")
	if len(key) == 0 {
		k, err := inUser.GetConfigSecretKey()
		if err != nil {
			logger.Fatal("Error getting config secret: %v", err)
		}
		os.Setenv("MICRO_CONFIG_SECRET_KEY", k)
	}
}
