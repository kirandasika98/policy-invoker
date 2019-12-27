package invoker

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/google/uuid"
	"github.com/kirandasika98/policy-invoker/pkg/policy"
)

// Invoker is a interface that invokes a policy
type Invoker interface {
	Invoke() (int, error)
}

const (
	sentinel = "/usr/local/bin/sentinel"
	apply    = "apply"
	config   = "-config"
	trace    = "-trace"
)

// sentinelPolicyInvoker invokes a sentinel policy
type sentinelPolicyInvoker struct {
	cfg    *os.File
	policy *os.File
	cmd    *exec.Cmd
}

// New builds a new Invoker
func New(p *policy.SentinelPolicy) (Invoker, error) {
	if p == nil {
		return nil, errors.New("invoker: policy cannot be nil")
	}
	invoker := &sentinelPolicyInvoker{}
	if len(p.Cfg) == 0 || len(p.Policy) == 0 {
		return nil, errors.New("invoker: either configuration or the policy is empty")
	}
	// sets the cfg variable
	invoker.buildConfig([]byte(p.Cfg))
	// sets the policy variable
	invoker.buildPolicyFile([]byte(p.Policy))
	// sets the cmd variable
	invoker.buildSentinelCommand()
	return invoker, nil
}

// Invoke invokes a policy and return the exitCode and an error if any
func (pi *sentinelPolicyInvoker) Invoke() (int, error) {
	//log.Println("Invoking run")
	//log.Printf("cfg path: %s", pi.cfg.Name())
	//log.Printf("policy: %s", pi.policy.Name())
	//log.Printf("%s", pi.cmd.String())
	//bb, _ := ioutil.ReadFile(pi.cfg.Name())
	//log.Printf("data from cfg: %s", bb)
	defer pi.destroyRun()
	bb, err := pi.cmd.CombinedOutput()
	if err != nil {
		log.Printf("failed with: %s", bb)
		return -1, err
	}
	return 0, nil
}

func (pi *sentinelPolicyInvoker) destroyRun() error {
	defer os.Remove(pi.cfg.Name())
	defer os.Remove(pi.policy.Name())
	if err := pi.cfg.Close(); err != nil {
		return nil
	}
	if err := pi.policy.Close(); err != nil {
		return nil
	}
	return nil
}

func (pi *sentinelPolicyInvoker) buildConfig(buf []byte) error {
	f, err := ioutil.TempFile(os.TempDir(), "policy-invoker-")
	if err != nil {
		return err
	}
	pi.cfg = f
	if _, err := f.Write([]byte(buf)); err != nil {
		return err
	}
	return nil
}

func (pi *sentinelPolicyInvoker) buildPolicyFile(buf []byte) error {
	f, err := ioutil.TempFile(os.TempDir(), "policy-invoker-")
	if err != nil {
		return err
	}
	pi.policy = f
	if _, err := f.Write(buf); err != nil {
		return err
	}
	return nil
}

func (pi *sentinelPolicyInvoker) buildSentinelCommand() {
	pi.cmd = exec.Command(sentinel, apply, config, pi.cfg.Name(), pi.policy.Name())
}

func getRandomFileName(ext string) string {
	return fmt.Sprintf("%s%s", uuid.New(), ext)
}
