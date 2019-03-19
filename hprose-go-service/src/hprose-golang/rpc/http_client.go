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
 * rpc/http_client.go                                     *
 *                                                        *
 * hprose http client for Go.                             *
 *                                                        *
 * LastModified: Jan 7, 2017                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"

	hio "github.com/hprose/hprose-golang/io"
	"github.com/openzipkin-contrib/zipkin-go-opentracing/examples/middleware"
	"context"
	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin-contrib/zipkin-go-opentracing"
	. "github.com/hprose/hprose-golang/rpc"
	"hprose-golang/register/zookeeper"
	"reflect"
	"fmt"
	"os"
)

var cookieJar, _ = cookiejar.New(nil)
var httpSchemes = []string{"http", "https"}
var tcpSchemes = []string{"tcp", "tcp4", "tcp6"}
var unixSchemes = []string{"unix"}
var allSchemes = []string{"http", "https", "tcp", "tcp4", "tcp6", "unix", "ws", "wss"}

// OkHTTPClient is hprose http client
type OkHTTPClient struct {
	BaseClient
	http.Transport
	Header     http.Header
	httpClient http.Client
	limiter    Limiter
	Tracer       opentracing.Tracer
	TraceRequest middleware.RequestFunc
	uri string
}
var serviceName = "go-hprose"

func GetServiceName() string{
	return serviceName
}
func SetServiceName(zipkin_address string){
	serviceName=zipkin_address
}
const (
	hostPort = "0.0.0.0:0"
	//zipkinHTTPEndpoint = "http://localhost:9411/api/v1/spans"
	zipkinHTTPEndpoint = "http://127.0.0.1:9411/api/v1/spans"
	// Debug mode.
	debug = false
	// Base endpoint of our SVC1 service.
	svc1Endpoint = "http://localhost:8091/hello.server2/hprose"
	// same span can be set to true for RPC style spans (Zipkin V1) vs Node style (OpenTracing)
	sameSpan = true
	// make Tracer generate 128 bit traceID's for root spans.
	traceID128Bit = true
)

// NewHTTPClient is the constructor of OkHTTPClient
func NewOkHTTPClient(uri ...string) (client *OkHTTPClient) {
	client = new(OkHTTPClient)
	client.InitBaseClient()
	client.limiter.InitLimiter()
	client.httpClient.Transport = &client.Transport
	client.Header = make(http.Header)
	client.DisableCompression = true
	client.DisableKeepAlives = false
	client.MaxIdleConnsPerHost = 10
	client.httpClient.Jar = cookieJar
	if DisableGlobalCookie {
		client.httpClient.Jar, _ = cookiejar.New(nil)
	}
	client.SetURIList(uri)
	return
}

