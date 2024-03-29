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
	"github.com/stretchr/testify/assert"
	"github.com/taouniverse/tao"
	"net/http"
	"testing"
)

func TestTao(t *testing.T) {
	err := tao.SetConfigPath("./test.yaml")
	assert.Nil(t, err)

	group := Engine.Group("/")

	group.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"msg": "hello tao-gin",
		})
	})

	group.GET("/index", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{
			"IP": context.ClientIP(),
		})
	})

	err = tao.Run(nil, nil)
	assert.Nil(t, err)
}
