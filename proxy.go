package proxima

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Server struct {
	// Target url, for example: https://www.google.com
	Target string

	targetURL         *url.URL
	requestListeners  []func(req *http.Request)
	responseListeners []func(res *http.Response)
	reverseProxy      *httputil.ReverseProxy
}

func (s *Server) OnRequest(callback func(req *http.Request)) {
	s.requestListeners = append(s.requestListeners, callback)
}

func (s *Server) OnResponse(callback func(res *http.Response)) {
	s.responseListeners = append(s.responseListeners, callback)
}

func (s *Server) director(req *http.Request) {
	// Set referer to the original target
	if referer := req.Referer(); referer != "" {
		referer = strings.Replace(referer, req.Host, s.targetURL.Host, 1)
		req.Header.Set("Referer", referer)
	}

	// Set origin to the original target
	if origin := req.Header.Get("Origin"); origin != "" {
		origin = strings.Replace(origin, req.Host, s.targetURL.Host, 1)
		req.Header.Set("Origin", origin)
	}

	req.URL.Scheme = s.targetURL.Scheme
	req.URL.Host = s.targetURL.Host
	req.Host = s.targetURL.Host

	if _, ok := req.Header["User-Agent"]; !ok {
		// explicitly disable User-Agent so it's not set to default value
		req.Header.Set("User-Agent", "")
	}

	req.Header.Del("Accept-Encoding")
	req.Header.Del("Content-Encoding")

	for _, listener := range s.requestListeners {
		listener(req)
	}
}

func (s *Server) modifyResponse(res *http.Response) (err error) {
	location, err := res.Location()
	if err != nil && err != http.ErrNoLocation {
		return err
	} else if err == nil {
		location.Scheme = ""
		location.Host = ""
		res.Header.Set("Location", location.String())
	}

	for _, listener := range s.responseListeners {
		listener(res)
	}

	return nil
}

func (s *Server) HandleRequests(w http.ResponseWriter, r *http.Request) {
	s.reverseProxy.ServeHTTP(w, r)
}

func (s *Server) StartProxima() (err error) {
	if s.targetURL, err = url.Parse(s.Target); err != nil {
		return err
	}

	s.reverseProxy = &httputil.ReverseProxy{
		Director:       s.director,
		ModifyResponse: s.modifyResponse,
	}

	return nil
}
