// +build small

package runner

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

///////////////////////////////////////////////////////////////////////////////

type parseScenario struct {
	inputCmdLine   string
	shouldBeParsed bool
	shouldBeValid  bool
}

var parseScenarios = []parseScenario{
	{ // 0
		inputCmdLine:   "-grpc-ip=1.2.3.4 --grpc-port=456 --log-level=warning",
		shouldBeParsed: true,
		shouldBeValid:  true,
	},
	{ // 1
		inputCmdLine:   "-grpc-ip=1.2.3.56 --log-level=4 --grpc-ping-timeout=5s --grpc-ping-max-missed=3",
		shouldBeParsed: true,
		shouldBeValid:  true,
	},
	{ // 2
		inputCmdLine:   "",
		shouldBeParsed: true,
		shouldBeValid:  true,
	},
	{ // 3
		inputCmdLine:   "--grpc-port=abc",
		shouldBeParsed: false,
		shouldBeValid:  false,
	},
	{ // 4
		inputCmdLine:   "--debug-level=invalid",
		shouldBeParsed: false,
		shouldBeValid:  false,
	},
	{ // 5
		inputCmdLine:   "--debug-level=8",
		shouldBeParsed: false,
		shouldBeValid:  false,
	},
}

func TestParseCmdLineOptions(t *testing.T) {
	Convey("Validate that options can be parsed", t, func() {
		for i, testCase := range parseScenarios {
			Convey(fmt.Sprintf("Scenario %d [%s]", i, testCase.inputCmdLine), func() {
				// Arrange
				inputCmd := []string{}
				if len(testCase.inputCmdLine) > 0 {
					inputCmd = strings.Split(testCase.inputCmdLine, " ")
				}

				// Act
				opt, err := ParseCmdLineOptions("plugin", inputCmd)

				fmt.Printf("opt=%#v\n", opt)
				fmt.Printf("err=%#v\n", err)

				// Assert
				if testCase.shouldBeParsed {
					So(err, ShouldBeNil)

					// Act
					validErr := ValidateOptions(opt)

					// Assert
					if testCase.shouldBeValid {
						So(validErr, ShouldBeNil)
					} else {
						So(validErr, ShouldBeError)
					}
				} else {
					So(err, ShouldBeError)
				}
			})
		}
	})
}

///////////////////////////////////////////////////////////////////////////////
