package example_logic

import (
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	data "github.com/pip-services3-go/pip-services3-swagger-go/example/data"
)

type IDummyController interface {
	GetPageByFilter(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *data.DummyDataPage, err error)
	GetOneById(correlationId string, id string) (result *data.Dummy, err error)
	Create(correlationId string, entity data.Dummy) (result *data.Dummy, err error)
	Update(correlationId string, entity data.Dummy) (result *data.Dummy, err error)
	DeleteById(correlationId string, id string) (result *data.Dummy, err error)
}
