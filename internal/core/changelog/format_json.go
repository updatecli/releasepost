package changelog

import (
	"encoding/json"
	"fmt"
)

func toJsonFile(s interface{}, filename string) error {
	var data []byte
	var err error

	data, err = json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling json: %v", err)
	}

	err = dataToFile(data, filename)
	if err != nil {
		return fmt.Errorf("creating file %s: %v", filename, err)
	}

	return nil
}
