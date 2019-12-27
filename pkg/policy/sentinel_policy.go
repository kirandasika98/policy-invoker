package policy

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

//SentinelPolicy is a struct
type SentinelPolicy struct {
	Cfg    string `json:"params"`
	Policy string `json:"policy"`
}

// NewSentinelPolicy is a function that builds a SentinelPolicy
func NewSentinelPolicy(cfgBuf, policyBuf []byte) (*SentinelPolicy, error) {
	return &SentinelPolicy{
		Cfg:    string(cfgBuf),
		Policy: string(policyBuf),
	}, nil
}

func NewFromPolicyName(policyPath, policyName, cfg string) (*SentinelPolicy, error) {
	if len(policyName) == 0 {
		return nil, errors.New("policy: policyName cannot be empty")
	}
	if len(cfg) == 0 {
		return nil, errors.New("policy: config for the policy cannot be nil")
	}
	fp := filepath.Join(policyPath, policyName)
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	p := &SentinelPolicy{
		Policy: string(b),
		Cfg:    cfg,
	}
	return p, nil
}
