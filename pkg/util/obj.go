package util

import "gopkg.in/yaml.v2"

// InitObjViaYaml marshal "config" to yaml data, then unMarshal data to "obj"
func InitObjViaYaml(obj interface{}, config interface{}) error {
	if obj == nil || config == nil {
		return nil
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, obj)
}
