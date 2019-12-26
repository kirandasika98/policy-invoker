package invoker

import (
	"log"
	"testing"

	"github.com/kirandasika98/policy-invoker/pkg/policy"
)

func TestBuildConfig(t *testing.T) {
	cfgJSON := `{
		"policy": "",
		"params": "{\"DOB\":\"1998-04-30T00:00:00Z\",\"consents\":[\"data_exchange\",\"marketing\",\"analytics\"],\"target_consent\":\"data_exchange\"}"
	}`
	p, err := policy.NewFromReader([]byte(cfgJSON))
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}
	if len(p.Cfg) == 0 {
		t.Fail()
	}
}

type mockInvoker struct{}

func (mi *mockInvoker) Invoke() (int, error) {
	return 0, nil
}

func TestInvoke(t *testing.T) {

}
