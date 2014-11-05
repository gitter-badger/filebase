package codec

import (
	"bytes"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type YAML struct{}

func (y YAML) NewDecoder(r io.Reader) decoder {
	return yaml_decoder{r}
}

func (y YAML) NewEncoder(r io.Writer) encoder {
	return yaml_encoder{r}
}

type yaml_decoder struct {
	r io.Reader
}

func (y yaml_decoder) Decode(v interface{}) error {

	data, err := ioutil.ReadAll(y.r)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, v)
}

type yaml_encoder struct {
	w io.Writer
}

func (y yaml_encoder) Encode(v interface{}) error {
	data, err := yaml.Marshal(v)

	if err != nil {
		return err
	}
	_, err = bytes.NewBuffer(data).WriteTo(y.w)
	return err
}
