package goshopee

import "encoding/json"

func ToMapData(in interface{}) (map[string]interface{}, error) {
	byts, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(byts, &result); err != nil {
		return nil, err
	}
	return result, nil
}
