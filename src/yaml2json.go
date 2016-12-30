package yaml2json

import (
	"encoding/json"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"strconv"
)

// YAMLToJSON converts YAML unmarshaled data into json compatible data
func YAMLToJSON(data interface{}) (json.RawMessage, error) {
	jm, err := transformData(data)
	if err != nil {
		return nil, err
	}
	b, err := json.MarshalIndent(jm, "", "  ")
	return json.RawMessage(b), err
}

func BytesToYAMLDoc(data []byte) (interface{}, error) {
	var document map[interface{}]interface{}
	if err := yaml.Unmarshal(data, &document); err != nil {
		return nil, err
	}

	return document, nil
}

func transformData(in interface{}) (out interface{}, err error) {
	switch in.(type) {
	case map[interface{}]interface{}:
		o := make(map[string]interface{})
		for k, v := range in.(map[interface{}]interface{}) {
			sk := ""
			switch k.(type) {
			case string:
				sk = k.(string)
			case int:
				sk = strconv.Itoa(k.(int))
			default:
				return nil, fmt.Errorf("types don't match: expect map key string or int get: %T", k)
			}
			v, err = transformData(v)
			if err != nil {
				return nil, err
			}
			o[sk] = v
		}
		return o, nil
	case []interface{}:
		in1 := in.([]interface{})
		len1 := len(in1)
		o := make([]interface{}, len1)
		for i := 0; i < len1; i++ {
			o[i], err = transformData(in1[i])
			if err != nil {
				return nil, err
			}
		}
		return o, nil
	}
	return in, nil
}
