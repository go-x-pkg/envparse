package envparse_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/go-x-pkg/envparse"
)

func TestParse(t *testing.T) {
	tests := []struct {
		raw string

		wEnv envparse.Env
		werr error
	}{
		{
			`foo=bar  myDog=Rex\ The\ Dog   mycat="dqwdqwd dwqdqdwq " bar=d12`,

			envparse.Env{
				"foo":   []string{`bar`},
				"myDog": []string{`Rex\ The\ Dog`},
				"mycat": []string{`dqwdqwd dwqdqdwq `},
				"bar":   []string{`d12`},
			},
			nil,
		},
		{
			`LD_LIBRARY_PATH=/usr/lib:/usr/nvidia/cuda/lib64 PATH=/usr/bin`,

			envparse.Env{
				"LD_LIBRARY_PATH": []string{`/usr/lib:/usr/nvidia/cuda/lib64`},
				"PATH":            []string{`/usr/bin`},
			},
			nil,
		},
		{
			`LD_LIBRARY_PATH /usr/lib:/usr/nvidia/cuda/lib64`,

			envparse.Env{
				"LD_LIBRARY_PATH": []string{`/usr/lib:/usr/nvidia/cuda/lib64`},
			},
			nil,
		},

		{
			`LD_LIBRARY_PATH`,

			nil,
			envparse.ErrMustHaveTwoArguments,
		},

		{
			`LD_LIBRARY_PATH=foo bar buz else ever`,

			nil,
			envparse.ErrMissingEqualsSign,
		},
	}

	for i, tt := range tests {
		func() {
			env, err := envparse.Parse(tt.raw)

			if !errors.Is(err, tt.werr) {
				t.Errorf("#%d: err = %v, want %v", i, err, tt.werr)
			}

			if eq := reflect.DeepEqual(env, tt.wEnv); !eq {
				t.Errorf("#%d: got = %v, want %v", i, env, tt.wEnv)
			}
		}()
	}
}
