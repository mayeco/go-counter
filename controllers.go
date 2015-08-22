// Copyright 2015 Mario Young. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package counter

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/gin-gonic/gin"
)

type Counter struct {
	Hits int64
	Last time.Time
}

func Create(context *gin.Context) {
	appenginecontext := appengine.NewContext(context.Request)
	counter := Counter{
		Hits: 0,
		Last: time.Now(),
	}

	key, err := datastore.Put(
		appenginecontext,
		datastore.NewIncompleteKey(
			appenginecontext,
			"counters",
			nil,
		),
		&counter,
	)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		appenginecontext.Debugf("ERROR: %s", err)
		return
	}

	context.JSON(200, gin.H{"ID": key.IntID()})
}

func Count(context *gin.Context) {
	appenginecontext := appengine.NewContext(context.Request)
	image := context.Param("image")
	imageId, err := strconv.Atoi(image)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		appenginecontext.Debugf("%s", err)
		return
	}

	var counter Counter
	key := datastore.NewKey(appenginecontext, "counters", "", int64(imageId), nil)
	if err := datastore.Get(appenginecontext, key, &counter); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		appenginecontext.Debugf("%s", err)
		return
	}

	counter.Hits = counter.Hits + 1
	counter.Last = time.Now()

	_, err = datastore.Put(appenginecontext, key, &counter)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		appenginecontext.Debugf("%s", err)
		return
	}

	output, _ := base64.StdEncoding.DecodeString("R0lGODlhAQABAIAAAP///wAAACH5BAEAAAAALAAAAAABAAEAAAICRAEAOw==")
	context.Data(200, "image/gif", output)
}

func Inspect(context *gin.Context) {
	appenginecontext := appengine.NewContext(context.Request)
	image := context.Param("image")
	imageId, err := strconv.Atoi(image)
	if err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		appenginecontext.Debugf("%s", err)
		return
	}

	var counter Counter
	key := datastore.NewKey(appenginecontext, "counters", "", int64(imageId), nil)
	appenginecontext.Infof("KEY: %v", key)
	if err := datastore.Get(appenginecontext, key, &counter); err != nil {
		context.AbortWithError(http.StatusInternalServerError, err)
		appenginecontext.Debugf("%s", err)
		return
	}

	context.JSON(200, counter)

}
