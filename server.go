package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

const (
	darkskyAPIURI = "https://api.darksky.net/forecast"
	geocodeAPIURI = "https://maps.googleapis.com/maps/api/geocode/json"
)

const serverHelp = `Run a static UI server for a registry.`

func (cmd *serverCommand) Name() string      { return "server" }
func (cmd *serverCommand) Args() string      { return "[OPTIONS]" }
func (cmd *serverCommand) ShortHelp() string { return serverHelp }
func (cmd *serverCommand) LongHelp() string  { return serverHelp }
func (cmd *serverCommand) Hidden() bool      { return false }

func (cmd *serverCommand) Register(fs *flag.FlagSet) {
	fs.StringVar(&cmd.darkskyAPIKey, "darksky-apikey", "", "Key for darksky.net API")
	fs.StringVar(&cmd.geocodeAPIKey, "geocode-apikey", "", "Key for Google Maps Geocode API")

	fs.StringVar(&cmd.cert, "cert", "", "path to ssl cert")
	fs.StringVar(&cmd.key, "key", "", "path to ssl key")
	fs.StringVar(&cmd.port, "port", "1234", "port for server to run on")
}

type serverCommand struct {
	darkskyAPIKey string
	geocodeAPIKey string

	cert string
	key  string
	port string
}

func (cmd *serverCommand) Run(ctx context.Context, args []string) error {
	// On ^C, or SIGTERM handle exit.
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, syscall.SIGTERM)
	var cancel context.CancelFunc
	_, cancel = context.WithCancel(ctx)
	go func() {
		for sig := range signals {
			cancel()
			logrus.Infof("Received %s, exiting.", sig.String())
			os.Exit(0)
		}
	}()

	if len(cmd.darkskyAPIKey) < 1 {
		return errors.New("please pass a darksky.net API Key")
	}

	if len(cmd.geocodeAPIKey) < 1 {
		logrus.Fatalf("Please pass a Google Maps Geocode API Key")
	}

	// Create mux server.
	mux := http.NewServeMux()

	mux.HandleFunc("/forecast", cmd.forecastHandler) // forecast handler
	mux.HandleFunc("/geocode", cmd.geocodeHandler)   // geocode handler
	mux.HandleFunc("/", failHandler)                 // everything else fail handler

	// Set up the server.
	server := &http.Server{
		Addr:    ":" + cmd.port,
		Handler: mux,
	}
	logrus.Infof("Starting server on port %q", cmd.port)
	if len(cmd.cert) > 0 && len(cmd.key) > 0 {
		return server.ListenAndServeTLS(cmd.cert, cmd.key)
	}
	return server.ListenAndServe()
}
