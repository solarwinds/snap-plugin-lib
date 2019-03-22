// +build small

package proxy

import (
	"encoding/json"
	"testing"
)
import . "github.com/smartystreets/goconvey/convey"

type basicConfig struct {
	Address struct {
		Ip   string
		Port int
	}
	Rights []string
	User   string
}

func TestContextAPI_Config(t *testing.T) {
	// Arrange
	jsonConfig := `{
    	"address": {
        	"ip": "192.153.25.123",
        	"port": 34245
    	},
    	"rights": ["admin", "logger", "runner", "reader", "writer"], 
    	"user": "admin"
	}`

	ctx, cErr := NewPluginContext(jsonConfig, []string{})

	Convey("Validate Context API for handling configuration", t, func() {

		So(cErr, ShouldBeNil)

		Convey("Validate Context::Config", func() {

			Convey("User can read correct configuration field", func() {

				// Act
				val, ok := ctx.Config("address.ip")

				// Assert
				So(ok, ShouldBeTrue)
				So(val, ShouldEqual, "192.153.25.123")

			})

			Convey("User can read incorrect configuration field", func() {

				// Act
				_, ok := ctx.Config("address.protocol")

				// Assert
				So(ok, ShouldBeFalse)

			})

			Convey("User can read correct configuration field (2)", func() {

				// Act
				val, ok := ctx.Config("address.port")

				// Assert
				So(ok, ShouldBeTrue)
				So(val, ShouldEqual, "34245")

			})
		})

		Convey("Validate Context::ConfigKeys", func() {

			keyList := []string{}

			Convey("User can read allowed config fields", func() {

				// Act
				keyList = ctx.ConfigKeys()

				// Assert
				So(len(keyList), ShouldEqual, 4)
				So(keyList, ShouldContain, "address.ip")
				So(keyList, ShouldContain, "address.port")
				So(keyList, ShouldContain, "user")
				So(keyList, ShouldContain, "rights")
			})

			Convey("User can use each element of config keys", func() {

				// Assert
				for _, k := range keyList {
					_, ok := ctx.Config(k)
					So(ok, ShouldBeTrue)
				}

			})

		})

		Convey("Validated Context::RawConfig", func() {

			Convey("User can unmarshal complicated configuration structures into custom type", func() {

				// Act
				rawJson := ctx.RawConfig()
				cfg := basicConfig{}
				err := json.Unmarshal([]byte(rawJson), &cfg)

				// Assert
				So(err, ShouldBeNil)
				So(cfg.Address.Ip, ShouldEqual, "192.153.25.123")
				So(cfg.Address.Port, ShouldEqual, 34245)
				So(cfg.User, ShouldEqual, "admin")
				So(cfg.Rights, ShouldResemble, []string{"admin", "logger", "runner", "reader", "writer"})

			})

		})

	})
}

func TestContextAPI_Storage(t *testing.T) {

}
