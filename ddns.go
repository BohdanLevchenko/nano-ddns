package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"net/http"
	"encoding/json"
	"github.com/miekg/dns"
)

var (
	dnsAddress  = flag.String("dns-address", ":5354", "DNS Address to listen to (TCP and UDP)")
	httpAddress = flag.String("http-address", ":8080", "HTTP Address to listen to")
	updateURL   = flag.String("update-url", "/updateip.php", "URL for updating ddns info")
	dbf         = flag.String("database", "db.json", "Database file")
	db          = map[string]string{}
)

func main() {
	flag.Parse()

	udpServer := &dns.Server{Addr: *dnsAddress, Net: "udp"}
	tcpServer := &dns.Server{Addr: *dnsAddress, Net: "tcp"}
	db = loadJson(*dbf)
	if len(db) != 0 {
		log.Println("Found addresses: ", db)
	}
	dns.HandleFunc(".", route)
	go func() {
		if err := udpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		if err := tcpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	go func() {
		http.HandleFunc(*updateURL, func(w http.ResponseWriter, r *http.Request) {
			params := r.URL.Query()
			var hostname = "modem" + params.Get("modem") + "."
			var ipAddr = params.Get("ip")
			db[hostname] = ipAddr
			log.Printf("Registered %s as %s", hostname, ipAddr)
			saveAsJson(db, *dbf)
		})
		log.Fatal(http.ListenAndServe(*httpAddress, nil))
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	udpServer.Shutdown()
	tcpServer.Shutdown()

}

func saveAsJson(v interface{}, path string) {
	fo, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer fo.Close()
	e := json.NewEncoder(fo)
	if err := e.Encode(v); err != nil {
		panic(err)
	}
}

func loadJson(path string) map[string]string {
	var db = map[string]string{}
	fo, err := os.Open(path)
	if err != nil {
		return db
	}
	e := json.NewDecoder(fo)
	if err := e.Decode(&db); err != nil {
		panic(err)
	}
	return db
}

func route(w dns.ResponseWriter, req *dns.Msg) {
	var ipAddr = db[req.Question[0].Name]
	log.Println("Resolving", req.Question[0].Name, ipAddr)
	m := new(dns.Msg)
	m.SetReply(req)

	m.Answer = make([]dns.RR, 1)

	ip, _, _ := net.ParseCIDR(ipAddr + "/24")
	m.Answer[0] = &dns.A{
		Hdr: dns.RR_Header{
			Name:   m.Question[0].Name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    0},
		A: ip}
	w.WriteMsg(m)
}
