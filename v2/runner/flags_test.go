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
		inputCmdLine:   "-grpc-ip=1.2.3.56 --log-level=4 --grpc-ping-timeout=5s --grpc-ping-max-missed=3 --plugin-config={}",
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
	{
		inputCmdLine:   "--grpc-ip=1.2.3.4.5",
		shouldBeParsed: true,
		shouldBeValid:  false,
	},
	{
		inputCmdLine:   "--pprof-port=5678",
		shouldBeParsed: true,
		shouldBeValid:  false,
	},
	{
		inputCmdLine:   "--enable-pprof=1 --pprof-port=5678",
		shouldBeParsed: true,
		shouldBeValid:  true,
	},
	{
		inputCmdLine:   "--stats-port=5678",
		shouldBeParsed: true,
		shouldBeValid:  false,
	},
	{
		inputCmdLine:   "--enable-stats=1 --stats-port=5678",
		shouldBeParsed: true,
		shouldBeValid:  true,
	},
	{
		inputCmdLine:   "--debug-collect-counts=11",
		shouldBeParsed: true,
		shouldBeValid:  false,
	},
	{
		inputCmdLine:   "--debug-mode=1 --debug-collect-counts=11",
		shouldBeParsed: true,
		shouldBeValid:  true,
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
