package gowl

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/lokhman/gowl/events"
	"github.com/lokhman/gowl/helpers"
	"github.com/lokhman/gowl/types"
)

const (
	defaultState types.Flag = 1 << iota

	HandleOPTIONS
	HandleMethodNotAllowed
	RedirectTrailingSlash
)

const RouteParamRequirement = `[^/]+`

var reRoutePathParams = regexp.MustCompile(`{([a-z0-9_]+)(?:<(.*?)>)?(?:\?([^}]+))?}`)
var reRouteParamRequirement = regexp.MustCompile(`^` + RouteParamRequirement + `$`)

// Handler
type Handler func(r *Request) ResponseInterface

// ParamAttributes
type ParamAttributes struct {
	Requirement  string
	DefaultValue string

	reRequirement *regexp.Regexp
	compiled      bool
}

// Params
type Params map[string]ParamAttributes

// RouteInterface
type RouteInterface interface {
	SetName(name string) RouteInterface
	AddParam(name string, attr ParamAttributes) RouteInterface
	SetParams(params Params) RouteInterface
	SetFlag(flag types.Flag) RouteInterface
	On(eventType events.EventType, listener func(event EventInterface)) RouteInterface
	String() string

	compile() bool
}

// route
type route struct {
	emitter events.Emitter

	name    string
	methods []string
	path    string
	handler Handler
	params  Params
	flags   types.Flag

	rePath   *regexp.Regexp
	rePrefix string

	trailingSlash bool
	compiled      bool
}

func (r *route) SetName(name string) RouteInterface {
	r.name = name
	return r
}

func (r *route) AddParam(name string, attr ParamAttributes) RouteInterface {
	if _, ok := r.params[name]; ok {
		panic(fmt.Sprintf(`gowl: path "%s" has a duplicated parameter "%s"`, r.path, name))
	}
	r.params[name] = attr
	return r
}

func (r *route) SetParams(params Params) RouteInterface {
	r.params = make(Params)
	for name, attr := range params {
		r.AddParam(name, attr)
	}
	return r
}

func (r *route) SetFlag(flag types.Flag) RouteInterface {
	if r.flags.Has(defaultState) {
		r.flags = 0
	}
	r.flags.Set(flag)
	return r
}

func (r *route) On(eventType events.EventType, listener func(event EventInterface)) RouteInterface {
	r.emitter.On(eventType, listener)
	return r
}

func (r *route) String() string {
	method := "[*]"
	if len(r.methods) == 1 {
		method = r.methods[0]
	} else if len(r.methods) > 1 {
		method = "[" + strings.Join(r.methods, "|") + "]"
	}
	return fmt.Sprintf("%s %s (%s)", method, r.path, r.name)
}

