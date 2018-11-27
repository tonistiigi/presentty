package localcommand

import (
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/yudai/gotty/server"
)

type Options struct {
	CloseSignal  int `hcl:"close_signal" flagName:"close-signal" flagSName:"" flagDescribe:"Signal sent to the command process when gotty close it (default: SIGHUP)" default:"1"`
	CloseTimeout int `hcl:"close_timeout" flagName:"close-timeout" flagSName:"" flagDescribe:"Time in seconds to force kill process after client is disconnected (default: -1)" default:"-1"`
}

type Factory struct {
	m       map[string][]string
	options *Options
	opts    []Option
}

func NewFactory(m map[string][]string, options *Options) (*Factory, error) {
	opts := []Option{WithCloseSignal(syscall.Signal(options.CloseSignal))}
	if options.CloseTimeout >= 0 {
		opts = append(opts, WithCloseTimeout(time.Duration(options.CloseTimeout)*time.Second))
	}

	return &Factory{
		m: m,
		// options: options,
		opts: opts,
	}, nil
}

func (factory *Factory) Name() string {
	return "docker command"
}

func (factory *Factory) New(params map[string][]string) (server.Slave, error) {
	// argv := make([]string, len(factory.argv))
	// copy(argv, factory.argv)
	// if params["arg"] != nil && len(params["arg"]) > 0 {
	// 	argv = append(argv, params["arg"]...)
	// }

	id, ok := params["demo"]
	if !ok {
		return nil, errors.Errorf("missing demo param")
	}

	args, ok := factory.m[id[0]]
	if !ok {
		return nil, errors.Errorf("invalid demo id %s", id)
	}

	return New(args[0], args[1])
}
