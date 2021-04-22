package cli

import (
	"github.com/AndreZiviani/MyIP/internal/serve"
	flags "github.com/jessevdk/go-flags"
)

var Parser = flags.NewParser(nil, flags.Default)

func Run() {
	serve.Init(Parser)
	Parser.Parse()
}