func (r *route) compile() bool {
	if r.compiled {
		return false
	}

	if r.name == "" {
		name := helpers.GetFuncName(r.handler)
		name = strings.TrimPrefix(name, "main.")
		name = strings.TrimSuffix(name, "-fm")
		name = helpers.ToUnderscore(name)

		// trim special characters and suffixes
		for i, s := range strings.Split(name, ".") {
			if i != 0 {
				r.name += "."
			}
			s = strings.Trim(s, "(*_)")
			s = strings.TrimSuffix(s, "_controller")
			s = strings.TrimSuffix(s, "_handler")
			s = strings.TrimSuffix(s, "_action")
			r.name += s
		}
	}

	for _, method := range r.methods {
		assertMethod(method, r.path)
	}

	var err error
	path, rePath, lastIndex, paramCount := r.path, "", 0, 0
	if strings.IndexByte(path, '{') != strings.IndexByte(path, '}') { // -1 != -1
		r.path = helpers.ReplaceAllStringSubmatchFunc(reRoutePathParams, path, func(m []string, i int) string {
			attr, ok := r.params[m[1]]
			if ok && attr.compiled {
				panic(fmt.Sprintf(`gowl: path "%s" has a duplicated parameter "%s"`, path, m[1]))
			} else if !ok || attr.Requirement == "" {
				attr.Requirement = RouteParamRequirement
			}

			if m[2] != "" {
				// extract from path
				attr.Requirement = m[2]
			}
			if attr.Requirement == RouteParamRequirement {
				attr.reRequirement = reRouteParamRequirement
			} else { // compile parameter requirement if not default given
				if attr.reRequirement, err = regexp.Compile(`^` + attr.Requirement + `$`); err != nil {
					panic(fmt.Sprintf(`gowl: path "%s" has invalid parameter "%s": %s`, path, m[1], err.Error()))
				}
			}

			if m[3] != "" {
				// extract from path
				attr.DefaultValue = m[3]
			}
			if attr.DefaultValue != "" && !attr.reRequirement.MatchString(attr.DefaultValue) {
				panic(fmt.Sprintf(`gowl: path "%s" has invalid default value "%s"`, path, attr.DefaultValue))
			}

			attr.compiled = true
			r.params[m[1]] = attr

			if paramCount == 0 {
				r.rePrefix = path[:strings.Index(path, m[0])]
			}

			rePath += regexp.QuoteMeta(path[lastIndex:i])
			rePath += "(?P<" + m[1] + ">" + attr.Requirement + ")"

			lastIndex = i + len(m[0])
			paramCount++

			return "{" + m[1] + "}"
		})
	}

	if n := len(r.params) - paramCount; n != 0 {
		panic(fmt.Sprintf(`gowl: path "%s" has %d unused parameter(s)`, path, n))
	}

	assertPath(r.path)

	if paramCount > 0 {
		rePath += regexp.QuoteMeta(path[lastIndex:])
		if r.rePath, err = regexp.Compile("^" + rePath + "$"); err != nil {
			panic(fmt.Sprintf(`gowl: path "%s" cannot be compiled: %s`, path, err.Error()))
		}
	}

	if r.path[len(r.path)-1] == '/' {
		r.trailingSlash = true
	}

	// compile flags
	r.flags.Clear(defaultState)
	r.compiled = true
	return true
}

func newRoute(methods []string, path string, handler Handler) *route {
	return &route{
		emitter: make(events.Emitter),
		methods: methods,
		path:    path,
		handler: handler,
		params:  make(map[string]ParamAttributes),
		flags:   defaultState,
	}
}

// RouterInterface
type RouterInterface interface {
	SetPrefix(path string)
	SetFlag(flag types.Flag)

	Match(path string, handler Handler, method ...string) RouteInterface
	HEAD(path string, handler Handler) RouteInterface
	GET(path string, handler Handler) RouteInterface
	POST(path string, handler Handler) RouteInterface
	PUT(path string, handler Handler) RouteInterface
	PATCH(path string, handler Handler) RouteInterface
	DELETE(path string, handler Handler) RouteInterface
	OPTIONS(path string, handler Handler) RouteInterface
	TRACE(path string, handler Handler) RouteInterface
	CONNECT(path string, handler Handler) RouteInterface

	On(eventType events.EventType, listener func(event EventInterface))

	compile() (routes []*route, ok bool)
}

// router
type router struct {
	emitter  events.Emitter
	routes   []*route
	prefix   string
	flags    types.Flag
	compiled bool
}

func (r *router) SetPrefix(path string) {
	assertPath(path)
	r.prefix = path
}

func (r *router) SetFlag(flag types.Flag) {
	if r.flags.Has(defaultState) {
		r.flags = 0
	}
	r.flags.Set(flag)
}

func (r *router) Match(path string, handler Handler, method ...string) RouteInterface {
	route := newRoute(method, path, handler)
	r.routes = append(r.routes, route)
	return route
}

func (r *router) HEAD(path string, handler Handler) RouteInterface {
	return r.Match(path, handler, HEAD)
}

func (r *router) GET(path string, handler Handler) RouteInterface {
	return r.Match(path, handler, GET)
}

func (r *router) POST(path string, handler Handler) RouteInterface {
	return r.Match(path, handler, POST)
}

func (r *router) PUT(path string, handler Handler) RouteInterface {
	return r.Match(path, handler, PUT)
}

