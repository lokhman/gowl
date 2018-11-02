package gowl

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

// ControllerInterface
type ControllerInterface interface {
	init(name string, server *server)

	Name() string
	Routing(r RouterInterface)
}

// Controller
type Controller struct {
	name   string
	server *server
}

func (c *Controller) init(name string, server *server) {
	c.name = name
	c.server = server
}

func (c *Controller) Name() string {
	return c.name
}

func (c *Controller) Routing(r RouterInterface) {
	r.GET("/", c.IndexAction)
}

func (c *Controller) TemplateResponse(statusCode int, templateName string, content interface{}) ResponseInterface {
	template := c.server.templates[templateName]
	if template == nil {
		panic(fmt.Sprintf(`gowl: cannot find template with name "%s"`, templateName))
	}
	return c.Response(statusCode, func(w io.Writer) error {
		return template.Execute(w, content)
	})
}

func (c *Controller) Response(statusCode int, content interface{}) ResponseInterface {
	return NewResponse(statusCode, content)
}

func (c *Controller) OK(content interface{}) ResponseInterface {
	return c.Response(http.StatusOK, content)
}

func (c *Controller) IndexAction(r *Request) ResponseInterface {
	return c.OK("Welcome!")
}

// ...
func getControllerName(controller ControllerInterface) string {
	name := getTypeName(controller)
	name = strings.TrimLeft(name, "*")
	return strings.TrimPrefix(name, "main.")
}

func getTypeName(i interface{}) string {
	return reflect.TypeOf(i).String()
}
