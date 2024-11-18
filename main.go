package main

import (
	"flag"
	"fmt"
	"github.com/kouhin/envflag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var listenAddress = flag.String("listen-address", ":8080", "Default listening address for HTTP server")
var webroot = flag.String("webroot", "/webroot", "Root directory for static files")

type flagStringArray []string

func (arr *flagStringArray) String() string {
	return strings.Join(*arr, ",")
}

func (arr *flagStringArray) Set(value string) error {
	*arr = append(*arr, value)
	return nil
}

func isUrlPrefixInArray(arr flagStringArray, url string) bool {
	for _, prefix := range arr {
		if strings.HasPrefix(url, prefix) {
			return true
		}
	}
	return false
}

func main() {
	var proxyPrefixAllow flagStringArray
	var proxyPrefixDeny flagStringArray

	proxyPrefixLocal := flag.String("proxy-prefix-local", "", "Public URL prefix for the site")
	flag.Var(&proxyPrefixAllow, "proxy-prefix-allow", "List of URL prefixes to allow for proxying")
	flag.Var(&proxyPrefixDeny, "proxy-prefix-deny", "List of URL prefixes to deny for proxying")
	envflag.Parse()

	fileServer := http.FileServer(http.Dir(*webroot))

	http.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		fmt.Fprintf(writer, "{\"status\":\"UP\"}")
	})

	// Logic from https://github.com/jgraph/drawio/blob/dev/src/main/java/com/mxgraph/online/ProxyServlet.java
	http.HandleFunc("/proxy", func(writer http.ResponseWriter, request *http.Request) {
		urlParam := request.URL.Query().Get("url")

		if *proxyPrefixLocal != "" && strings.HasPrefix(urlParam, *proxyPrefixLocal) {
			localPath := "/" + strings.TrimLeft(strings.TrimPrefix(urlParam, *proxyPrefixLocal), "/")
			log.Printf("Mapping request for URL %s to local path %s", urlParam, localPath)
			request.URL.Path = localPath
			request.URL.RawQuery = ""
			fileServer.ServeHTTP(writer, request)
			return
		}

		if !isUrlPrefixInArray(proxyPrefixAllow, urlParam) || isUrlPrefixInArray(proxyPrefixDeny, urlParam) {
			log.Printf("Rejecting proxy request for forbidden URL %s", urlParam)
			writer.Header().Set("Content-type", "text/plain")
			writer.WriteHeader(403)
			fmt.Fprint(writer, "Access to this URL is not allowed!\n")
			return
		}

		parsedUrl, err := url.Parse(urlParam)
		if err != nil {
			log.Printf("Rejecting proxy request for unparsable URL %s", urlParam)
			writer.Header().Set("Content-type", "text/plain")
			writer.WriteHeader(400)
			fmt.Fprintf(writer, "Cannot parse URL: %s\n", err)
			return
		}

		log.Printf("Proxying URL %s", urlParam)
		proxy := httputil.NewSingleHostReverseProxy(parsedUrl)
		request.Method = "GET"
		request.URL.Host = parsedUrl.Host
		request.URL.Scheme = parsedUrl.Scheme
		request.Host = parsedUrl.Host
		proxy.ServeHTTP(writer, request)
	})

	http.Handle("/", fileServer)

	log.Printf("Listening on \"%s\"", *listenAddress)
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
