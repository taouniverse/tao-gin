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
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/taouniverse/tao"
	"net/http"
	"net/http/pprof"
)

// ConfigKey for this repo
const ConfigKey = "gin"

// Config implements tao.Config
type Config struct {
	Schema       string   `json:"schema"`
	Host         string   `json:"host"`
	Port         int      `json:"port"`
	Listen       string   `json:"listen"`
	Mode         string   `json:"mode"`
	TrustProxies []string `json:"trust_proxies"`
	HTMLPattern  string   `json:"html_pattern"`
	StaticPath   string   `json:"static_path"`
	Pprof        *Pprof   `json:"pprof"`
	RunAfters    []string `json:"run_after,omitempty"`
}

// Pprof config of gin
type Pprof struct {
	Enable bool   `json:"enable"`
	Prefix string `json:"prefix"`
}

var defaultPprof = &Pprof{
	Enable: true,
	Prefix: "/pprof",
}

var defaultGin = &Config{
	Schema:    "http",
	Host:      "localhost",
	Port:      8080,
	Listen:    "127.0.0.1",
	Mode:      gin.DebugMode,
	Pprof:     defaultPprof,
	RunAfters: []string{},
}

// Name of Config
func (g *Config) Name() string {
	return ConfigKey
}

// ValidSelf with some default values
func (g *Config) ValidSelf() {
	if g.Schema != "http" && g.Schema != "https" {
		g.Schema = "http"
	}
	if g.Host == "" {
		g.Host = defaultGin.Host
	}
	if g.Port == 0 {
		g.Port = defaultGin.Port
	}
	if g.Listen == "" {
		g.Listen = defaultGin.Listen
	}
	if g.Mode == "" {
		g.Mode = defaultGin.Mode
	}
	if g.Pprof == nil {
		g.Pprof = defaultPprof
	} else {
		if g.Pprof.Enable {
			if g.Pprof.Prefix == "" {
				g.Pprof.Prefix = defaultPprof.Prefix
			}
		}
	}
	if g.RunAfters == nil {
		g.RunAfters = defaultGin.RunAfters
	}
}

// ToTask transform itself to Task
func (g *Config) ToTask() tao.Task {
	return tao.NewTask(
		ConfigKey,
		func(ctx context.Context, param tao.Parameter) (tao.Parameter, error) {
			// non-block check
			select {
			case <-ctx.Done():
				return param, tao.NewError(tao.ContextCanceled, "%s: context has been canceled", ConfigKey)
			default:
			}
			if Engine == nil {
				return param, tao.NewError(tao.Unknown, "%s: engine is nil", ConfigKey)
			}
			// gin middlewares after
			if g.Pprof != nil && g.Pprof.Enable {
				pprofOption(Engine, g)
			}
			// gin run
			tao.Add(1)
			go func() {
				defer tao.Done()
				err := Engine.Run(fmt.Sprintf("%s:%d", g.Listen, g.Port))
				if err != nil {
					tao.Panic(err)
				}
			}()
			return param, nil
		})
}

// RunAfter defines pre task names
func (g *Config) RunAfter() []string {
	return g.RunAfters
}

/**
middlewares
*/

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
