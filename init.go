package requests

import (
	"github.com/Zcentury/gologger"
	"github.com/Zcentury/requests/config"
)

func init() {
	gologger.LoggerOptions.SetFormatter(config.NewCLI(true))
	gologger.LoggerOptions.SetTimestamp(true)
}
