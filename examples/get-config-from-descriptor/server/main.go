package server

import (
	"fmt"
	"html"
	"net/http"

	"github.com/typical-go/typical-go/pkg/typcore"
)

const (
	// ConfigName is lookup key in ConfigManager
	ConfigName = "SERVER"
)

// Config of app
type Config struct {
	Address string `default:":8080" required:"true"`
}

// Main function to run the server
func Main(d *typcore.Descriptor) (err error) {
	var spec interface{}
	if spec, err = d.RetrieveConfig(ConfigName); err != nil {
		return
	}

	// type assertion to Config type
	cfg := spec.(*Config)

	fmt.Printf("Get Config From Descriptor -- Serve http at %s\n", cfg.Address)
	return http.ListenAndServe(cfg.Address, &handler{})
}

// Handler responsible to handle server request
type handler struct{}

func (*handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}