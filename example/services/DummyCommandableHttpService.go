package example_services

import (
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cservices "github.com/pip-services3-go/pip-services3-rpc-go/services"
)

type DummyCommandableHttpService struct {
	*cservices.CommandableHttpService
}

func NewDummyCommandableHttpService() *DummyCommandableHttpService {
	c := DummyCommandableHttpService{
		CommandableHttpService: cservices.NewCommandableHttpService("dummies2"),
	}
	c.DependencyResolver.Put("controller", cref.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	return &c
}
