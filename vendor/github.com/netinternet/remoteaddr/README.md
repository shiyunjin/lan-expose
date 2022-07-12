# remoteaddr
Go http real ip header parser module

A forwarders such as a reverse proxy or Cloudflare find the real IP address from the requests made to the http server behind it. Local IP addresses and CloudFlare ip addresses are defined by default within the module. It is possible to define more forwarder IP addresses.

In Turkey, it is obligatory to keep the port information of IP addresses shared with cgnat by the law no 5651. For this reason, dst port information is also given along with the IP addresses. **If the IP address is behind a proxy, the dst port information is returned as -1.**

## Usage

```
go get -u github.com/netinternet/remoteaddr
```

```go
// remoteaddr.Parse().IP(*http.Request) return to string IPv4 or IPv6 address
```

## Example

Run a simple web server and get the real IP address to string format

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/netinternet/remoteaddr"
)

func root(w http.ResponseWriter, r *http.Request) {
	ip, port := remoteaddr.Parse().IP(r)
	fmt.Fprintf(w, "Your IP address is "+ip+" and dst port "+port)
}

func main() {
	http.HandleFunc("/", root)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

```

## Example 2 (Nginx or another web service forwarder address)

**AddForwarders([]string{"8.8.8.0/24"})** = Add a new multiple forwarder prefixes

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/netinternet/remoteaddr"
)

func root(w http.ResponseWriter, r *http.Request) {
	ip, port := remoteaddr.Parse().AddForwarders([]string{"8.8.8.0/24"}).IP(r)
	fmt.Fprintf(w, "Your IP address is "+ip+" and dst port "+port)
}

func main() {
	http.HandleFunc("/", root)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

```

## Example 3 (Add an alternative header for real IP address)

**AddHeaders([]string{"True-Client-IP"})** = Add a new multiple real ip headers

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/netinternet/remoteaddr"
)

func root(w http.ResponseWriter, r *http.Request) {
	ip, port := remoteaddr.Parse().AddHeaders([]string{"True-Client-IP"}).IP(r)
	fmt.Fprintf(w, "Your IP address is "+ip+" and dst port "+port)
}

func main() {
	http.HandleFunc("/", root)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

```
