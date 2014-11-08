package codec

import (
	"bytes"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type YAML struct{}

func (y YAML) NewDecoder(r io.Reader) Decoder {
	return yaml_codec{r: r}
}

func (y YAML) NewEncoder(w io.Writer) Encoder {
	return yaml_codec{w: w}
}

type yaml_codec struct {
	r io.Reader
	w io.Writer
}

func (y yaml_codec) Decode(v interface{}) error {

	data, err := ioutil.ReadAll(y.r)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, v)
}

func (y yaml_codec) Encode(v interface{}) error {
	data, err := yaml.Marshal(v)

	if err != nil {
		return err
	}
	_, err = bytes.NewBuffer(data).WriteTo(y.w)
	return err
}
