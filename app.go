package koala

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/cors"
	"github.com/tralus/koala/config"
	"github.com/tralus/koala/knife"
)

// Config holds the app config
var Config config.Config

func init() {
	var err error

	// Loads the application config
	Config, err = config.LoadConfig()

	// (Config == config.Config{}) avoids compiler error
	if err != nil && (Config == config.Config{}) {
		panic(err.Error())
	}
}

// Module defines the interface for modules.
type Module interface {
	Up()
}

// App represents a Koala Application
type App struct {
	cors    *cors.Cors
	router  *knife.Router
	modules []Module
}

// NewApplication creates an instance of App
func NewApplication(r *knife.Router, cors *cors.Cors) App {
	return App{cors: cors, router: r}
}

// SetRouter sets the application router
func (a *App) SetRouter(r *knife.Router) {
	a.router = r
}

// SetCors sets CORS support for the application
func (a *App) SetCors(c *cors.Cors) {
	a.cors = c
}

// AddModules adds the modules for the application
func (a *App) AddModules(m []Module) {
	for _, z := range m {
		a.modules = append(a.modules, z)
	}
}

// AddModule adds a module for the application
func (a *App) AddModule(m Module) {
	a.modules = append(a.modules, m)
}

// Gets the separator for the output
func sep() string {
	return "***"
}

// Run starts the application
func (a *App) Run() error {
	// The router can not be nil
	if a.router == nil {
		panic("Set a not nil router to run.")
	}

	fmt.Println("Starting app...")

	// Up application modules
	for _, m := range a.modules {
		m.Up()
	}

	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "9003"
	}

	fmt.Println(sep())
	fmt.Printf("Debug: %t\n", Config.Debug)
	fmt.Printf("Port: %s\n", port)

	fmt.Println(sep())
	fmt.Printf("On http://localhost:%s\n", port)
	fmt.Println("To shut down, press <CTRL> + C.")

	// Starts the router
	handler := a.router.Start()

	if a.cors != nil {
		// Starts the server with CORS support
		return http.ListenAndServe(
			fmt.Sprintf(":%s", port), a.cors.Handler(handler))
	}

	// Starts the server
	return http.ListenAndServe(
		fmt.Sprintf(":%s", port), handler)
}
