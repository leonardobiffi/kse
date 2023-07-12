package encode

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"reflect"

	"github.com/leonardobiffi/kse/package/entities"
	"gopkg.in/yaml.v2"
)

type service struct{}

var _ Service = (*service)(nil)

func (s *service) Execute(content []byte) (output []byte, err error) {
	dec := yaml.NewDecoder(bytes.NewReader(content))

	var outputFiles []yaml.MapSlice
	for {
		var m yaml.MapSlice
		if err = dec.Decode(&m); err != nil {
			// Break when there are no more documents to decode
			if err != io.EOF {
				return nil, err
			}
			break
		}

		var secretData yaml.MapItem
		for _, item := range m {
			if item.Key == "data" {
				secretData = item
			}
		}

		data, ok := cast(secretData.Value, false)
		if !ok || len(data) == 0 {
			return nil, fmt.Errorf("invalid secret data")
		}

		secretData.Value = encode(data)
		for i, item := range m {
			if item.Key == "data" {
				m[i] = secretData
			}
		}

		outputFiles = append(outputFiles, m)
	}

	for i, file := range outputFiles {
		out, err := yaml.Marshal(file)
		if err != nil {
			return nil, err
		}

		output = append(output, out...)
		// add separator between documents, and skip last separator
		if i < len(outputFiles)-1 {
			output = append(output, []byte("---\n")...)
		}
	}

	return output, nil
}

func cast(data interface{}, isJSON bool) (map[string]interface{}, bool) {
	if isJSON {
		d, ok := data.(map[string]interface{})
		return d, ok
	}

	val := reflect.ValueOf(data)
	d := make(map[string]interface{}, val.Len())
	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)
		// print key of map item
		key := fmt.Sprintf("%v", item.Field(0))
		// print value of map item
		value := fmt.Sprintf("%v", item.Field(1))

		d[key] = value
	}

	return d, true
}

func encode(data map[string]interface{}) map[string]string {
	length := len(data)
	secrets := make(chan entities.Data, length)
	encoded := make(map[string]string, length)
	for key, decoded := range data {
		go encodeSecret(key, decoded.(string), secrets)
	}
	for i := 0; i < length; i++ {
		secret := <-secrets
		encoded[secret.Key] = secret.Value
	}
	return encoded
}

func encodeSecret(key, secret string, secrets chan entities.Data) {
	var value string

	// avoid wrong encoded secrets
	if encoded := base64.StdEncoding.EncodeToString([]byte(secret)); encoded != secret {
		value = encoded
	} else {
		value = secret
	}
	secrets <- entities.Data{Key: key, Value: value}
}

func New() Service {
	svc := &service{}
	return svc
}
