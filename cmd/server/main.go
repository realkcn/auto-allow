package main

import (
	firewallCtl "auto-allow/pkg/firewallctl"
	"flag"
	"fmt"
	"github.com/ReneKroon/ttlcache/v2"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	AutoAllowChain = "AUTOALLOW"
	AllowIPChain   = "ALLOWIP"
)

var (
	displayHelp       bool
	allowRelatedInput bool
	cleanFirewall     bool
	listenPort        int
	protectPorts      string
	accessKey         string
	openPorts         string
	allowIPs          string
	durationString    string
)

func argumentFlag() {
	flag.BoolVar(&allowRelatedInput, "allowRelated", true, "allow RELATED packet on INPUT")
	flag.BoolVar(&cleanFirewall, "clean", false, "clean iptables rule")
	flag.BoolVar(&displayHelp, "h", false, "display help")
	flag.StringVar(&durationString, "t", "12h", "allow access duration")
	flag.StringVar(&accessKey, "k", "", "access key(must have)")
	flag.IntVar(&listenPort, "p", 8080, "listen port")
	flag.StringVar(&protectPorts, "protect", "", "special ports to open(default is all)")
	flag.StringVar(&openPorts, "open", "22", "ports always open")
	flag.StringVar(&allowIPs, "allow", "127.0.0.1", "IPs always allow")
}

func initFirewall() {
	if allowRelatedInput {
		firewallCtl.ExecIptables("-I", "INPUT", "-m", "state", "--state", "ESTABLISHED,RELATED", "-j", "ACCEPT")
	}
	firewallCtl.AddChain(AutoAllowChain)
	firewallCtl.AddJumpChain("INPUT", AutoAllowChain)
	firewallCtl.AddChain(AllowIPChain)
	firewallCtl.AddJumpChain(AutoAllowChain, AllowIPChain)
	//open our listen port
	firewallCtl.AddRule(AutoAllowChain, "", strconv.Itoa(listenPort), "tcp", "ACCEPT")
	//ports always open
	ports := strings.Split(openPorts, ",")
	if len(ports) > 0 {
		log.Infof("Allow port list:%s", openPorts)
		for _, port := range ports {
			if port != "" {
				firewallCtl.AddRule(AutoAllowChain, "", port, "tcp", "ACCEPT")
			}
		}
	}
	//ip always allow
	ips := strings.Split(allowIPs, ",")
	if len(ips) > 0 {
		log.Infof("Allow ip list:%s", allowIPs)
		for _, ip := range ips {
			if ip != "" {
				firewallCtl.AddRule(AutoAllowChain, ip, "", "tcp", "ACCEPT")
			}
		}
	}
	//protect ports
	ports = strings.Split(protectPorts, ",")
	if len(ports) > 0 && protectPorts != "" {
		log.Infof("Protect ports:%s", protectPorts)
		for _, port := range ports {
			if port != "" {
				firewallCtl.AddRule(AutoAllowChain, "", port, "tcp", "DROP")
			}
		}
	} else {
		//all DROP
		log.Infof("Drop all other ports")
		firewallCtl.AddRule(AutoAllowChain, "", "", "tcp", "DROP")
	}
}

func removeFirewall() {
	if allowRelatedInput {
		firewallCtl.ExecIptables("-D", "INPUT", "-m", "state", "--state", "ESTABLISHED,RELATED", "-j", "ACCEPT")
	}
	firewallCtl.RemoveJumpChain("INPUT", AutoAllowChain)
	firewallCtl.RemoveChain(AutoAllowChain)
	firewallCtl.RemoveChain(AllowIPChain)
}

func main() {
	argumentFlag()
	flag.Parse()

	if displayHelp {
		flag.Usage()
		return
	}
	if cleanFirewall {
		removeFirewall()
		return
	}
	allowDuration, err := time.ParseDuration(durationString)
	if err != nil {
		log.Fatalf("parse duration error:%s", durationString)
		return
	}

	if accessKey == "" {
		log.Fatalf("Must assign access key:use -k key")
		return
	}

	allowIpList := NewTTLList(allowDuration,
		func(key string, reason ttlcache.EvictionReason, value interface{}) {
			log.Tracef("IP %s has expired because of %s", key, reason)
			firewallCtl.RemoveRule(AllowIPChain, key, "", "tcp", "ACCEPT")
		})
	handlers := Handlers{AccessKey: accessKey, AllowIpList: allowIpList}

	initFirewall()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", handlers.DefaultPage)
	router.HandleFunc("/add", handlers.HandleAddAllowIP)
	router.HandleFunc("/remove", handlers.HandleRemoveAllowIP)
	router.HandleFunc("/get", handlers.HandleGetAllowIPs)
	log.Infof("Server start Port: %d", listenPort)
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		// Run Cleanup
		signal.Ignore(syscall.SIGTERM)
		removeFirewall()
		os.Exit(0)
	}()
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", listenPort), router)
	if err != nil {
		log.Fatalf("Server error: %s", err.Error())
	}
	removeFirewall()
}
