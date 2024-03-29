// Copyright 2022 huija
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/taouniverse/tao"
)

/**
import _ "github.com/taouniverse/tao-gin"
*/

// G config of gin
var G = new(Config)

func init() {
	err := tao.Register(ConfigKey, G, setup)
	if err != nil {
		panic(err.Error())
	}
}

// Engine of gin implements http.Handler
var Engine *gin.Engine

// setup with gin config
// execute when init tao universe
func setup() (err error) {
	gin.SetMode(G.Mode)

	Engine = gin.New()
	writer := tao.GetWriter(G.Writer)
	if writer == nil {
		return tao.NewError(tao.ParamInvalid, "gin: writer of %q not found", G.Writer)
	}
	Engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output: writer,
	}))
	Engine.Use(gin.RecoveryWithWriter(writer))

	// html & static
	if G.HTMLPattern != "" {
		Engine.LoadHTMLGlob(G.HTMLPattern)
	}
	if G.StaticPath != "" {
		Engine.Static("/static", G.StaticPath)
	}

	err = Engine.SetTrustedProxies(G.TrustProxies)
	if err != nil {
		return err
	}

	return
}
