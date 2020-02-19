package main

import (
	"log"

	"github.com/crabtree/defeway-toolbox/internal/scanner"
	"github.com/crabtree/defeway-toolbox/pkg/cmdtoolbox"
)

func main() {
	params, err := NewParams()
	cmdtoolbox.DieOnError(err)

	log.Println(params.Dump())

	command := scanner.NewCommand(
		paramsToCommandParams(params))

	err = command.Run()
	cmdtoolbox.DieOnError(err)
}

func paramsToCommandParams(params *params) scanner.ScannerParams {
	return scanner.ScannerParams{
		Concurrent: params.Concurrent,
		LogDir:     params.LogDir,
		NetAddr:    params.NetAddr,
		NetMask:    params.NetMask,
		Password:   params.Password,
		Ports:      params.Ports,
		Username:   params.Username,
	}
}
