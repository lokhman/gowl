package gowl

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/lokhman/gowl/events"
	"github.com/lokhman/gowl/helpers"
	"github.com/lokhman/gowl/httputil"
	"github.com/lokhman/gowl/templates"
	"github.com/lokhman/gowl/types"
	"github.com/pkg/errors"
)

// ServerInterface
type ServerInterface interface {
	Config() *Config
	NewRouter() RouterInterface
	RegisterRouter(router RouterInterface, routers ...RouterInterface)
	RegisterController(controller ControllerInterface, controllers ...ControllerInterface)
	On(eventType events.EventType, listener func(event EventInterface))
	LoadTemplates()
	Listen() error
	String() string
}

// server
type server struct {
	config    *Config
	router    *compiledRouter
	templates map[string]*template.Template
}

func (s *server) Config() *Config {
	return s.config
}

func (s *server) NewRouter() RouterInterface {
	return NewRouter(s.defaultRouterFlags())
}

func (s *server) RegisterRouter(router RouterInterface, routers ...RouterInterface) {
	s.registerRouters(append([]RouterInterface{router}, routers...))
}

func (s *server) RegisterController(controller ControllerInterface, controllers ...ControllerInterface) {
	s.registerControllers(append([]ControllerInterface{controller}, controllers...))
}

func (s *server) On(eventType events.EventType, listener func(event EventInterface)) {
	s.router.emitter.On(eventType, listener)
}

func (s *server) LoadTemplates() {
	if s.templates != nil {
		return
	}
	funcMap := make(template.FuncMap)
	for name, fn := range s.config.TemplateFunc {
		funcMap[name] = fn
	}
	t, err := templates.Load(s.config.TemplatePath, s.config.TemplateFileExt, funcMap)
	if err != nil {
		panic(fmt.Sprintf("gowl: cannot load templates: %s", err.Error()))
	}
	s.templates = t
}

func (s *server) Listen() error {
	if !s.router.compile() {
		return nil
	}

	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: s,
	}

	if s.config.EnableTLS {
		return server.ListenAndServeTLS(s.config.CertFile, s.config.KeyFile)
	}
	return server.ListenAndServe()
}

func (s *server) String() string {
	buf := new(strings.Builder)
	buf.WriteString(s.config.String())
	buf.WriteByte('\n')
	if len(s.router.routes) > 0 {
		out := make([]string, len(s.router.routes))
		pad := make([]int, len(s.router.routes))
		max := 0

		// calculate padding after method
		for i, route := range s.router.routes {
			str := route.String()
			pos := strings.IndexByte(str, ' ')
			if pos > max {
				max = pos
			}
			out[i] = str
			pad[i] = pos
		}

		// apply padding
		for i, str := range out {
			p := strings.Repeat(" ", max-pad[i]+1)
			out[i] = "  " + strings.Replace(str, " ", p, 1)
		}
		buf.WriteString("Routes:\n")
		buf.WriteString(strings.Join(out, "\n"))
	} else {
		buf.WriteString("Routes: <no routes>")
	}
	return buf.String()
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var start = time.Now()
	var path = r.URL.Path
	var handler Handler

	var request = &Request{Request: r, server: s}
	var response ResponseInterface

	// redirect request to lowercase path if configured
	if s.config.RedirectUpperCasePath && helpers.IndexUpper(path) != -1 {
		response = s.redirect(request, strings.ToLower(path))
		s.serve(w, request, response, start)
		return
	}

	// match request by method and path
	route, params, flag := s.router.match(r.Method, path)

	switch flag {
	case HandleOPTIONS:
		// handle OPTIONS automatically
		s.setAllowHeaderForPath(w, path)
		response = NewResponse(http.StatusOK, nil)
		s.serve(w, request, response, start)
		return
	case HandleMethodNotAllowed:
		s.setAllowHeaderForPath(w, path)

		// if handler not configured, return plain HTTP error
		if handler = s.config.MethodNotAllowedHandler; handler == nil {
			response = s.error(http.StatusMethodNotAllowed, "")
			s.serve(w, request, response, start)
			return
		}
	case RedirectTrailingSlash:
		// fix trailing slash
		response = s.redirect(request, r.URL.Path+"/")
		s.serve(w, request, response, start)
		return
	}

	if route != nil {
		handler = route.handler
	} else {
		// if handler not configured, return plain HTTP error
		if handler = s.config.NotFoundHandler; handler == nil {
			response = s.error(http.StatusNotFound, "")
			s.serve(w, request, response, start)
			return
		}
	}

	// if handler is still not defined
	if handler == nil {
		response = s.error(http.StatusNotImplemented, "")
		s.serve(w, request, response, start)
		return
	}

	// add special parameters
	params.Set(":route", route.name)
	params.Set(":path", route.path)
	request.params = params

	// resolve request scheme
	if scheme := r.Header.Get("X-Scheme"); scheme != "" {
		request.URL.Scheme = strings.ToLower(scheme)
	} else if s.config.EnableTLS {
		request.URL.Scheme = "https"
	} else {
		request.URL.Scheme = "http"
	}

	// TODO: validate host
	// specify host in URL
	if r.URL.Host == "" {
		r.URL.Host = r.Host
	}

	// recover
	defer func() {
		if err := recover(); err != nil {
			var e error
			switch err := err.(type) {
			case string:
				e = errors.New(err)
			case error:
				e = errors.Wrap(err, "gowl")
			default:
				e = errors.New(fmt.Sprintf("%+v", err))
			}

			// emit "panic" events
			if route.emitter.HasListeners(EventPanic) {
				event := &PanicEvent{error: e}
				route.emitter.Emit(EventPanic, event)
			}

			// display error 500 with stack trace
			stack := getMainStackTrace(e.(stackTracer).StackTrace())
			debug := fmt.Sprintf("%s\n%+v", err, stack)
			response = s.error(http.StatusInternalServerError, debug)
			s.serve(w, request, response, start)
		}
	}()

	// emit "request" events
	if route.emitter.HasListeners(EventRequest) {
		event := &RequestEvent{request: request}
		route.emitter.Emit(EventRequest, event)
		response = event.response
	}

	// handle request
	if response == nil {
		response = handler(request)
	}

	// emit "response" events
	if route.emitter.HasListeners(EventResponse) {
		event := &ResponseEvent{request: request, response: response}
		route.emitter.Emit(EventResponse, event)
		response = event.response
	}

	// still no response?
	if response == nil {
		response = EmptyResponse()
	}

	// write response to the connection
	s.serve(w, request, response, start)
}

