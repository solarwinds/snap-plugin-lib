package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Covert json to string map (key are json path to access element, values are element), ie
// input json = `{ "credentials": { "account_types": [ "admin", "management", "service", "debug" ], "name": "admin" }, "server": { "ip": "192.168.56.101", "port": 1234 }}`,
// output     = map[string]string{"credentials.name":"admin", "credentials.account_types":"admin,management,service,debug", "server.ip":"192.168.56.101", "server.port":"1234"}
func JSONToFlatMap(rawConfig string) (map[string]string, error) {
	retMap := map[string]string{}
	m := map[string]interface{}{}

	err := json.Unmarshal([]byte(rawConfig), &m)
	if err != nil {
		return nil, fmt.Errorf("can't create configuration map based on provided json: %v", err)
	}

	parseMap(m, retMap)

	return retMap, nil
}

func parseMap(aMap map[string]interface{}, retMap map[string]string, path ...string) {
	for key, val := range aMap {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			newPath := append(path, key)
			parseMap(val.(map[string]interface{}), retMap, newPath...)
		case []interface{}:
			rv := parseArray(val.([]interface{}))
			if len(rv) > 0 {
				newPath := append(path, key)
				retMap[strings.Join(newPath, ".")] = strings.Join(rv, ",")
			}
		default:
			newPath := append(path, key)
			retMap[strings.Join(newPath, ".")] = fmt.Sprintf("%v", concreteVal)
		}
	}
}

func parseArray(anArray []interface{}) []string {
	res := []string{}

	for _, val := range anArray {
		switch val.(type) {
		case map[string]interface{}:
			// ignored
		case []interface{}:
			// ignored
		default:
			res = append(res, val.(string))
		}
	}

	return res
}