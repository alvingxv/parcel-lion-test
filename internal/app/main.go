package app

import (
	"context"
	"lion-parcel-test/config"
	"lion-parcel-test/pkg/httpclient"
	"lion-parcel-test/pkg/log"

	"github.com/opentracing/opentracing-go"
	"go.elastic.co/apm/module/apmot"
)

type App struct {
	Usecases     *Usecases
	Repos        *Repositories
	Dependencies *Dependencies
}

func NewApp(ctx context.Context) (*App, error) {
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	err = log.Initialize()
	if err != nil {
		panic(err)
	}

	httpclient.Init()
	httpclient.Client.NewCbSource(
		httpclient.Client.CbWithCommand("otherservice"),             // Name of the circuit breaker for this service
		httpclient.Client.CbWithErrorPercentThreshold(50),           // Open the circuit if 50% of requests fail
		httpclient.Client.CbWithFallbackMsg("otherservice Timeout"), // Message to return when the circuit is open
		httpclient.Client.CbWithMaxConcurrentRequests(50),           // Allow up to 50 requests at the same time
		httpclient.Client.CbWithSleepWindow(5),                      // Wait 5 seconds before trying to close the circuit again
		httpclient.Client.CbWithTimeout(4000),                       // Timeout for each request is 4 seconds
		httpclient.Client.CbWithRequestVolumeThreshold(100),         // Start checking errors after 100 requests
	)
	opentracing.SetGlobalTracer(apmot.New())

	dependencies, err := NewDependencies()
	if err != nil {
		panic(err)
	}

	repos := NewRepos(dependencies)

	usecases := NewUsecases(repos)

	return &App{
		Repos:        repos,
		Usecases:     usecases,
		Dependencies: dependencies,
	}, nil
}

func (a *App) Close(ctx context.Context) error {
	err := a.Dependencies.Close(ctx)
	if err != nil {
		return err
	}

	return nil
}
