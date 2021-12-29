// Copyright 2021
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

package tao_gin

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/taouniverse/tao"
)

// Engine of gin implements http.Handler
var Engine *gin.Engine

/**
import _ "github.com/taouniverse/tao_gin"
*/
func init() {
	// 1. transfer config bytes to object
	g := new(GinConfig)
	bytes, err := tao.GetConfigBytes(ConfigKey)
	if err != nil {
		g = g.Default().(*GinConfig)
	} else {
		err = json.Unmarshal(bytes, &g)
		if err != nil {
			panic(err)
		}
	}

	// gin config
	g.ValidSelf()

	// 2. set object to tao
	err = tao.SetConfig(ConfigKey, g)
	if err != nil {
		panic(err)
	}

	// gin setup
	err = setup(g)
	if err != nil {
		panic(err)
	}
}

// setup with gin config
func setup(g *GinConfig) (err error) {
	gin.SetMode(g.Mode)

	Engine = gin.New()
	Engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output: tao.TaoWriter,
	}))
	Engine.Use(gin.RecoveryWithWriter(tao.TaoWriter))
	// html & static
	if g.HtmlPattern != "" {
		Engine.LoadHTMLGlob(g.HtmlPattern)
	}
	if g.StaticPath != "" {
		Engine.Static("/static", g.StaticPath)
	}

	err = Engine.SetTrustedProxies(g.TrustProxies)
	if err != nil {
		return err
	}

	return
}
