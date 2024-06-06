package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/motemen/go-loghttp"
)

type ServerConfig struct {
	From string `yaml:"from" env:"GO_PROXY_FROM" env-default:"localhost:8090"`
	To   string `yaml:"to" env:"GO_PROXY_TO" env-default:"localhost:9000"`
}

type Server struct {
	proxy  *httputil.ReverseProxy
	url    *url.URL
	logger Logger
}

func main() {
	fromFlag := flag.String("from", "", "proxy's listening address")
	toFlag := flag.String("to", "", "the address this proxy will forward to")
	verboseFlag := flag.Bool("v", true, "more logs")
	flag.Parse()

	var l Logger = emptyLog{}
	if *verboseFlag {
		l = logger{}
		http.DefaultTransport = loghttp.DefaultTransport
	}

	var config ServerConfig
	err := cleanenv.ReadConfig("config.yml", &config)
	if err != nil {
		l.Log("Failed to load config.")
	}
	if *fromFlag != "" {
		config.From = *fromFlag
	}
	if *toFlag != "" {
		config.To = *toFlag
	}
	l.Log("Created server config:", config)

	startProxy(config.From, config.To, l)
}

func startProxy(from, to string, l Logger) {
	toUrl := parseToUrl(to)
	proxy := httputil.NewSingleHostReverseProxy(toUrl)
	proxy.Transport = http.DefaultTransport
	l.Log("starting proxy server on: ", from)
	l.Log("redirecting to: ", to)

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
