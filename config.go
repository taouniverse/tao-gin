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
	"context"
	"github.com/gin-gonic/gin"
	"github.com/taouniverse/tao"
)

// Engine of gin implements http.Handler
var Engine *gin.Engine

// ConfigKey for this repo
const ConfigKey = "gin"

// GinConfig implements tao.Config
type GinConfig struct {
	Host         string   `json:"host"`
	Port         string   `json:"port"`
	Listen       string   `json:"listen"`
	Mode         string   `json:"mode"`
	TrustProxies []string `json:"trust_proxies"`
	RunAfter_    []string `json:"run_after,omitempty"`
}

var defaultGin = &GinConfig{
	Host:         "localhost",
	Port:         "8080",
	Listen:       "127.0.0.1",
	Mode:         gin.DebugMode,
	TrustProxies: nil,
	RunAfter_:    []string{},
}

// Default config
func (g GinConfig) Default() tao.Config {
	return defaultGin
}

// ValidSelf with some default values
func (g *GinConfig) ValidSelf() {
	if g.Host == "" {
		g.Host = defaultGin.Host
	}
	if g.Port == "" {
		g.Port = defaultGin.Port
	}
	if g.Listen == "" {
		g.Listen = defaultGin.Listen
	}
	if g.Mode == "" {
		g.Mode = defaultGin.Mode
	}
	if g.RunAfter_ == nil {
		g.RunAfter_ = defaultGin.RunAfter_
	}
}

// ToTask transform itself to Task
func (g *GinConfig) ToTask() tao.Task {
	return tao.NewTask(ConfigKey, func(ctx context.Context, param tao.Parameter) (tao.Parameter, error) {
		// gin run
		return param, Engine.Run(g.Listen + ":" + g.Port)
	})
}

// RunAfter defines pre task names
func (g *GinConfig) RunAfter() []string {
	return g.RunAfter_
}
