package main

import (
	"fmt"
	"log"

	"github.com/cppforlife/go-patch/patch"
)

func main() {
	man := map[interface{}]interface{}{
		"releases": []interface{}{
			map[string]interface{}{
				"name": "asdf",
			},
		},
	}

	varr, err := patch.ReplaceOp{
		Path:  patch.MustNewPointerFromString("/releases/url?"),
		Value: "asdfasdf",
	}.Apply(man)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("%#+v", varr)
}
