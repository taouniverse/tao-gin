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
	"net/http"
	"net/http/pprof"
)

// pprofOption for gin.Engine
func pprofOption(engin *gin.Engine, g *Config) {
	pp := engin.Group(g.Pprof.Prefix)

	var pprofHandler = func(h http.HandlerFunc) gin.HandlerFunc {
		return func(c *gin.Context) {
			h.ServeHTTP(c.Writer, c.Request)
		}
	}

	pp.GET("/", pprofHandler(pprof.Index))
	pp.GET("/cmdline", pprofHandler(pprof.Cmdline))
	pp.GET("/profile", pprofHandler(pprof.Profile))
	pp.POST("/symbol", pprofHandler(pprof.Symbol))
	pp.GET("/symbol", pprofHandler(pprof.Symbol))
	pp.GET("/trace", pprofHandler(pprof.Trace))
	pp.GET("/allocs", pprofHandler(pprof.Handler("allocs").ServeHTTP))
	pp.GET("/block", pprofHandler(pprof.Handler("block").ServeHTTP))
	pp.GET("/goroutine", pprofHandler(pprof.Handler("goroutine").ServeHTTP))
	pp.GET("/heap", pprofHandler(pprof.Handler("heap").ServeHTTP))
	pp.GET("/mutex", pprofHandler(pprof.Handler("mutex").ServeHTTP))
	pp.GET("/threadcreate", pprofHandler(pprof.Handler("threadcreate").ServeHTTP))
}
