package profile

import (
	"github.com/micro-community/micro/v3/service/auth"
	"github.com/micro-community/micro/v3/service/client"
	"github.com/micro-community/micro/v3/service/config"
	"github.com/micro-community/micro/v3/service/events"
	"github.com/micro-community/micro/v3/service/logger"
	"github.com/micro-community/micro/v3/service/registry"
	"github.com/micro-community/micro/v3/service/router"
	"github.com/micro-community/micro/v3/service/runtime"
	"github.com/micro-community/micro/v3/service/server"
	"github.com/micro-community/micro/v3/service/store"

	mAuthNoop "github.com/micro-community/micro/v3/service/auth/noop"
	mBrokerMemory "github.com/micro-community/micro/v3/service/broker/memory"
	mConfigEnv "github.com/micro-community/micro/v3/service/config/env"
	mEventStream "github.com/micro-community/micro/v3/service/events/stream/memory"
	mRegistryNoop "github.com/micro-community/micro/v3/service/registry/noop"
	mRouterStatic "github.com/micro-community/micro/v3/service/router/static"
	mRuntimeLocal "github.com/micro-community/micro/v3/service/runtime/local"
	mStoreMemory "github.com/micro-community/micro/v3/service/store/memory"

	"github.com/urfave/cli/v2"
)

// Simple profile to run service in simple config
var Simple = &Profile{
	Name: "simple",
	Setup: func(ctx *cli.Context) error {
		auth.DefaultAuth = mAuthNoop.NewAuth()
		runtime.DefaultRuntime = mRuntimeLocal.NewRuntime()
		//store.DefaultStore = fstore.NewStore()
		store.DefaultStore = mStoreMemory.NewStore()
		config.DefaultConfig, _ = mConfigEnv.NewConfig()
		var err error
		events.DefaultStream, err = mEventStream.NewStream()
		if err != nil {
			logger.Fatalf("Error configuring stream for simple profile: %v", err)
		}

		SetupBroker(mBrokerMemory.NewBroker())
		//turn off Registry
		setupRegistry(mRegistryNoop.NewRegistry())
		// store.DefaultBlobStore, err = fstore.NewBlobStore()
		// if err != nil {
		// 	logger.Fatalf("Error configuring file blob store: %v", err)
		// }

		return nil
	},
}

// setupRegistry configures the registry
func setupRegistry(reg registry.Registry) {
	registry.DefaultRegistry = reg
	router.DefaultRouter = mRouterStatic.NewRouter(router.Registry(reg))
	_ = client.DefaultClient.Init(client.Registry(reg))
	_ = server.DefaultServer.Init(server.Registry(reg))
}