func (r *router) PATCH(path string, handler Handler) RouteInterface {
	return r.Match(path, handler, PATCH)
}

func (r *router) DELETE(path string, handler Handler) RouteInterface {
	return r.Match(path, handler, DELETE)
}

func (r *router) OPTIONS(path string, handler Handler) RouteInterface {
	return r.Match(path, handler, OPTIONS)
}

func (r *router) TRACE(path string, handler Handler) RouteInterface {
	return r.Match(path, handler, TRACE)
}

func (r *router) CONNECT(path string, handler Handler) RouteInterface {
	return r.Match(path, handler, CONNECT)
}

func (r *router) On(eventType events.EventType, listener func(event EventInterface)) {
	r.emitter.On(eventType, listener)
}

func (r *router) compile() (routes []*route, ok bool) {
	routes = r.routes
	if r.compiled {
		return
	}

	for _, route := range routes {
		if n := len(r.prefix); n > 1 {
			path := route.path
			route.path = r.prefix
			if r.prefix[n-1] == '/' {
				route.path += path[1:]
			} else if path != "/" {
				route.path += path
			}
		}

		// inherit flags if not set
		if route.flags.Has(defaultState) {
			route.flags = r.flags
		}

		// bind events from router emitter
		for eventType, listeners := range r.emitter {
			for _, listener := range listeners {
				route.emitter.On(eventType, listener)
			}
		}
		route.compile()
	}

	ok = true
	return
}

func NewRouter(flags types.Flag) RouterInterface {
	return &router{
		emitter: make(events.Emitter),
		routes:  make([]*route, 0),
		flags:   flags | defaultState,
	}
}

// compiledRouter
type compiledRouter struct {
	routes   []*route
	names    map[string]int
	emitter  events.Emitter
	compiled bool
}

func (r *compiledRouter) normalizeName(name string) string {
	if n, ok := r.names[name]; ok {
		n++
		r.names[name] = n
		name += "_" + strconv.Itoa(n)
		return r.normalizeName(name)
	} else {
		r.names[name] = 0
	}
	return name
}

func (r *compiledRouter) addRouter(router RouterInterface) {
	if routes, ok := router.compile(); ok && len(routes) > 0 {
		for _, route := range routes {
			route.name = r.normalizeName(route.name)
		}
		r.routes = append(r.routes, routes...)
	}
}

func (r *compiledRouter) allowedMethods(path string) (methods []string) {
	methodSet := make(map[string]struct{})
	for _, route := range r.routes {
		if route.rePath == nil {
			// check static path
			if path != route.path {
				continue
			}
		} else {
			// check regex prefix
			if !strings.HasPrefix(path, route.rePrefix) {
				continue
			}

			// match expensive regex
			if !route.rePath.MatchString(path) {
				continue
			}
		}

		// can match any method
		if len(route.methods) == 0 {
			return make([]string, 0)
		}

		for _, method := range route.methods {
			methodSet[method] = struct{}{}
		}

		if route.flags.Has(HandleOPTIONS) {
			methodSet[OPTIONS] = struct{}{}
		}
	}

	if i, n := 0, len(methodSet); n > 0 {
		methods = make([]string, n)
		for method := range methodSet {
			methods[i] = method
			i++
		}
	}
	return
}