// SetURIList sets a list of server addresses
func (client *OkHTTPClient) SetURIList(uriList []string) {
	if CheckAddresses(uriList, httpSchemes) == "https" {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	client.BaseClient.SetURIList(uriList)
}

// TLSClientConfig returns the tls.Config in hprose client
func (client *OkHTTPClient) TLSClientConfig() *tls.Config {
	return client.Transport.TLSClientConfig
}

// SetTLSClientConfig sets the tls.Config
func (client *OkHTTPClient) SetTLSClientConfig(config *tls.Config) {
	client.Transport.TLSClientConfig = config
}

// MaxConcurrentRequests returns max concurrent request count
func (client *OkHTTPClient) MaxConcurrentRequests() int {
	return client.limiter.MaxConcurrentRequests
}

// SetMaxConcurrentRequests sets max concurrent request count
func (client *OkHTTPClient) SetMaxConcurrentRequests(value int) {
	client.limiter.MaxConcurrentRequests = value
}

// KeepAlive returns the keepalive status of hprose client
func (client *OkHTTPClient) KeepAlive() bool {
	return !client.DisableKeepAlives
}

// SetKeepAlive sets the keepalive status of hprose client
func (client *OkHTTPClient) SetKeepAlive(enable bool) {
	client.DisableKeepAlives = !enable
}

// Compression returns the compression status of hprose client
func (client *OkHTTPClient) Compression() bool {
	return !client.DisableCompression
}

// SetCompression sets the compression status of hprose client
func (client *OkHTTPClient) SetCompression(enable bool) {
	client.DisableCompression = !enable
}

func (client *OkHTTPClient) readAll(
	response *http.Response) (data []byte, err error) {
	if response.ContentLength > 0 {
		data = make([]byte, response.ContentLength)
		_, err = io.ReadFull(response.Body, data)
		return data, err
	}
	if response.ContentLength < 0 {
		return ioutil.ReadAll(response.Body)
	}
	return nil, nil
}

func (client *OkHTTPClient) limit() {
	client.limiter.L.Lock()
	client.limiter.Limit()
	client.limiter.L.Unlock()
}

func (client *OkHTTPClient) unlimit() {
	client.limiter.L.Lock()
	client.limiter.Unlimit()
	client.limiter.L.Unlock()
}
func (client *OkHTTPClient) sendAndReceive(
	data []byte, context1 *ClientContext) ([]byte, error) {
	client.Tracer=opentracing.GlobalTracer()
	client.TraceRequest=middleware.ToHTTPRequest(client.Tracer)
	client.limit()
	defer client.unlimit()
	//跟踪
	inst:=context1.Get("zipkinSpan")
	span,ok:=inst.(opentracing.Span)
	if !ok{
		panic("span 不存在")
	}
	inst=context1.Get("zipkinCtx")
	ctx,_:=inst.(context.Context)
	span.LogEvent("Call "+string(data))
	defer span.Finish()
	url:=client.BaseClient.URIList()[0]
	fmt.Println("url "+url)
	req, err := http.NewRequest("POST", url, hio.NewByteReader(data))
	if err != nil {
		return nil, err
	}
	for key, values := range client.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	header, ok := context1.Get("httpHeader").(http.Header)
	if ok && header != nil {
		for key, values := range header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}
	req.ContentLength = int64(len(data))
	req.Header.Set("Content-Type", "application/hprose")
	client.httpClient.Timeout = context1.Timeout
	//跟踪
	req = client.TraceRequest(req.WithContext(ctx))
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	context1.Set("httpHeader", resp.Header)
	data, err = ioutil.ReadAll(resp.Body)
	if err == nil {
		err = resp.Body.Close()
	}else{
		span.SetTag("error", err.Error())
		return nil, err
	}
	return data, err
}

var Zkaddress=""
func InitRegister(zkaddress string){
	Zkaddress=zkaddress
}
func Create(service interface{}) (client *OkHTTPClient){
	v:=reflect.ValueOf(service).Elem()
	mv := v.MethodByName("ServiceName")
	results:=mv.Call(nil)
	serviceName:=results[0].String()
	discovery:=zookeeper.GetZooKeeperServiceDiscovery(Zkaddress)
	address,err:=discovery.Discover(serviceName)
	if err!=nil{
		fmt.Println(err.Error())
	}
	fmt.Println(	"address  "+address)
	client = NewOkHTTPClient(address)
	client.SendAndReceive = client.sendAndReceive
	client.UseService(service)
	//client.BeforeFilter.Use(logFilter{}.handler)
	client.AddFilter(new(LogFilter))
	return client
}

func InitZipkin() (zipkin.Collector,opentracing.Tracer){
	collector, err := zipkin.NewHTTPCollector(zipkinHTTPEndpoint)
	if err != nil {
		fmt.Printf("unable to create Zipkin HTTP collector: %+v\n", err)
		os.Exit(-1)
	}
	recorder := zipkin.NewRecorder(collector, debug, hostPort, serviceName)
	tracer, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(sameSpan),
		zipkin.TraceID128Bit(traceID128Bit),
	)
	if err != nil {
		fmt.Printf("unable to create Zipkin tracer: %+v\n", err)
		os.Exit(-1)
	}
	opentracing.InitGlobalTracer(tracer)
	return collector,tracer
}


// LogFilter ...
type LogFilter struct {
	Prompt string
}

// InputFilter ...
func (lf LogFilter) InputFilter(data []byte, context1 Context) []byte {
	fmt.Printf("%v: %s\r\n", lf.Prompt, data)
	inst:=context1.Get("collector")
	collector,ok:=inst.(zipkin.Collector)
	if !ok{
		panic("collector 不存在")
	}
	inst=context1.Get("zipkinSpan")
	span,ok:=inst.(opentracing.Span)
	if !ok{
		panic("span 不存在")
	}
	span.Finish()
	collector.Close()
	return data
}

// OutputFilter ...
func (lf LogFilter) OutputFilter(data []byte, context1 Context) []byte {
	fmt.Printf("%v: %s\r\n", lf.Prompt, data)
	collector,_:=InitZipkin()
	span := opentracing.StartSpan("Run")
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	context1.Set("zipkinCtx",ctx)
	context1.Set("zipkinSpan",span)
	context1.Set("collector",collector)
	return data
}
