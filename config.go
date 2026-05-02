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
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/taouniverse/tao"
)

// ConfigKey for this repo
const ConfigKey = "gin"

// InstanceConfig 单实例配置
type InstanceConfig struct {
	Schema       string   `json:"schema"`
	Host         string   `json:"host"`
	Port         int      `json:"port"`
	Listen       string   `json:"listen"`
	Mode         string   `json:"mode"`
	TrustProxies []string `json:"trust_proxies"`
	HTMLPattern  string   `json:"html_pattern"`
	StaticPath   string   `json:"static_path"`
	Pprof        *Pprof   `json:"pprof"`
	Writer       string   `json:"writer"`
}

// Pprof config of gin
type Pprof struct {
	Enable bool   `json:"enable"`
	Prefix string `json:"prefix"`
}

// Config 总配置，实现 tao.MultiConfig 接口
type Config struct {
	tao.BaseMultiConfig[InstanceConfig]
	RunAfters []string `json:"run_after,omitempty" yaml:"run_after,omitempty"`
}

var defaultPprof = &Pprof{
	Enable: true,
	Prefix: "/pprof",
}

var defaultInstance = &InstanceConfig{
	Schema: "http",
	Host:   "localhost",
	Port:   8080,
	Listen: "127.0.0.1",
	Mode:   gin.DebugMode,
	Pprof:  defaultPprof,
	Writer: tao.ConfigKey,
}

// Name of Config
func (g *Config) Name() string {
	return ConfigKey
}

// ValidSelf with some default values
func (g *Config) ValidSelf() {
	for i := range g.Instances {
		instance := &g.Instances[i].Cfg
		if instance.Schema != "http" && instance.Schema != "https" {
			instance.Schema = "http"
		}
		if instance.Host == "" {
			instance.Host = defaultInstance.Host
		}
		if instance.Port == 0 {
			instance.Port = defaultInstance.Port
		}
		if instance.Listen == "" {
			instance.Listen = defaultInstance.Listen
		}
		if instance.Mode == "" {
			instance.Mode = defaultInstance.Mode
		}
		if instance.Pprof == nil {
			instance.Pprof = defaultPprof
		} else {
			if instance.Pprof.Enable {
				if instance.Pprof.Prefix == "" {
					instance.Pprof.Prefix = defaultPprof.Prefix
				}
			}
		}
		if instance.Writer == "" {
			instance.Writer = defaultInstance.Writer
		}
	}
	if g.RunAfters == nil {
		g.RunAfters = []string{}
	}
}

// ToTask transform itself to Task
func (g *Config) ToTask() tao.Task {
	return tao.NewTask(
		ConfigKey,
		func(ctx context.Context, param tao.Parameter) (tao.Parameter, error) {
			select {
			case <-ctx.Done():
				return param, tao.NewError(tao.ContextCanceled, "%s: context has been canceled", ConfigKey)
			default:
			}
			for _, inst := range g.Instances {
				name := inst.Name
				engine, err := Factory.Get(name)
				if err != nil {
					return param, err
				}
				cfg := inst.Cfg
				if cfg.Pprof != nil && cfg.Pprof.Enable {
					pprofOption(engine, &cfg)
				}
				tao.Add(1)
				go func() {
					defer tao.Done()
					err := engine.Run(fmt.Sprintf("%s:%d", cfg.Listen, cfg.Port))
					if err != nil {
						tao.Panic(err)
					}
				}()
			}
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
func pprofOption(engin *gin.Engine, inst *InstanceConfig) {
	pp := engin.Group(inst.Pprof.Prefix)

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
