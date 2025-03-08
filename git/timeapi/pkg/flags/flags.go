package flags

import (
	"time"

	"github.com/spf13/pflag"
)

type Options struct {
	Addr        string
	ReadTimeout time.Duration
}

func NewOptions() *Options {
	return &Options{}
}

func Init(opts *Options) {
	pflag.StringVar(&opts.Addr, "addr", ":8080", "Address to listen on")
	pflag.DurationVar(&opts.ReadTimeout, "read-timeout", 5*time.Second, "Read timeout")
	pflag.Parse()
}
