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

/**
import _ "github.com/taouniverse/tao_gin"
*/
func init() {
	// 1. transfer config bytes to object
	h := new(GinConfig)
	bytes, err := tao.GetConfigBytes(ConfigKey)
	if err != nil {
		h = h.Default().(*GinConfig)
	} else {
		err = json.Unmarshal(bytes, &h)
		if err != nil {
			panic(err)
		}
	}

	// gin config
	h.ValidSelf()

	// 2. set object to tao
	err = tao.SetConfig(ConfigKey, h)
	if err != nil {
		panic(err)
	}

	// setup with gin config
	Engine = gin.New()
	Engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output: tao.TaoWriter,
	}))
	Engine.Use(gin.RecoveryWithWriter(tao.TaoWriter))

	if h.Mode != gin.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	err = Engine.SetTrustedProxies(h.TrustProxies)
	if err != nil {
		panic(err)
	}
}
