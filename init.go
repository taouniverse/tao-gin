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
	"encoding/json"
	"github.com/taouniverse/tao"
)

/**
import _ "github.com/taouniverse/tao-gin"
*/

// G config of gin
var G = new(Config)

func init() {
	err := tao.Register(ConfigKey, func() error {
		// 1. transfer config bytes to object
		bytes, err := tao.GetConfigBytes(ConfigKey)
		if err != nil {
			G = G.Default().(*Config)
		} else {
			err = json.Unmarshal(bytes, &G)
			if err != nil {
				return err
			}
		}

		// gin config
		G.ValidSelf()

		// 2. set object to tao
		err = tao.SetConfig(ConfigKey, G)
		if err != nil {
			return err
		}

		// gin setup
		return setup()
	})
	if err != nil {
		panic(err.Error())
	}
}
