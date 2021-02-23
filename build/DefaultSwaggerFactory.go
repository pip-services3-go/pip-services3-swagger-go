package build

import (
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cbuild "github.com/pip-services3-go/pip-services3-components-go/build"
	"github.com/pip-services3-go/pip-services3-swagger-go/services"
)

// DefaultSwaggerFactory are creates RPC components by their descriptors.

// See Factory
// See HttpEndpoint
// See HeartbeatRestService
// See StatusRestService
type DefaultSwaggerFactory struct {
	cbuild.Factory
	Descriptor               *cref.Descriptor
	SwaggerServiceDescriptor *cref.Descriptor
}

// NewDefaultSwaggerFactorymethod create a new instance of the factory.
func NewDefaultSwaggerFactory() *DefaultSwaggerFactory {
	c := DefaultSwaggerFactory{}
	c.Factory = *cbuild.NewFactory()
	c.Descriptor = cref.NewDescriptor("pip-services", "factory", "swagger", "default", "1.0")
	c.SwaggerServiceDescriptor = cref.NewDescriptor("pip-services", "swagger-service", "*", "*", "1.0")

	c.RegisterType(c.SwaggerServiceDescriptor, services.NewSwaggerService)
	return &c
}
