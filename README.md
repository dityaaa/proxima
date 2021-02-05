# Proxima 7
One of the reasons this reverse proxy is made because I need to modify client request and
response from the server. Feel free to contribute 😉

## Usage
Using Proxima 7 is as simple as this:
```go
package main

import (
    "github.com/dityaaa/proxima"
    "net/http"
)

func main() {
    proxy := proxima.New("http://example.com/")
    http.HandleFunc("/", proxy.HandleRequests)
    http.ListenAndServe(":80", nil)
}
```
After that, you can start browse `localhost`. You will get respond like the target URL you have entered.

To modify request or response, you can use `OnRequest` and `OnResponse`
```go
proxy.OnRequest(func(req *http.Request) {
    // do anything to request
})

proxy.OnResponse(func(res *http.Response) {
    // do anything to response
})
```


[golang]: https://go.dev/
