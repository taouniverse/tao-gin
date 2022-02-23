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

package tao_gin

import (
	"github.com/gin-gonic/gin"
	"github.com/taouniverse/tao"
)

// Engine of gin implements http.Handler
var Engine *gin.Engine

// setup with gin config
// execute when init tao universe
func setup() (err error) {
	gin.SetMode(g.Mode)

	Engine = gin.New()
	writer := tao.GetWriter(tao.ConfigKey)
	Engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output: writer,
	}))
	Engine.Use(gin.RecoveryWithWriter(writer))

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
