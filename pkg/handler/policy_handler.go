package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/kirandasika98/policy-invoker/pkg/invoker"
	"github.com/kirandasika98/policy-invoker/pkg/policy"
)

// PolicyInvokeHandler handles calls to a Invoker
type PolicyInvokeHandler struct {
}

type invokePolicyData struct {
	Cfg string `json:"cfg"`
}

type invokePolicyErr struct {
	Error string `json:"error"`
}

// NewPolicyInvokeHandler is a function that returns a new handler for invoking policies
func NewPolicyInvokeHandler() *PolicyInvokeHandler {
	return &PolicyInvokeHandler{}
}

// ServeHTTP is a function that should be implemented as part of the http.handler
func (h *PolicyInvokeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	glog.Infof("policy: %s", vars["policy"])
	var data invokePolicyData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		glog.Errorf("error while decoding json: %v", err)
		handleError(w, err)
		return
	}
	glog.Infof("cfg: %s", data.Cfg)
	p, err := policy.NewFromPolicyName(".", vars["policy"]+".sentinel", data.Cfg)
	if err != nil {
		glog.Errorf("error while building new policy: %v", err)
		handleError(w, err)
		return
	}
	invoker, err := invoker.New(p)
	if err != nil {
		glog.Errorf("error while build invoker: %v", err)
		handleError(w, err)
		return
	}
	exitCode, err := invoker.Invoke()
	if err != nil {
		glog.Errorf("error while invoking %T policy: %v", p, err)
		handleError(w, err)
		return
	}
	w.Write([]byte(strconv.Itoa(exitCode)))
}

func handleError(w http.ResponseWriter, err error) {
	res := &invokePolicyErr{
		Error: fmt.Sprintf("%v", err),
	}
	b, _ := json.Marshal(res)
	http.Error(w, string(b), http.StatusBadRequest)
}
