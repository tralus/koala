package knife

import (
	"fmt"
	"strconv"
	"net/http"
	"encoding/json"
	"time"
	
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/gorilla/context"
)

type RouteParams struct {
	Params httprouter.Params
}

func (p RouteParams) AsString(name string) string {
	return p.Params.ByName(name)
}

func (p RouteParams) AsInt(name string) int {
	i, err := strconv.Atoi(p.Params.ByName(name))
	
	if (err != nil) {
		return 0
	}
	
	return i
}

func GetRouteParams(request *http.Request) RouteParams {
	var params httprouter.Params
	
	if ps := context.Get(request, "params"); ps != nil {
		params = ps.(httprouter.Params)
	}

	return RouteParams{params}
}

func HttpRouterWrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		context.Set(r, "params", ps)
		h.ServeHTTP(w, r)
	}
}

type HttpMethod func(string, httprouter.Handle)

type Router struct {
	*httprouter.Router
	
	Debug bool
	Routes map[string][]*Route
	
	Middlewares []Middleware
	MiddlewaresMap MiddlewaresMap
	
	ErrorHandler ErrorHandler
}

type Routes map[string][]*Route

func NewRouter() *Router {
	return &Router{
		Router: httprouter.New(),
		Routes: make(Routes),
	}
}

func (r *Router) SetErrorHandler(e ErrorHandler) {
	r.ErrorHandler = e
}

func (router *Router) AddRoutes(group string, newRoutes ...*Route) {
	oldRoutes := router.Routes[group]
	
	for _, newRoute := range newRoutes {
		for _, oldRoute := range oldRoutes {
			if oldRoute.Token == newRoute.Token {
				panic("kinife: multiple registrations for route token: " + newRoute.Token)
			}
		}
		
		newRoute.Token = group + "." + newRoute.Token
	}
	
	router.Routes[group] = append(oldRoutes, newRoutes...)
}

type Response struct {
	Status int
	Bytes []byte
}

func NewResponse(s int, b []byte) Response {
	return Response{s, b}
}

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) (Response, error)
}

type HandlerFunc func(http.ResponseWriter, *http.Request) (Response, error)

func (r *Router) applyErrorHandler(h HandlerFunc) HandlerFunc {
	return r.ErrorHandler(h)
}

func (r *Router) responseMiddleware(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		response, err := h(w, req)
		
		status := response.Status
		
		if (status == 0) {
			if (err != nil) {
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

func (router *Router) Start() {
	for _, routes := range router.Routes {
		for _, route := range routes {
			middlewares := router.Middlewares
			
			middlewaresMap := router.MiddlewaresMap
			
			chain := alice.New(context.ClearHandler)
			
			if middlewareTokens, ok := middlewaresMap[route.Token]; ok {
				for _, middlewareToken := range middlewareTokens {
					for _, middleware := range middlewares {
						if (middleware.Token == middlewareToken) {
							chain = chain.Append(middleware.Constructor)
						}
					}
				}
			} else {
				for _, middleware := range middlewares {
					chain = chain.Append(middleware.Constructor)
				}
			}

			handler := router.responseMiddleware(
				router.applyErrorHandler(route.Handler))

			route.Method(route.Path, HttpRouterWrapHandler(chain.Then(handler)))
		}
	}
}

type Route struct {
	Token string
	Method HttpMethod
	Path string
	Handler HandlerFunc
}

func NewRoute(token string, method HttpMethod, path string, handler Handler) *Route {
	return &Route{token, method, path, handler.ServeHTTP}
}

func NewRouteFunc(token string, method HttpMethod, path string, handler HandlerFunc) *Route {
	return &Route{token, method, path, handler}
}

type Middleware struct {
	Token string
	Constructor alice.Constructor
}

type MiddlewaresMap map[string][]string

type MiddlewareMapper struct {
	MiddlewaresMap MiddlewaresMap
}

func NewMiddlewareMapper() *MiddlewareMapper {
	return &MiddlewareMapper{
		MiddlewaresMap: make(MiddlewaresMap),
	}
}

func (mapper *MiddlewareMapper) Map(routeToken string, middlewareTokens ...string) {
	mapper.MiddlewaresMap[routeToken] = middlewareTokens;
}

func NewMiddlewareManager() *MiddlewareManager {
	return &MiddlewareManager{}
}

type MiddlewareManager struct {
	Middlewares []Middleware
}

func (m *MiddlewareManager) Add(token string, constructor alice.Constructor) {
	m.Middlewares = append(m.Middlewares, Middleware{token, constructor})
}

type ErrorHandler func(HandlerFunc) HandlerFunc

func JsonAcceptMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			status := http.StatusNotAcceptable
			
			HttpError(w, status, NewError(status, http.StatusText(status)))
			
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

type Error struct {
	Status int    `json:"status"`
	Message []string `json:"errors"`
}

func NewError(status int, message ...string) *Error {
	return &Error{status, message}
}

func HttpError(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func PanicRecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
	        if r := recover(); r != nil {
	            if err, ok := r.(error); ok {
	                JsonError(http.StatusInternalServerError, err.Error())
	            }
	        }
	    }()

		next.ServeHTTP(w, r)
	})
}

func JsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ((r.Method == "POST" || r.Method == "PUT") && 
			r.Header.Get("Content-Type") != "application/json") {
			
			status := http.StatusUnsupportedMediaType
			
			HttpError(w, status, NewError(status, http.StatusText(status)))
			
			return
		}

		next.ServeHTTP(w, r)
	})
}

func JsonUnMarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if (err != nil) {
		format := fmt.Sprintf("It was not possible to decode json. Origin - %s", err.Error())
		return NewJsonUnMarshalError(format)
	}
	return nil
}

func JsonMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func JsonError(status int, message ...string) ([]byte, error) {
	return json.Marshal(&Error{status, message})
}

type Module struct {
	ContextPath string
}

func (module Module) NewRoute(token string, m HttpMethod, path string, h Handler) *Route {
	contextPath := ""
	
	if (module.ContextPath != "") {
		contextPath = module.ContextPath
	}
	
	return NewRoute(token, m, contextPath + path, h)
}

func SetContextPath(contextPath string, modules ...*Module) {
	for _, module := range modules {
		module.ContextPath = contextPath
	}
}

type JsonUnMarshalError struct {
	Msg string
}

func IsJsonUnMarshalError(err error) bool {
	_, ok := err.(JsonUnMarshalError)
	return ok
}

func (v JsonUnMarshalError) Error() string {
	return v.Msg
}

func NewJsonUnMarshalError(msg string) JsonUnMarshalError {
	return JsonUnMarshalError{msg}
}

func GetSimpleNetClient(secs int) *http.Client {
	return &http.Client{
	  	Timeout: time.Second * time.Duration(secs),
	}
}

func UrlForParams(url string, p ...interface{}) string {
	return fmt.Sprintf(url, p...)
}