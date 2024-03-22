package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/motemen/go-loghttp"
)

type Server struct {
	proxy  *httputil.ReverseProxy
	url    *url.URL
	logger Logger
}

func main() {
	fromAddr := flag.String("from", ":8090", "proxy's listening address")
	toAddr := flag.String("to", "http://localhost:9000", "the address this proxy will forward to")
	verbose := flag.Bool("v", true, "more logs")
	flag.Parse()

	var l Logger = emptyLog{}
	if *verbose {
		l = logger{}
		http.DefaultTransport = loghttp.DefaultTransport
	}

	startProxy(*fromAddr, *toAddr, l)
}

func startProxy(from, to string, l Logger) {
	toUrl := parseToUrl(to)
	proxy := httputil.NewSingleHostReverseProxy(toUrl)
	proxy.Transport = http.DefaultTransport
	l.Log("starting proxy server from ", from)
	l.Log("redirect to ", to)

	proxy.ErrorLog = l.GetLogger()
	t := Server{
		proxy:  proxy,
		url:    toUrl,
		logger: l,
	}

	if err := http.ListenAndServe(from, t); err != nil {
		l.Log("ListenAndServe: ", err)
	}
}

func (t Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	r.Host = t.url.Host

	if method := r.Header.Get("x-http-method-override"); method != "" {
		r.Method = method
	}

	t.proxy.ServeHTTP(rw, r)
}

func parseToUrl(addr string) *url.URL {
	toUrl, err := url.Parse(addr)
	if err != nil {
		log.Fatal(err)
	}
	return toUrl
}

type Logger interface {
	Log(...any)
	GetLogger() *log.Logger
}

type emptyLog struct{}

func (emptyLog) Log(...any) {
	// do nothing
}

func (emptyLog) GetLogger() *log.Logger {
	return nil
}

type logger struct{}

func (logger) Log(args ...any) {
	log.Print(args...)
}

func (logger) GetLogger() *log.Logger {
	return log.Default()
}
