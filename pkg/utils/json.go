package utils

import (
	"encoding/json"
	"errors"
	"log"
)

var (
	ErrMarshalingArgs      = errors.New("error marshaling args")
	ErrUnmarshalingRequest = errors.New("error unmarshaling request")
)

func ArgsToStruct(args map[string]any, response any) error {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		log.Fatal(ErrMarshalingArgs)
	}

	if err := json.Unmarshal(argsJSON, &response); err != nil {
		log.Fatal(ErrUnmarshalingRequest)
	}

	return nil
}
