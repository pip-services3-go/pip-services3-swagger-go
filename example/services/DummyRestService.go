package example_services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cconv "github.com/pip-services3-go/pip-services3-commons-go/convert"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	crefer "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cvalid "github.com/pip-services3-go/pip-services3-commons-go/validate"
	cservices "github.com/pip-services3-go/pip-services3-rpc-go/services"
	data "github.com/pip-services3-go/pip-services3-swagger-go/example/data"
	logic "github.com/pip-services3-go/pip-services3-swagger-go/example/logic"
	"github.com/rakyll/statik/fs"

	_ "github.com/pip-services3-go/pip-services3-swagger-go/example/resources"
)

type DummyRestService struct {
	cservices.RestService
	controller logic.IDummyController
}

func NewDummyRestService() *DummyRestService {
	c := &DummyRestService{}
	c.RestService = *cservices.InheritRestService(c)
	c.DependencyResolver.Put("controller", crefer.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	return c
}

func (c *DummyRestService) Configure(config *cconf.ConfigParams) {
	c.RestService.Configure(config)
}

func (c *DummyRestService) SetReferences(references crefer.IReferences) {
	c.RestService.SetReferences(references)
	depRes, depErr := c.DependencyResolver.GetOneRequired("controller")
	if depErr == nil && depRes != nil {
		c.controller = depRes.(logic.IDummyController)
	}
}

func (c *DummyRestService) getPageByFilter(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	paginParams := make(map[string]string, 0)

	paginParams["skip"] = params.Get("skip")
	paginParams["take"] = params.Get("take")
	paginParams["total"] = params.Get("total")

	delete(params, "skip")
	delete(params, "take")
	delete(params, "total")

	result, err := c.controller.GetPageByFilter(
		c.GetCorrelationId(req),
		cdata.NewFilterParamsFromValue(params), // W! need test
		cdata.NewPagingParamsFromTuples(paginParams),
	)
	c.SendResult(res, req, result, err)
}

func (c *DummyRestService) getOneById(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	vars := mux.Vars(req)
	dummyId := params.Get("dummy_id")
	if dummyId == "" {
		dummyId = vars["dummy_id"]
	}
	result, err := c.controller.GetOneById(
		c.GetCorrelationId(req),
		dummyId)
	c.SendResult(res, req, result, err)
}

func (c *DummyRestService) create(res http.ResponseWriter, req *http.Request) {
	correlationId := c.GetCorrelationId(req)
	var dummy data.Dummy

	body, bodyErr := ioutil.ReadAll(req.Body)
	if bodyErr != nil {
		err := cerr.NewInternalError(correlationId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(bodyErr)
		c.SendError(res, req, err)
		return
	}
	defer req.Body.Close()
	jsonErr := json.Unmarshal(body, &dummy)

	if jsonErr != nil {
		err := cerr.NewInternalError(correlationId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(jsonErr)
		c.SendError(res, req, err)
		return
	}

	result, err := c.controller.Create(
		correlationId,
		dummy,
	)
	c.SendCreatedResult(res, req, result, err)
}

func (c *DummyRestService) update(res http.ResponseWriter, req *http.Request) {
	correlationId := c.GetCorrelationId(req)

	var dummy data.Dummy

	body, bodyErr := ioutil.ReadAll(req.Body)
	if bodyErr != nil {
		err := cerr.NewInternalError(correlationId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(bodyErr)
		c.SendError(res, req, err)
		return
	}
	defer req.Body.Close()
	jsonErr := json.Unmarshal(body, &dummy)

	if jsonErr != nil {
		err := cerr.NewInternalError(correlationId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(jsonErr)
		c.SendError(res, req, err)
		return
	}
	result, err := c.controller.Update(
		correlationId,
		dummy,
	)
	c.SendResult(res, req, result, err)
}

func (c *DummyRestService) deleteById(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	vars := mux.Vars(req)

	dummyId := params.Get("dummy_id")
	if dummyId == "" {
		dummyId = vars["dummy_id"]
	}

	result, err := c.controller.DeleteById(
		c.GetCorrelationId(req),
		dummyId,
	)
	c.SendDeletedResult(res, req, result, err)
}

func (c *DummyRestService) Register() {
	statikFS, err := fs.NewWithNamespace("example")
	if err != nil {
		panic(err)
	}
	r, err := statikFS.Open("/dummies.yml")
	if err != nil {
		panic(err)
	}
	defer r.Close()
	content, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	c.SwaggerRoute = "/dummies/swagger"
	c.RegisterOpenApiSpec(string(content))

	c.RegisterRoute(
		"get", "/dummies",
		&cvalid.NewObjectSchema().
			WithOptionalProperty("key", cconv.String).
			WithOptionalProperty("skip", cconv.Long).
			WithOptionalProperty("take", cconv.Long).
			WithOptionalProperty("total", cconv.Boolean).Schema,
		c.getPageByFilter,
	)

	c.RegisterRoute(
		"get", "/dummies/{dummy_id}",
		&cvalid.NewObjectSchema().
			WithRequiredProperty("dummy_id", cconv.String).Schema,
		c.getOneById,
	)

	c.RegisterRoute(
		"post", "/dummies",
		&cvalid.NewObjectSchema().
			WithRequiredProperty("body", data.NewDummySchema()).Schema,
		c.create,
	)

	c.RegisterRoute(
		"put", "/dummies",
		&cvalid.NewObjectSchema().
			WithRequiredProperty("body", data.NewDummySchema()).Schema,
		c.update,
	)

	c.RegisterRoute(
		"delete", "/dummies/{dummy_id}",
		&cvalid.NewObjectSchema().
			WithRequiredProperty("dummy_id", cconv.String).Schema,
		c.deleteById,
	)
}
