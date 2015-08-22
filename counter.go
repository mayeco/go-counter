// Copyright 2015 Mario Young. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package counter

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

func configRuntime() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func init() {
	configRuntime()

	gin.SetMode("release")
	engine := gin.New()

	engine.GET("/", Create)
	engine.GET("/img/:image", Count)
	engine.GET("/inspect/:image", Inspect)

	http.Handle("/", engine)
}