func (s *server) serve(w http.ResponseWriter, request *Request, response ResponseInterface, start time.Time) {
	defer func() { _ = recover() }()

	if s.config.ServerName != "" {
		w.Header().Set("Server", s.config.ServerName)
	}

	statusCode := response.StatusCode()
	if _, ok := response.(ResponseWriterInterface); !ok {
		httputil.CopyHeader(w.Header(), response.Header())
		w.WriteHeader(statusCode)
	}

	if err := response.Write(w); err != nil {
		Error.Print(err)
	}

	buf := new(strings.Builder)
	fmt.Fprintf(buf, "%3d %s %s", statusCode, request.Method, request.URL.Path)
	if name := request.Param(":route"); name != "" {
		buf.WriteString(" [")
		buf.WriteString(name)
		buf.WriteByte(']')
	}
	fmt.Fprintf(buf, " (%v)", time.Now().Sub(start))
	Debug.Print(buf.String())
}

func (_ *server) redirect(request *Request, url string) ResponseInterface {
	statusCode := http.StatusMovedPermanently
	if request.Method != GET {
		statusCode = http.StatusPermanentRedirect
	}
	return RedirectResponse(request, statusCode, url)
}

func (_ *server) error(statusCode int, debug string) ResponseInterface {
	if !*_debug && debug != "" {
		Error.Print(debug)
		debug = ""
	}
	return ErrorResponse(statusCode, debug)
}

func (s *server) setAllowHeaderForPath(w http.ResponseWriter, path string) {
	if allow := s.router.allowedMethods(path); len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
	}
}

func (s *server) registerRouters(routers []RouterInterface) {
	for _, router := range routers {
		s.router.addRouter(router)
	}
}

func (s *server) registerControllers(controllers []ControllerInterface) {
	flags := s.defaultRouterFlags()
	for _, controller := range controllers {
		if name := controller.Name(); name != "" {
			panic(fmt.Sprintf(`gowl: controller with name "%s" is already registered`, name))
		}

		controller.init(getControllerName(controller), s)

		router := NewRouter(flags)
		controller.Routing(router)
		s.router.addRouter(router)
	}
}

func (s *server) defaultRouterFlags() (flags types.Flag) {
	if s.config.HandleOptions {
		flags.Set(HandleOPTIONS)
	}
	if s.config.HandleMethodNotAllowed {
		flags.Set(HandleMethodNotAllowed)
	}
	if s.config.RedirectTrailingSlash {
		flags.Set(RedirectTrailingSlash)
	}
	return
}

func NewServer(config *Config) ServerInterface {
	return &server{
		config: config,
		router: newCompiledRouter(),
	}
}
