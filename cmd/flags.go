package cmd

import (
	"flag"

	"github.com/alexdebril/figaro/log"
)

type Flags struct {
	logJson bool
	Debug   bool
}

func (f *Flags) GetLogFormat() log.Format {
	if f.logJson {
		return log.JSON
	}
	return log.TEXT
}

func GetCoreFlags() (*Flags, *flag.FlagSet) {
	flg := &Flags{
		logJson: false,
		Debug:   false,
	}
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.BoolVar(&flg.Debug, "debug", flg.Debug, "Run the server in debug mode")
	fs.BoolVar(&flg.logJson, "log-json", flg.logJson, "Log in JSON format")
	return flg, fs
}
