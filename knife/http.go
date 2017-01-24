package knife

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// SupressError sets if http error should be sent
var SupressError bool

func init() {
	flag.BoolVar(&SupressError, "koala_knife_supress_error", false, "suppress error")
}

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

	routes         Routes
	middlewares    []Middleware
	middlewaresMap MiddlewaresMap
	errorHandler   ErrorHandler
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

// Request represents the server request
type Request struct {
	httpRequest *http.Request
}

// Body represents the request body
type Body struct {
	reader io.Reader
}

// NewBody creates a Body instance
func NewBody(r io.Reader) Body {
	return Body{r}
}

// UnMarshalJSON parses the JSON into the body
func (b Body) UnMarshalJSON(v interface{}) error {
	body, err := ioutil.ReadAll(b.reader)
	if err != nil {
		format := fmt.Sprintf("It was not possible to read body json. Origin - %s", err.Error())
		return NewUnMarshalError(format)
	}
	return UnMarshalJSON(body, v)
}

// Body gets the request body
func (r Request) Body() Body {
	return NewBody(r.httpRequest.Body)
}

// URL gets the request URL
func (r Request) URL() *url.URL {
	return r.httpRequest.URL
}

// HTTPRequest gests the pointer of *http.Request
func (r Request) HTTPRequest() *http.Request {
	return r.httpRequest
}

// Params gets the request params
func (r Request) Params() RouteParams {
	var params httprouter.Params

	if ps := context.Get(r.HTTPRequest(), "params"); ps != nil {
		params = ps.(httprouter.Params)
	}

	return RouteParams{params}
}

// NewRequest creates an instance of Request
func NewRequest(r *http.Request) *Request {
	return &Request{r}
}

// Response represents the server response
type Response struct {
	writer      http.ResponseWriter
	contentType string
	status      int
	bytes       []byte
}

// NewResponse creates an instance of Response
func NewResponse(w http.ResponseWriter) Response {
	return Response{contentType: "text/html", writer: w}
}

// Writer gets the http response writer
func (r Response) Writer() http.ResponseWriter {
	return r.writer
}

// SetContentType sets the response content type
func (r *Response) SetContentType(s string) {
	r.contentType = s
}

// ContentType gets the response content type
func (r Response) ContentType() string {
	return r.contentType
}

// SetBytes sets the response bytes
func (r *Response) SetBytes(b []byte) {
	r.bytes = b
}

// Bytes gets the response bytes
func (r Response) Bytes() []byte {
	return r.bytes
}

// SetStatus sets the response status
func (r *Response) SetStatus(s int) {
	r.status = s
}

// Status gets the response status
func (r Response) Status() int {
	return r.status
}

// Ok creates an ok response
func (r Response) Ok(bytes []byte) (Response, error) {
	r.SetBytes(bytes)
	return r, nil
}

// JSON creates a JSON response from v
func (r Response) JSON(v interface{}) (Response, error) {
	bytes, err := MarshalJSON(v)

	if err != nil {
		return r, err
	}

	r.SetContentType("application/json")

	return r.Ok(bytes)
}

// NoContent creates a NoContent response
func (r Response) NoContent() (Response, error) {
	r.SetStatus(http.StatusNoContent)
	return r, nil
}

// NotFound creates a NotFound response
func (r Response) NotFound() (Response, error) {
	r.SetStatus(http.StatusNotFound)
	return r, nil
}

// ServerError creates an InternalServerError response
func (r Response) ServerError(err error) (Response, error) {
	r.SetStatus(http.StatusInternalServerError)
	return r, err
}

// BadRequest creates a BadRequest response
func BadRequest(err error) (Response, error) {
	r := Response{}

	r.SetStatus(http.StatusBadRequest)

	return r, err
}

// Handler defines an interface to be used by structs
type Handler interface {
	ServeHTTP(Response, *Request) (Response, error)
}

// HandlerFunc represents the routes created from a function
type HandlerFunc func(Response, *Request) (Response, error)

func (r *Router) applyErrorHandler(h HandlerFunc) HandlerFunc {
	return r.errorHandler(h)
}

func (r *Router) responseMiddleware(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp, err := h(NewResponse(w), NewRequest(req))

		s := resp.Status()

		// Ensures that the Internal Server can be defined without response body
		if (s == 0 && err != nil) || s == http.StatusInternalServerError {
			if err != nil && SupressError == false {
				bytes := []byte(err.Error())
				resp.SetBytes(bytes)
			}
		} else if s == 0 {
			s = http.StatusOK
		}

		w.Header().Set("Content-Type", resp.ContentType())
		w.WriteHeader(s)

		if bytes := resp.Bytes(); len(bytes) > 0 {
			w.Write(resp.Bytes())
		}
	}
}

// Start configures all necessary steps for each route.
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
					if middleware.Silent == false {
						chain = chain.Append(middleware.Constructor)
					}
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
	Silent      bool
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
	m.Middlewares = append(m.Middlewares, Middleware{token, constructor, false})
}

// AddSilent adds a silent middleware to the middlewares chain
func (m *MiddlewareManager) AddSilent(token string, constructor alice.Constructor) {
	m.Middlewares = append(m.Middlewares, Middleware{token, constructor, true})
}

// ErrorHandler represents the global error handler
type ErrorHandler func(HandlerFunc) HandlerFunc

// ErrorMessage represents errors sent as response
type ErrorMessage struct {
	Message []string `json:"errors"`
}

// NewErrorMessage creates a ErrorMessages instance
func NewErrorMessage(message ...string) ErrorMessage {
	return ErrorMessage{message}
}

// PanicRecoverMiddleware recovers a panic error
func PanicRecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				if _, ok := r.(error); ok {
					w.WriteHeader(http.StatusInternalServerError)

					if SupressError == false {
						w.Write(debug.Stack())
					}
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// JSONContentTypeMiddleware forces application/json Content-Type
func JSONContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mediatype, params, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
		charset, ok := params["charset"]

		if !ok {
			charset = "UTF-8"
		}

		if r.ContentLength > 0 &&
			!(mediatype == "application/json" && strings.ToUpper(charset) == "UTF-8") {

			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("Bad Content-Type or charset, expected 'application/json' and 'UTF-8'."))

			return
		}

		next.ServeHTTP(w, r)
	})
}

// UnMarshalJSON parses the JSON-encoded data and stores in v interface
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
