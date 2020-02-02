package main

import (
	"log"

	"github.com/crabtree/defeway-toolbox/pkg/cmdtoolbox"
)

func main() {
	params, err := NewParams()
	cmdtoolbox.DieOnError(err)

	log.Println(params.Dump())
}
