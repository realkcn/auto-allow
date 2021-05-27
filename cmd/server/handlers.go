package main

import (
	firewallCtl "auto-allow/pkg/firewallctl"
	"fmt"
	"net/http"
	"strings"
)

type Handlers struct {
	AccessKey   string
	AllowIpList *TTLList
}

func (handlers *Handlers) getAllowIP(request *http.Request) string {
	remoteIP := strings.Split(request.RemoteAddr, ":")[0]
	err := request.ParseForm()
	if err != nil {
		return remoteIP
	}
	ip := request.Form["ip"]
	if len(ip) > 0 {
		return ip[0]
	} else {
		return remoteIP
	}
}

func (handlers *Handlers) verifyKey(request *http.Request) bool {
	err := request.ParseForm()
	if err != nil {
		return false
	}
	key := request.Form["key"]
	if len(key) > 0 && key[0] == accessKey {
		return true
	} else {
		return false
	}
}

func (handlers *Handlers) DefaultPage(response http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(response, "Test Page")
}

func (handlers *Handlers) HandleAddAllowIP(response http.ResponseWriter, request *http.Request) {
	if handlers.verifyKey(request) {
		remoteIP := handlers.getAllowIP(request)
		if !handlers.AllowIpList.Exist(remoteIP) {
			firewallCtl.AddRule(AllowIPChain, remoteIP, "", "tcp", "ACCEPT")
		}
		handlers.AllowIpList.Add(remoteIP)
		_, _ = fmt.Fprintf(response, "allow ip:%s", remoteIP)
	} else {
		handlers.DefaultPage(response, request)
	}
}

func (handlers *Handlers) HandleRemoveAllowIP(response http.ResponseWriter, request *http.Request) {
	if handlers.verifyKey(request) {
		remoteIP := handlers.getAllowIP(request)
		if handlers.AllowIpList.Exist(remoteIP) {
			firewallCtl.RemoveRule(AllowIPChain, remoteIP, "", "tcp", "ACCEPT")
			handlers.AllowIpList.Remove(remoteIP)
			_, _ = fmt.Fprintf(response, "remove ip:%s", remoteIP)
		} else {
			_, _ = fmt.Fprintf(response, "ip not exist:%s", remoteIP)
		}
	} else {
		handlers.DefaultPage(response, request)
	}
}

func (handlers *Handlers) HandleGetAllowIPs(response http.ResponseWriter, request *http.Request) {
	if handlers.verifyKey(request) {
		ips := handlers.AllowIpList.GetAll()
		for _, ip := range ips {
			_, _ = fmt.Fprintf(response, "%s\n", ip)
		}
	}
}
