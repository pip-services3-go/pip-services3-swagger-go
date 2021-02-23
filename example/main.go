package main

import (
	"bufio"
	"fmt"
	"os"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	crun "github.com/pip-services3-go/pip-services3-commons-go/run"
	ccount "github.com/pip-services3-go/pip-services3-components-go/count"
	clog "github.com/pip-services3-go/pip-services3-components-go/log"
	cservices "github.com/pip-services3-go/pip-services3-rpc-go/services"
	logic "github.com/pip-services3-go/pip-services3-swagger-go/example/logic"
	services "github.com/pip-services3-go/pip-services3-swagger-go/example/services"
	sservices "github.com/pip-services3-go/pip-services3-swagger-go/services"
)

func main() {
	// Create components
	logger := clog.NewConsoleLogger()
	counter := ccount.NewLogCounters()
	controller := logic.NewDummyController()
	httpEndpoint := cservices.NewHttpEndpoint()
	restService := services.NewDummyRestService()
	httpService := services.NewDummyCommandableHttpService()
	statusService := cservices.NewStatusRestService()
	heartbeatService := cservices.NewHeartbeatRestService()
	swaggerService := sservices.NewSwaggerService()

	components := []interface{}{
		logger,
		counter,
		controller,
		httpEndpoint,
		restService,
		httpService,
		statusService,
		heartbeatService,
		swaggerService,
	}

	// Configure components
	logger.Configure(cconf.NewConfigParamsFromTuples(
		"level", "trace",
	))

	httpEndpoint.Configure(cconf.NewConfigParamsFromTuples(
		"connection.prototol", "http",
		"connection.host", "localhost",
		"connection.port", 8080,
	))

	restService.Configure(cconf.NewConfigParamsFromTuples(
		"swagger.enable", true,
	))

	httpService.Configure(cconf.NewConfigParamsFromTuples(
		"base_route", "dummies2",
		"swagger.enable", true,
	))

	// Set references
	references := cref.NewReferencesFromTuples(
		cref.NewDescriptor("pip-services", "logger", "console", "default", "1.0"), logger,
		cref.NewDescriptor("pip-services", "counter", "log", "default", "1.0"), counter,
		cref.NewDescriptor("pip-services", "endpoint", "http", "default", "1.0"), httpEndpoint,
		cref.NewDescriptor("pip-services-dummies", "controller", "default", "default", "1.0"), controller,
		cref.NewDescriptor("pip-services-dummies", "service", "rest", "default", "1.0"), restService,
		cref.NewDescriptor("pip-services-dummies", "service", "commandable-http", "default", "1.0"), httpService,
		cref.NewDescriptor("pip-services", "status-service", "rest", "default", "1.0"), statusService,
		cref.NewDescriptor("pip-services", "heartbeat-service", "rest", "default", "1.0"), heartbeatService,
		cref.NewDescriptor("pip-services", "swagger-service", "http", "default", "1.0"), swaggerService,
	)

	cref.Referencer.SetReferences(references, components)

	// Open components
	err := crun.Opener.Open("", components)
	if err != nil {
		logger.Error("", err, "Failed to open components")
		return
	}

	// Wait until user presses ENTER
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press ENTER to stop the microservice...")
	reader.ReadString('\n')

	// Close components
	err = crun.Closer.Close("", components)
	if err != nil {
		logger.Error("", err, "Failed to close components")
	}
}
