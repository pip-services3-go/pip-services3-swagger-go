package services

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	cconv "github.com/pip-services3-go/pip-services3-commons-go/convert"
	cservices "github.com/pip-services3-go/pip-services3-rpc-go/services"
	"github.com/rakyll/statik/fs"

	_ "github.com/pip-services3-go/pip-services3-swagger-go/resources"
)

type SwaggerService struct {
	*cservices.RestService
	routes map[string]string
	fs     http.FileSystem
}

func NewSwaggerService() *SwaggerService {
	c := &SwaggerService{}
	c.RestService = cservices.InheritRestService(c)
	c.BaseRoute = "swagger"
	c.routes = map[string]string{}

	sfs, err := fs.NewWithNamespace("swagger")
	if err != nil {
		panic(err)
	}
	c.fs = sfs

	// c.RegisterOpenApiSpec("dummies", "/dummies/swagger")
	// c.RegisterOpenApiSpec("dummies2", "/dummies2/swagger")

	return c
}

func (c *SwaggerService) calculateContentType(fileName string) string {
	if strings.HasSuffix(fileName, ".html") {
		return "text/html"
	}
	if strings.HasSuffix(fileName, ".css") {
		return "text/css"
	}
	if strings.HasSuffix(fileName, ".js") {
		return "application/javascript"
	}
	if strings.HasSuffix(fileName, ".png") {
		return "image/png"
	}
	return "text/plain"
}

func (c *SwaggerService) getSwaggerFile(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	fileName := strings.ToLower(vars["file_name"])
	r, err := c.fs.Open("/" + fileName)
	if err != nil {
		res.WriteHeader(404)
		io.WriteString(res, err.Error())
		return
	}
	defer r.Close()
	content, err := ioutil.ReadAll(r)
	if err != nil {
		res.WriteHeader(500)
		io.WriteString(res, err.Error())
		return
	}

	res.Header().Set("Content-Length", cconv.StringConverter.ToString(len(content)))
	res.Header().Set("Content-Type", c.calculateContentType(fileName))
	res.WriteHeader(200)
	res.Write(content)
}

func (c *SwaggerService) getIndex(res http.ResponseWriter, req *http.Request) {
	r, err := c.fs.Open("/index.html")
	if err != nil {
		panic(err)
	}
	defer r.Close()
	contentBytes, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	content := string(contentBytes)

	builder := strings.Builder{}
	builder.WriteString("[")
	for k, v := range c.routes {
		if builder.Len() > 1 {
			builder.WriteString(",")
		}
		builder.WriteString("{name:\"")
		name := strings.ReplaceAll(k, "\"", "\\\"")
		builder.WriteString(name)
		builder.WriteString("\",url:\"")
		url := strings.ReplaceAll(v, "\"", "\\\"")
		builder.WriteString(url)
		builder.WriteString("\"}")
	}
	builder.WriteString("]")

	content = strings.ReplaceAll(content, "[/*urls*/]", builder.String())

	res.Header().Add("Content-Type", "text/html")
	res.Header().Add("Content-Length", cconv.StringConverter.ToString(len(content)))
	res.WriteHeader(200)
	io.WriteString(res, content)
}

func (c *SwaggerService) redirectToIndex(res http.ResponseWriter, req *http.Request) {
	url := req.RequestURI
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	http.Redirect(res, req, url+"index.html", http.StatusSeeOther)
}

func (c *SwaggerService) composeSwaggerRoute(baseRoute string, route string) string {
	if baseRoute != "" {
		if route == "" {
			route = "/"
		}
		if !strings.HasPrefix(route, "/") {
			route = "/" + route
		}
		if !strings.HasPrefix(baseRoute, "/") {
			baseRoute = "/" + baseRoute
		}
		route = baseRoute + route
	}

	return route
}

func (c *SwaggerService) RegisterOpenApiSpec(baseRoute string, swaggerRoute string) {
	route := c.composeSwaggerRoute(baseRoute, swaggerRoute)
	if baseRoute == "" {
		baseRoute = "default"
	}
	c.routes[baseRoute] = route
}

func (c *SwaggerService) Register() {
	// A hack to redirect default base route
	baseRoute := c.BaseRoute
	c.BaseRoute = ""
	c.RegisterRoute(
		"get", baseRoute, nil, c.redirectToIndex,
	)
	c.BaseRoute = baseRoute

	c.RegisterRoute(
		"get", "/", nil, c.redirectToIndex,
	)

	c.RegisterRoute(
		"get", "/index.html", nil, c.getIndex,
	)

	c.RegisterRoute(
		"get", "/{file_name}", nil, c.getSwaggerFile,
	)
}