func (r *compiledRouter) match(method, path string) (*route, types.StringMap, types.Flag) {
	var pathTrailingSlash = path[len(path)-1] == '/'
	var methodNotAllowed bool
	var flag types.Flag

	for _, route := range r.routes {
		var p = path
		var match []string
		var redirectTrailingSlash bool

		// path is used unchanged for CONNECT requests
		if method != CONNECT && route.flags.Has(RedirectTrailingSlash) {
			// handle redirect for "/test" -> "/test/"
			if route.trailingSlash && !pathTrailingSlash {
				redirectTrailingSlash = true
				p += "/"
			}
		}

		if route.rePath == nil {
			// check static path
			if p != route.path {
				continue // fail
			}
		} else {
			// check regex prefix
			if !strings.HasPrefix(p, route.rePrefix) {
				continue // fail
			}

			// match expensive regex
			match = route.rePath.FindStringSubmatch(p)
			if match == nil {
				continue // fail
			}
		}

		// match if any method allowed or method is explicitly defined
		if len(route.methods) == 0 || helpers.IndexString(method, route.methods) != -1 {
			if redirectTrailingSlash {
				return nil, nil, RedirectTrailingSlash
			}

			// extract parameters from path
			params := make(types.StringMap)
			if match != nil {
				subexpNames := route.rePath.SubexpNames()
				for i, name := range subexpNames[1:] {
					value := match[i+1]
					if value == "" {
						attr := route.params[name]
						value = attr.DefaultValue
					}
					params.Set(name, value)
				}
			}
			return route, params, flag
		}

		if method == OPTIONS && route.flags.Has(HandleOPTIONS) {
			return nil, nil, HandleOPTIONS
		} else if route.flags.Has(HandleMethodNotAllowed) {
			methodNotAllowed = true
		}
	}

	// if method was not found
	if methodNotAllowed {
		return nil, nil, HandleMethodNotAllowed
	}

	return nil, nil, flag
}

func (r *compiledRouter) find(name string) *route {
	for _, route := range r.routes {
		if route.name == name {
			return route
		}
	}
	return nil
}

func (r *compiledRouter) url(name string, params types.StringMap) *url.URL {
	route := r.find(name)
	if route == nil {
		panic(fmt.Sprintf(`gowl: cannot find route with name "%s"`, name))
	}

	q := make(url.Values)
	for name, value := range params {
		q.Set(name, value)
	}

	path := route.path
	if strings.IndexByte(path, '{') != strings.IndexByte(path, '}') { // -1 != -1
		path = helpers.ReplaceAllStringSubmatchFunc(reRoutePathParams, path, func(m []string, _ int) string {
			value, ok := params[m[1]]
			if !ok { // if parameter not given, try to pick default
				if value = route.params[m[1]].DefaultValue; value == "" {
					panic(fmt.Sprintf(`gowl: parameter "%s" is missing in path "%s"`, m[1], name))
				}
			} else if !route.params[m[1]].reRequirement.MatchString(value) {
				panic(fmt.Sprintf(`gowl: parameter "%s" in path "%s" has invalid value "%s"`, m[1], name, value))
			}
			q.Del(m[1])
			return value
		})
	}
	return &url.URL{
		Path:     path,
		RawQuery: q.Encode(),
	}
}

func (r *compiledRouter) compile() bool {
	if r.compiled {
		return false
	}

	for _, route := range r.routes {
		// bind events from compiledRouter emitter
		for eventType, listeners := range r.emitter {
			for _, listener := range listeners {
				route.emitter.On(eventType, listener)
			}
		}
	}

	r.compiled = true
	return true
}

func newCompiledRouter() *compiledRouter {
	return &compiledRouter{
		routes:  make([]*route, 0),
		names:   make(map[string]int),
		emitter: make(events.Emitter),
	}
}

// ...
func assertMethod(method string, path string) {
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			panic(fmt.Sprintf(`gowl: invalid method "%s" for path "%s"`, method, path))
		}
	}
}

func assertPath(path string) {
	if path == "" || path[0] != '/' {
		panic(fmt.Sprintf(`gowl: path "%s" must begin with "/"`, path))
	}

	n := len(path)

	i := 1
	for i < n {
		if path[i] == '\\' {
			panic(fmt.Sprintf(`gowl: path "%s" contains invalid separator "\\"`, path))
		}
		if path[i-1] != '/' {
			goto next
		}
		switch {
		case path[i] == '/':
			panic(fmt.Sprintf(`gowl: path "%s" must not contain empty element`, path))
		case path[i] == '.' && (i+1 == n || path[i+1] == '/'):
			panic(fmt.Sprintf(`gowl: path "%s" must not contain "." element`, path))
		case path[i] == '.' && path[i+1] == '.' && (i+2 == n || path[i+2] == '/'):
			panic(fmt.Sprintf(`gowl: path "%s" must not contain ".." element`, path))
		}
	next:
		i++
	}
}
