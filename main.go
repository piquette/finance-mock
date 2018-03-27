package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/piquette/finance-mock/fixture"
	"github.com/piquette/finance-mock/server"
	yaml "gopkg.in/yaml.v2"
)

const (
	defaultPort         = 12111
	defaultFixturesPath = "./fixture/resources.json"
	defaultSpecPath     = "./fixture/spec.yml"
)

// verbose tracks whether the program is operating in verbose mode
var verbose bool

// This is set to the actual version by GoReleaser (using `-ldflags "-X ..."`)
// as it's run. Versions built from source will always show master.
var version = "master"

func main() {
	var showVersion bool
	var port int
	var fixturesPath string
	var specPath string
	var unix string

	flag.IntVar(&port, "port", defaultPort, "Port to listen on")
	flag.StringVar(&fixturesPath, "fixtures", defaultFixturesPath, "Path to fixtures to use instead of bundled version")
	flag.StringVar(&specPath, "spec", defaultSpecPath, "Path to spec to use instead of bundled version")
	flag.StringVar(&unix, "unix", "", "Unix socket to listen on")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose mode")
	flag.BoolVar(&showVersion, "version", false, "Show version and exit")
	flag.Parse()

	if showVersion || len(flag.Args()) == 1 && flag.Arg(0) == "version" {
		fmt.Printf("%s\n", version)
		return
	}

	if unix != "" && port != 0 {
		flag.Usage()
		abort("Specify only one of -port or -unix\n")
	}

	// Get spec.
	spec, err := getSpec(specPath)
	if err != nil {
		abort(err.Error())
	}

	// Get fixtures.
	fixtures, err := getFixtures(fixturesPath)
	if err != nil {
		abort(err.Error())
	}

	// Stub server.
	stub := server.StubServer{Fixtures: fixtures, Spec: spec, Verbose: verbose}
	server.Version = version

	// Initialize server router.
	err = stub.InitRouter()
	if err != nil {
		abort(fmt.Sprintf("Error initializing router: %v\n", err))
	}

	// Set handler.
	http.HandleFunc("/", stub.HandleRequest)
	s := http.Server{}

	// Init listener.
	listener, err := getListener(port, unix)
	if err != nil {
		abort(err.Error())
	}

	// Serve.
	s.Serve(listener)
}

func abort(message string) {
	fmt.Fprintf(os.Stderr, message)
	os.Exit(1)
}

func getFixtures(fixturesPath string) (*fixture.Fixtures, error) {
	var data []byte
	var err error

	data, err = ioutil.ReadFile(fixturesPath)

	if err != nil {
		return nil, fmt.Errorf("error loading fixtures: %v\n", err)
	}

	var fixtures fixture.Fixtures
	err = json.Unmarshal(data, &fixtures)
	if err != nil {
		return nil, fmt.Errorf("error decoding fixtures: %v\n", err)
	}
	return &fixtures, nil
}

func getListener(port int, unix string) (net.Listener, error) {
	var err error
	var listener net.Listener

	if unix != "" {
		listener, err = net.Listen("unix", unix)
		fmt.Printf("Listening on unix socket %v\n", unix)
	} else {
		if port == 0 {
			port = defaultPort
		}
		listener, err = net.Listen("tcp", ":"+strconv.Itoa(port))
		fmt.Printf("Listening on port %v\n", port)
	}

	if err != nil {
		return nil, fmt.Errorf("error listening on socket: %v\n", err)
	}

	return listener, nil
}

func getSpec(specPath string) (*fixture.Spec, error) {
	var data []byte
	var err error

	data, err = ioutil.ReadFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("Error loading spec: %v\n", err)
	}

	var spec fixture.Spec

	err = yaml.Unmarshal(data, &spec)
	if err != nil {
		return nil, fmt.Errorf("Error decoding spec: %v\n", err)
	}

	return &spec, nil
}
