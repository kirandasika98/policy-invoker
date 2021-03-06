package main

import (
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/kirandasika98/policy-invoker/pkg/handler"
)

var (
	addr               *string
	policiesPath       *string
	sentinelBinaryPath *string
)

func main() {
	addr = flag.String("addr", "localhost:9000", "the host:port that the application should run")
	policiesPath = flag.String("policies-path", ".", "directory where all the policies exist")
	sentinelBinaryPath = flag.String("sentinel-binary-path", "/usr/local/bin/sentinel", "absolute path to the sentinel binary")
	flag.Parse()

	defer glog.Flush()
	r := mux.NewRouter()
	r.Handle("/invoke/{policy}", handler.NewPolicyInvokeHandler(*sentinelBinaryPath)).Methods("POST")
	glog.Infof("starting policy-invoke server on %s", *addr)
	glog.Fatal(http.ListenAndServe(*addr, r))
}
