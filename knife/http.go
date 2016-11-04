package knife

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// RouteParams represents the params for a route
type RouteParams struct {
	Params httprouter.Params
}

// AsString gets the param as an String
func (p RouteParams) AsString(name string) string {
	return p.Params.ByName(name)
}

// AsInt gets the param as an Int
func (p RouteParams) AsInt(name string) int {
	i, err := strconv.Atoi(p.Params.ByName(name))

	if err != nil {
		return 0
	}

	return i
}

// GetRouteParams gets the URL params of the route
func GetRouteParams(request *http.Request) RouteParams {
	var params httprouter.Params

	if ps := context.Get(request, "params"); ps != nil {
		params = ps.(httprouter.Params)
	}

	return RouteParams{params}
}

// HTTPRouterWrapHandler wraps the http.Handler with a httprouter.Handle
// The httprouter.Handle supports better parse of URL params
func HTTPRouterWrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		context.Set(r, "params", ps)
		h.ServeHTTP(w, r)
	}
}

// HTTPMethod represents a http method
type HTTPMethod func(string, httprouter.Handle)

// Router represents the router
type Router struct {
	*httprouter.Router

	routes Routes

	middlewares    []Middleware
	middlewaresMap MiddlewaresMap

	errorHandler ErrorHandler
}

// Routes represents the map de routes
type Routes map[string][]*Route

// NewRouter creates an instance of the *Router
func NewRouter() *Router {
	return &Router{
		Router: httprouter.New(),
		routes: make(Routes),
	}
}

// SetErrorHandler defines the error handler for the router
func (r *Router) SetErrorHandler(e ErrorHandler) {
	r.errorHandler = e
}

// SetMiddlewares defines the middlewares for the router
func (r *Router) SetMiddlewares(m []Middleware) {
	r.middlewares = m
}

// SetMiddlewaresMap defines the middlewares map fot the router
func (r *Router) SetMiddlewaresMap(m MiddlewaresMap) {
	r.middlewaresMap = m
}

// AddRoutes adds one or more routes for the routes
func (r *Router) AddRoutes(g string, newRoutes ...*Route) {
	oldRoutes := r.routes[g]

	x := "^[\\w-]+$"
	matched, err := regexp.MatchString(x, g)

	if err != nil {
		panic(err)
	}

	if !matched {
		m := "knife: group %s does not match to the %s regex."
		panic(fmt.Sprintf(m, g, x))
	}

	for _, nr := range newRoutes {
		for _, or := range oldRoutes {
			if or.Token == nr.Token {
				m := "knife: many registrations for route '%s' on group '%s'."
				panic(fmt.Sprintf(m, nr.Token, g))
			}
		}

		sep := ""
		if !strings.HasPrefix(nr.Path, "/") {
			sep = "/"
		}

		nr.Token = g + "." + nr.Token
		nr.Path = "/" + g + sep + nr.Path
	}

	r.routes[g] = append(oldRoutes, newRoutes...)
}

// GET creates a new route for HTTP GET
func (r *Router) GET(token string, path string, h Handler) *Route {
	return NewRoute(token, r.Router.GET, path, h)
}

// POST creates a new route for HTTP POST
func (r *Router) POST(token string, path string, h Handler) *Route {
	return NewRoute(token, r.Router.POST, path, h)
}

// DELETE creates a new route for HTTP DELETE
func (r *Router) DELETE(token string, path string, h Handler) *Route {
	return NewRoute(token, r.Router.DELETE, path, h)
}

// PUT creates a new route for HTTP PUT
func (r *Router) PUT(token string, path string, h Handler) *Route {
	return NewRoute(token, r.Router.PUT, path, h)
}

// PATCH creates a new route for HTTP PATCH
func (r *Router) PATCH(token string, path string, h Handler) *Route {
	return NewRoute(token, r.Router.PATCH, path, h)
}

// OPTIONS creates a new route for HTTP OPTIONS
func (r *Router) OPTIONS(token string, path string, h Handler) *Route {
	return NewRoute(token, r.Router.OPTIONS, path, h)
}

// HEAD creates a new route for HTTP HEAD
func (r *Router) HEAD(token string, path string, h Handler) *Route {
	return NewRoute(token, r.Router.HEAD, path, h)
}

// Response represents the server response
type Response struct {
	Status int
	Bytes  []byte
}

// NewResponse creates a instance of Response
func NewResponse(s int, b []byte) Response {
	return Response{s, b}
}

// Handler defines a interface to be used by structs
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) (Response, error)
}

// HandlerFunc represents the routes created from a function
type HandlerFunc func(http.ResponseWriter, *http.Request) (Response, error)

func (r *Router) applyErrorHandler(h HandlerFunc) HandlerFunc {
	return r.errorHandler(h)
}

func (r *Router) responseMiddleware(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		response, err := h(w, req)

		status := response.Status

		if status == 0 {
			if err != nil {
				status = http.StatusInternalServerError
			} else {
				status = http.StatusOK
			}
		}

		w.WriteHeader(status)

		if bytes := response.Bytes; len(bytes) > 0 {
			w.Write(response.Bytes)
		}
	}
}

