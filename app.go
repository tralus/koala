package koala

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/tralus/koala/config"
	"github.com/tralus/koala/knife"
)

// Config holds the app config
var Config config.Config

// ServerPort holds the port of tht server
var ServerPort string

func init() {
	var err error

	// Sets the server port of the application
	flag.StringVar(&ServerPort, "koala_server_port", ":9003", "server port")

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
	router  *knife.Router
	modules []Module
}

// NewApplication creates an instance of App
func NewApplication(r *knife.Router) App {
	return App{router: r}
}

// SetRouter sets the application router
func (a *App) SetRouter(r *knife.Router) {
	a.router = r
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

	fmt.Println(sep())
	fmt.Printf("Debug: %t\n", Config.Debug)
	fmt.Printf("Port: %s\n", ServerPort)

	fmt.Println(sep())
	fmt.Printf("On http://localhost%s\n", ServerPort)
	fmt.Println("To shut down, press <CTRL> + C.")

	// Starts the server
	return http.ListenAndServe(
		ServerPort,
		http.Handler(a.router.Start()),
	)
}
