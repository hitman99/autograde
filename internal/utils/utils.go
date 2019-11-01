package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func UnmarshalBody(b io.ReadCloser, target interface{}) error {
	body, err := ioutil.ReadAll(b)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, target)
	if err != nil {
		return err
	}
	return nil
}