// Start configures all steps for each route.
// The router starts the middlewares chain with the context.ClearHandler.
// It is responsible to clear all data in the request context.
// After, it configures specific middlewares for a route or adds all.
// So, the router adds the error handle as the last handler in the chain.
func (r *Router) Start() *Router {
	for _, routes := range r.routes {
		for _, route := range routes {
			middlewares := r.middlewares

			middlewaresMap := r.middlewaresMap

			chain := alice.New(context.ClearHandler)

			if middlewareTokens, ok := middlewaresMap[route.Token]; ok {
				for _, middlewareToken := range middlewareTokens {
					for _, middleware := range middlewares {
						if middleware.Token == middlewareToken {
							chain = chain.Append(middleware.Constructor)
						}
					}
				}
			} else {
				for _, middleware := range middlewares {
					chain = chain.Append(middleware.Constructor)
				}
			}

			var handler http.HandlerFunc

			if r.errorHandler != nil {
				handler = r.responseMiddleware(
					r.applyErrorHandler(route.Handler))
			} else {
				handler = r.responseMiddleware(route.Handler)
			}

			route.Method(route.Path, HTTPRouterWrapHandler(chain.Then(handler)))
		}
	}

	return r
}

// Route represents a route
type Route struct {
	Token   string
	Method  HTTPMethod
	Path    string
	Handler HandlerFunc
}

// NewRoute creates an instance of *Route using a struct
func NewRoute(token string, method HTTPMethod, path string, handler Handler) *Route {
	return &Route{token, method, path, handler.ServeHTTP}
}

// NewRouteFunc creates an instance of *Route using a func
func NewRouteFunc(token string, method HTTPMethod, path string, handler HandlerFunc) *Route {
	return &Route{token, method, path, handler}
}

// Middleware represents a middleware
type Middleware struct {
	Token       string
	Constructor alice.Constructor
}

// MiddlewaresMap represents the middlewares map.
type MiddlewaresMap map[string][]string

// MiddlewareMapper represents the mapper of middlewares and routes
type MiddlewareMapper struct {
	MiddlewaresMap MiddlewaresMap
}

// NewMiddlewareMapper creates an instance of NewMiddlewareMapper
func NewMiddlewareMapper() *MiddlewareMapper {
	return &MiddlewareMapper{
		MiddlewaresMap: make(MiddlewaresMap),
	}
}

// Map maps middlewares for routes.
// A route can have specifics middlewares and not to use the globals.
func (mapper *MiddlewareMapper) Map(routeToken string, middlewareTokens ...string) {
	mapper.MiddlewaresMap[routeToken] = middlewareTokens
}

// NewMiddlewareManager creates an instance of MiddlewareManager
func NewMiddlewareManager() *MiddlewareManager {
	return &MiddlewareManager{}
}

// MiddlewareManager represents the middleware manager.
// It can stores built-in middlewares or external middlewares.
// The middlewares are supported via alice package.
type MiddlewareManager struct {
	Middlewares []Middleware
}

// Add adds a middleware to the middlewares chain
func (m *MiddlewareManager) Add(token string, constructor alice.Constructor) {
	m.Middlewares = append(m.Middlewares, Middleware{token, constructor})
}

// ErrorHandler represents the global error handler
type ErrorHandler func(HandlerFunc) HandlerFunc

// Error Represents an error sent as response
type Error struct {
	Message []string `json:"errors"`
}

// NewError creates an instance of Error
func NewError(message ...string) *Error {
	return &Error{message}
}

// HTTPError writes for the response object
func HTTPError(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// PanicRecoverMiddleware recovers a panic error
func PanicRecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					HTTPError(w, http.StatusInternalServerError,
						NewError(err.Error()))
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// JSONContentTypeMiddleware forces application/json Content-Type
// It forces only when the http method is PUT or POST
func JSONContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (r.Method == "POST" || r.Method == "PUT") &&
			r.Header.Get("Content-Type") != "application/json" {

			status := http.StatusUnsupportedMediaType

			HTTPError(w, status, NewError(http.StatusText(status)))

			return
		}

		next.ServeHTTP(w, r)
	})
}

// UnMarshalJSONFromReader parses the json from a Reader and stores in v interface
func UnMarshalJSONFromReader(r io.Reader, v interface{}) error {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		format := fmt.Sprintf("It was not possible to read body json. Origin - %s", err.Error())
		return NewUnMarshalError(format)
	}
	return UnMarshalJSON(body, v)
}

// UnMarshalJSON the JSON-encoded data and stores in v interface
func UnMarshalJSON(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		format := fmt.Sprintf("It was not possible to decode json. Origin - %s", err.Error())
		return NewUnMarshalError(format)
	}
	return nil
}

// MarshalJSON returns the JSON encoding of v
func MarshalJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// JSONError creates a json bytes of Error
func JSONError(message ...string) ([]byte, error) {
	return json.Marshal(&Error{message})
}

// UnMarshalError represents an unmarshal error
type UnMarshalError struct {
	Msg string
}

// IsUnMarshalError verifies if error is an UnMarshalError
func IsUnMarshalError(err error) bool {
	_, ok := err.(UnMarshalError)
	return ok
}

// Error gets the error message
func (v UnMarshalError) Error() string {
	return v.Msg
}

// NewUnMarshalError an instance of UnMarshalError
func NewUnMarshalError(msg string) UnMarshalError {
	return UnMarshalError{msg}
}

// GetSimpleNetClient gets an instance of http.Client
func GetSimpleNetClient(secs int) *http.Client {
	if secs == 0 {
		secs = 30
	}

	return &http.Client{
		Timeout: time.Second * time.Duration(secs),
	}
}

// URLParams generates an url with params
// The url is generated from one string format via Sprintf
func URLParams(url string, p ...interface{}) string {
	return fmt.Sprintf(url, p...)
}
