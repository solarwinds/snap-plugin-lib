// +build small

package proxy

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type jsonScenario struct {
	inputJSON      string
	expectedResult map[string]string
	description    string
}

var jsonScenarios = []jsonScenario{
	{ // 0
		inputJSON: `{
			"address": {
				"ip": "192.153.25.123",
				"port": 34245
 		   	},
    		"credentials": "password",
    		"user": "admin" 
		}`,
		expectedResult: map[string]string{
			"address.ip":   "192.153.25.123",
			"address.port": "34245",
			"credentials":  "password",
			"user":         "admin",
		},
		description: "Basic configuration",
	},
	{ // 1
		inputJSON: `{
			"address": {
       			"ip": "192.153.25.123",
				"port": 34245
    		},
  		    "user": "admin",
    		"rights": ["admin", "logger", "runner", "reader", "writer"]
		}`,
		expectedResult: map[string]string{
			"address.ip":   "192.153.25.123",
			"address.port": "34245",
			"user":         "admin",
			"rights":       "admin,logger,runner,reader,writer",
		},
		description: "Configuration with leaf array of a simple type",
	},
	{ // 2
		inputJSON: `{
 		   "addresses": [ 
				{
					"protocol": "tcp",
            		"ip": "192.153.25.123",
            		"port": 34245    
        		},
       			{
            		"protocol": "udp",
            		"ip": "122.178.11.1",
            		"port": 2316
        		}
    		],
    		"user": "admin"
		}`,
		expectedResult: map[string]string{
			// addresses doesn't contain simple type and is ignored
			"user": "admin",
		},
		description: "Configuration with array containing not simple elements",
	},
	{ // 3
		inputJSON: `{
			"address": {
				"ip": "192.153.25.123",
				"port": 34245
 		   	},
    		"credentials": "password",
    		"user": "admin" 
		}`,
		expectedResult: map[string]string{
			"address.ip":   "192.153.25.123",
			"address.port": "34245",
			"credentials":  "password",
			"user":         "admin",
		},
		description: "Basic configuration",
	},
}

func TestJSONToMap(t *testing.T) {
	Convey("Validate that json configuration can be properly flatten to map", t, func() {
		for i, testCase := range jsonScenarios {
			Convey(fmt.Sprintf("Scenario %d - %s", i+1, testCase.description), func() {
				// Act
				result, err := JSONToMap(testCase.inputJSON)

				// Assert
				So(err, ShouldBeNil)
				So(result, ShouldResemble, testCase.expectedResult)
			})
		}
	})
}
