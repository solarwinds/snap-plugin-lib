// +build small

/*
 Copyright (c) 2020 SolarWinds Worldwide, LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package proxy

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

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
	jsonConfig := []byte(`{
    	"address": {
        	"ip": "192.153.25.123",
        	"port": 34245
    	},
    	"rights": ["admin", "logger", "runner", "reader", "writer"], 
    	"user": "admin"
	}`)

	ctx, cErr := NewContext(jsonConfig)

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

			var keyList []string

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
				err := json.Unmarshal(rawJson, &cfg)

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

type storedClient struct {
	count int
}

func (sc *storedClient) Inc() {
	sc.count++
}

func (sc *storedClient) Count() int {
	return sc.count
}

func TestContextAPI_Storage(t *testing.T) {
	Convey("Validate Context API for handling storage", t, func() {
		// Arrange
		emptyConfig := []byte("{}")
		ctx, cErr := NewContext(emptyConfig)

		So(cErr, ShouldBeNil)

		Convey("Validate that object of basic type may be stored in context", func() {
			// Arrange
			ctx.Store("version", "1.0.1")
			ctx.Store("apiVersion", 12)
			ctx.Store("debugMode", true)

			Convey("Validated that object of basic type may be read from context (1) via Load method", func() {
				// Act
				ver, ok := ctx.Load("version")

				// Assert
				So(ok, ShouldBeTrue)
				So(ver, ShouldHaveSameTypeAs, "")
				So(ver, ShouldEqual, "1.0.1")
			})

			Convey("Validated that object of basic type may be read from context (1) via LoadTo method", func() {
				// Act
				ver := ""
				err := ctx.LoadTo("version", &ver)

				// Assert
				So(err, ShouldBeNil)
				So(ver, ShouldEqual, "1.0.1")
			})

			Convey("Validated that object of basic type may be read from context (2) via Load method", func() {
				// Act
				ver, ok := ctx.Load("apiVersion")

				// Assert
				So(ok, ShouldBeTrue)
				So(ver, ShouldHaveSameTypeAs, 11)
				So(ver, ShouldEqual, 12)
			})

			Convey("Validated that object of basic type may be read from context (2) via LoadTo method", func() {
				// Act
				v := 0
				err := ctx.LoadTo("apiVersion", &v)

				// Assert
				So(err, ShouldBeNil)
				So(v, ShouldEqual, 12)
			})

			Convey("Validated that object of basic type may be read from context (3) via Load method", func() {
				// Act
				ver, ok := ctx.Load("debugMode")

				// Assert
				So(ok, ShouldBeTrue)
				So(ver, ShouldHaveSameTypeAs, false)
				So(ver, ShouldEqual, true)
			})

			Convey("Validated that object of basic type may be read from context (3) via LoadTo method", func() {
				// Act
				v := false
				err := ctx.LoadTo("debugMode", &v)

				// Assert
				So(err, ShouldBeNil)
				So(v, ShouldEqual, true)
			})

			Convey("Validated that object of unknown key can't be read from context", func() {
				// Act
				_, ok := ctx.Load("serverAPI")

				// Assert
				So(ok, ShouldBeFalse)
			})

		})

		Convey("Validate that object of complex type may be stored", func() {
			// Arrange
			obj := &storedClient{}
			obj.Inc()

			ctx.Store("client", obj)

			Convey("Validated that object of complex type may be read from context", func() {
				// Act
				cli, ok := ctx.Load("client")

				// Assert
				So(ok, ShouldBeTrue)
				So(cli, ShouldHaveSameTypeAs, &storedClient{})
				So(cli.(*storedClient).Count(), ShouldEqual, 1)

				// Act
				sCli := cli.(*storedClient)
				sCli.Inc()
				sCli.Inc()

				// Assert
				So(cli.(*storedClient).Count(), ShouldEqual, 3)
			})

			Convey("Validated that object of complex type may be read from context directly to complex type", func() {
				// Act
				cli := &storedClient{}
				err := ctx.LoadTo("client", &cli)

				// Assert
				So(err, ShouldBeNil)
				So(cli.Count(), ShouldEqual, 1)

				// Act
				cli.Inc()
				cli.Inc()

				// Assert
				So(cli.Count(), ShouldEqual, 3)
			})
		})
	})
}
