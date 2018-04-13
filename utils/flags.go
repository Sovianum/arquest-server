package utils

import "flag"

func NewFlags(defaultConfig string) *Flags {
	return &Flags{
		defaultConfig: defaultConfig,
	}
}

type Flags struct {
	Config        string
	defaultConfig string
}

func (f *Flags) Parse() {
	flag.StringVar(&f.Config, "c", f.defaultConfig, "path to config file")
	flag.Parse()
}
