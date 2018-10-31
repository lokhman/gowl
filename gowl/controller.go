package gowl

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
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

func (c *Controller) JSONResponse(statusCode int, content interface{}) ResponseInterface {
	response := NewResponse(statusCode, func(w io.Writer) error {
		return json.NewEncoder(w).Encode(content)
	})
	response.header.Set("Content-Type", "application/json; charset=utf-8")
	return response
}

func (c *Controller) XMLResponse(statusCode int, content interface{}) ResponseInterface {
	response := NewResponse(statusCode, func(w io.Writer) error {
		return xml.NewEncoder(w).Encode(content)
	})
	response.header.Set("Content-Type", "application/xml; charset=utf-8")
	return response
}

func (c *Controller) TemplateResponse(statusCode int, templateName string, content interface{}) ResponseInterface {
	template := c.server.templates[templateName]
	if template == nil {
		panic(fmt.Sprintf(`gowl: cannot find template with name "%s"`, templateName))
	}
	return c.StreamResponse(statusCode, func(w io.Writer) error {
		return template.Execute(w, content)
	})
}

func (c *Controller) StreamResponse(statusCode int, content func(w io.Writer) error) ResponseInterface {
	return NewResponse(statusCode, content)
}

func (c *Controller) Response(statusCode int, content interface{}) ResponseInterface {
	return NewResponse(statusCode, content)
}

func (c *Controller) OK(content interface{}) ResponseInterface {
	return c.Response(http.StatusOK, content)
}

func (c *Controller) NegotiationOffers(request *Request) []NegotiationOffer {
	return []NegotiationOffer{
		{"text/plain", c.Response},
		{"application/json", c.JSONResponse},
		{"application/xml", c.XMLResponse},
		{"text/xml", c.XMLResponse},
		{"text/html", func(statusCode int, content interface{}) ResponseInterface {
			return c.TemplateResponse(statusCode, request.TemplateName(), content)
		}},
	}
}

func (c *Controller) NegotiateResponse(request *Request, response *Response, offers []NegotiationOffer) ResponseInterface {
	if len(offers) == 0 {
		offers = c.NegotiationOffers(request)
	}

	offerTypes := make([]string, len(offers))
	for i, offer := range offers {
		offerTypes[i] = offer.Type
		i++
	}

	offerType := NegotiateAcceptHeader(request.Header, "Accept", offerTypes)
	if offerType == "" && !c.server.config.NegotiateDefaultOffer {
		link := make(HeaderValues, len(offerTypes))
		for i, offerType := range offerTypes {
			link[i] = HeaderValue{
				Value:  "<" + request.URL.String() + ">",
				Params: StringMap{"type": offerType},
			}
		}

		response := NewErrorResponse(http.StatusNotAcceptable, c.server.config.ServerName, "")
		response.Header().Add("Link", link.String())
		return response
	} else if c.server.config.NegotiateDefaultOffer {
		offerType = offerTypes[0]
	}

	offerResponse := offers[StringInSliceIndex(offerType, offerTypes)].Response
	negotiatedResponse := offerResponse(response.statusCode, response.Content)
	copyHeader(negotiatedResponse.Header(), response.header)
	return negotiatedResponse
}

func (c *Controller) NegotiateResponseListener(event EventInterface) {
	ev := event.(*ResponseEvent)
	if response, ok := ev.Response().(*Response); ok {
		ev.SetResponse(c.NegotiateResponse(ev.Request(), response, nil))
	}
}

func (c *Controller) IndexAction(r *Request) ResponseInterface {
	return c.OK("Welcome!")
}

// NegotiationResponse
type NegotiationResponse func(statusCode int, content interface{}) ResponseInterface

// NegotiationOffer
type NegotiationOffer struct {
	Type     string
	Response NegotiationResponse
}

// ...
func getControllerName(controller ControllerInterface) string {
	name := getTypeName(controller)
	name = strings.TrimLeft(name, "*")
	return strings.TrimPrefix(name, "main.")
}
