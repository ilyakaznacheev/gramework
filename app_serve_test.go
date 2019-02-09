// Copyright 2017-present Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package gramework

import (
	"errors"
	"net"
	"testing"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func testBuildReqRes(method, uri string) (*fasthttp.Request, *fasthttp.Response) {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	req.Header.SetMethod(method)
	req.SetRequestURI(uri)
	return req, res
}

func TestAppServe(t *testing.T) {
	const uri = "http://test.request"

	testCases := []func(*App) (func(string,  interface{}) *App, string){
		// check GET request
		func (app *App) (func(string,  interface{}) *App, string) {
			return app.GET, GET
		},
		// check POST request
		func (app *App) (func(string,  interface{}) *App, string) {
			return app.POST, POST
		},
		// check PUT request
		func (app *App) (func(string,  interface{}) *App, string) {
			return app.PUT, PUT
		},
		// check PATCH request
		func (app *App) (func(string,  interface{}) *App, string) {
			return app.PATCH, PATCH
		},
		// check DELETE request
		func (app *App) (func(string,  interface{}) *App, string) {
			return app.DELETE, DELETE
		},
		// check HEAD request
		func (app *App) (func(string,  interface{}) *App, string) {
			return app.HEAD, HEAD
		},
		// check GET request
		func (app *App) (func(string,  interface{}) *App, string) {
			return app.OPTIONS, OPTIONS
		},
	}

	for _, test := range testCases {
		var handleOK bool 

		app := New()
		ln := fasthttputil.NewInmemoryListener()
		c := &fasthttp.Client{
			Dial: func(addr string) (net.Conn, error) {
				return ln.Dial()
			},
		}

		reg, method := test(app)

		go app.Serve(ln)

		reg("/", func() {
			handleOK = true
		})
		c.Do(testBuildReqRes(method, uri))
		
		ln.Close()

		if !handleOK {
			t.Errorf("%s request was not served correctly", method)
		}
	}
}

type testErrListener struct {
	fasthttputil.InmemoryListener
	err error
}

func (el *testErrListener) Accept() (net.Conn, error) {
	return nil, el.err
}

func TestAppServeErr(t *testing.T) {

	errTest := errors.New("test")

	el := &testErrListener{
		InmemoryListener: *fasthttputil.NewInmemoryListener(),
		err:              errTest,
	}
	app := New()
	err := app.Serve(el)
	defer el.Close()

	if err != errTest {
		t.Errorf("wrong serving error %s, but %s expected", err, errTest)
	}

}
