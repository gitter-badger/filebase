package codec

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type YAML struct {
}

type yaml_decoder struct {
	r io.Reader
}

func (d yaml_decoder) Decode(v interface{}) error {
	data, err := ioutil.ReadAll(d.r)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, v)
}

func (j YAML) NewDecoder(r io.Reader) decoder {
	return json.NewDecoder(r)
}

type yaml_encoder struct {
	w io.Writer
}

func (e yaml_encoder) Encode(v interface{}) error {
	data, err := yaml.Marshal(v)

	if err != nil {
		return err
	}

	_, err = bytes.NewBuffer(data).WriteTo(e.w)
	return err
}
