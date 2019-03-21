/**********************************************************\
|                                                          |
|                          hprose                          |
|                                                          |
| Official WebSite: http://www.hprose.com/                 |
|                   http://www.hprose.org/                 |
|                                                          |
\**********************************************************/
/**********************************************************\
 *                                                        *
 * rpc/http_service.go                                    *
 *                                                        *
 * hprose http service for Go.                            *
 *                                                        *
 * LastModified: Nov 24, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"github.com/hprose/hprose-golang/util"
	. "github.com/hprose/hprose-golang/rpc"
	"net"
	"github.com/opentracing/opentracing-go"
	"fmt"
	"github.com/opentracing/opentracing-go/ext"
)

// OkHTTPService is the hprose http service
type OkHTTPService struct {
	BaseHTTPService
	contextPool sync.Pool
}

type sendHeaderEvent interface {
	OnSendHeader(context *HTTPContext)
}

type sendHeaderEvent2 interface {
	OnSendHeader(context *HTTPContext) error
}

func httpFixArguments(args []reflect.Value, context ServiceContext) {
	i := len(args) - 1
	switch args[i].Type() {
	case httpContextType:
		if c, ok := context.(*HTTPContext); ok {
			args[i] = reflect.ValueOf(c)
		}
	case httpRequestType:
		if c, ok := context.(*HTTPContext); ok {
			args[i] = reflect.ValueOf(c.Request)
		}
	default:
		DefaultFixArguments(args, context)
	}
}

// InitHTTPService initializes OkHTTPService
func (service *OkHTTPService) InitHTTPService() {
	service.InitBaseHTTPService()
	service.contextPool = sync.Pool{
		New: func() interface{} { return new(HTTPContext) },
	}
	service.FixArguments = httpFixArguments
}

func (service *OkHTTPService) acquireContext() (context *HTTPContext) {
	return service.contextPool.Get().(*HTTPContext)
}

func (service *OkHTTPService) releaseContext(context *HTTPContext) {
	service.contextPool.Put(context)
}

func (service *OkHTTPService) xmlFileHandler(
	response http.ResponseWriter, request *http.Request,
	path string, context []byte) bool {
	if context == nil || strings.ToLower(request.URL.Path) != path {
		return false
	}
	if request.Header.Get("if-modified-since") == service.LastModified &&
		request.Header.Get("if-none-match") == service.Etag {
		response.WriteHeader(304)
	} else {
		contentLength := len(context)
		header := response.Header()
		header.Set("Last-Modified", service.LastModified)
		header.Set("Etag", service.Etag)
		header.Set("Content-Type", "text/xml")
		header.Set("Content-Length", util.Itoa(contentLength))
		response.Write(context)
	}
	return true
}

func (service *OkHTTPService) crossDomainXMLHandler(
	response http.ResponseWriter, request *http.Request) bool {
	path := "/crossdomain.xml"
	context := service.BaseHTTPService.CrossDomainXMLContent()
	return service.xmlFileHandler(response, request, path, context)
}

func (service *OkHTTPService) clientAccessPolicyXMLHandler(
	response http.ResponseWriter, request *http.Request) bool {
	path := "/clientaccesspolicy.xml"
	context := service.BaseHTTPService.ClientAccessPolicyXMLFile()
	return service.xmlFileHandler(response, request, path, []byte(context))
}

func (service *OkHTTPService) fireSendHeaderEvent(
	context *HTTPContext) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = NewPanicError(e)
		}
	}()
	switch event := service.Event.(type) {
	case sendHeaderEvent:
		event.OnSendHeader(context)
	case sendHeaderEvent2:
		err = event.OnSendHeader(context)
	}
	return err
}

func (service *OkHTTPService) sendHeader(context *HTTPContext) (err error) {
	if err = service.fireSendHeaderEvent(context); err != nil {
		return err
	}
	header := context.Response.Header()
	header.Set("Content-Type", "text/plain")
	if service.P3P {
		header.Set("P3P",
			`CP="CAO DSP COR CUR ADM DEV TAI PSA PSD IVAi IVDi `+
				`CONi TELo OTPi OUR DELi SAMi OTRi UNRi PUBi IND PHY ONL `+
				`UNI PUR FIN COM NAV INT DEM CNT STA POL HEA PRE GOV"`)
	}
	if service.CrossDomain {
		origin := context.Request.Header.Get("origin")
		if origin != "" && origin != "null" {
			if len(service.AccessControlAllowOrigins) == 0 ||
				service.AccessControlAllowOrigins[origin] {
				header.Set("Access-Control-Allow-Origin", origin)
				header.Set("Access-Control-Allow-Credentials", "true")
			}
		} else {
			header.Set("Access-Control-Allow-Origin", "*")
		}
	}
	return nil
}

func readAllFromHTTPRequest(request *http.Request) ([]byte, error) {
	if request.ContentLength > 0 {
		data := make([]byte, request.ContentLength)
		_, err := io.ReadFull(request.Body, data)
		return data, err
	}
	if request.ContentLength < 0 {
		return ioutil.ReadAll(request.Body)
	}
	return nil, nil
}

// ServeHTTP is the hprose http handler method
func (service *OkHTTPService) ServeHTTP(
	response http.ResponseWriter, request *http.Request) {
	if service.clientAccessPolicyXMLHandler(response, request) ||
		service.crossDomainXMLHandler(response, request) {
		return
	}
	context := service.acquireContext()
	context.InitHTTPContext(service, response, request)
	_,tracer:=InitZipkin()
	// Try to join to a trace propagated in `req`.
	req:=request
	wireContext, err := tracer.Extract(
		opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
	if err != nil {
		fmt.Printf("error encountered while trying to extract span: %+v\n", err)
	}

	// create span
	operationName := serviceName
	span := tracer.StartSpan(operationName, ext.RPCServerOption(wireContext))
	defer span.Finish()

	// store span in context
	ctx := opentracing.ContextWithSpan(req.Context(), span)

	// update request context to include our new span
	req = req.WithContext(ctx)

	var resp []byte
	err = service.sendHeader(context)
	if err == nil {
		switch request.Method {
		case "GET":
			if service.GET {
				resp = service.DoFunctionList(context)
			} else {
				response.WriteHeader(403)
			}
		case "POST":
			var req []byte
			if req, err = readAllFromHTTPRequest(request); err == nil {
				resp = service.Handle(req, context)
			}
		}
	}
	if err != nil {
		resp = service.EndError(err, context)
	}
	service.releaseContext(context)
	response.Header().Set("Content-Length", util.Itoa(len(resp)))
	response.Write(resp)
}

var stringType = reflect.TypeOf("")
var errorType = reflect.TypeOf((*error)(nil)).Elem()
var interfaceType = reflect.TypeOf((*interface{})(nil)).Elem()
var contextType = reflect.TypeOf((*Context)(nil)).Elem()
var serviceContextType = reflect.TypeOf((*ServiceContext)(nil)).Elem()
var httpContextType = reflect.TypeOf((*HTTPContext)(nil))
var httpRequestType = reflect.TypeOf((*http.Request)(nil))
var socketContextType = reflect.TypeOf((*SocketContext)(nil))
var netConnType = reflect.TypeOf((*net.Conn)(nil)).Elem()


// NewHTTPService is the constructor of OkHTTPService
func CreateHTTPService() (service *OkHTTPService) {
	service = new(OkHTTPService)
	service.InitHTTPService()
	return
}
