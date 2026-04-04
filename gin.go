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

package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/taouniverse/tao"
)

/**
import _ "github.com/taouniverse/tao-gin"
*/

// G is the global config instance for tao-gin
var G = &Config{}

// Factory is the global factory instance for managing gin.Engine
var Factory *tao.BaseFactory[*gin.Engine]

func init() {
	var err error
	Factory, err = tao.Register(ConfigKey, G, NewGin)
	if err != nil {
		panic(err.Error())
	}
}

// Engine returns the default gin engine instance
func Engine() (*gin.Engine, error) {
	return Factory.Get(G.GetDefaultInstanceName())
}

// GetEngine returns the gin engine instance with specified name
func GetEngine(name string) (*gin.Engine, error) {
	if name == "" {
		return Engine()
	}
	return Factory.Get(name)
}

// NewGin creates a new gin engine instance
func NewGin(name string, config InstanceConfig) (*gin.Engine, func() error, error) {
	gin.SetMode(config.Mode)

	engine := gin.New()
	writer := tao.GetWriter(config.Writer)
	if writer == nil {
		return nil, nil, tao.NewError(tao.ParamInvalid, "gin: writer of %q not found", config.Writer)
	}
	engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output: writer,
	}))
	engine.Use(gin.RecoveryWithWriter(writer))

	if config.HTMLPattern != "" {
		engine.LoadHTMLGlob(config.HTMLPattern)
	}
	if config.StaticPath != "" {
		engine.Static("/static", config.StaticPath)
	}

	err := engine.SetTrustedProxies(config.TrustProxies)
	if err != nil {
		return nil, nil, err
	}

	closer := func() error { return nil }

	return engine, closer, nil
}
