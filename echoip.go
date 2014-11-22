package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	geoip2 "github.com/oschwald/geoip2-golang"
)

const numWorkers = 4

var skipPrefix = []string{"vm", "vbox", "tun"}

var port = flag.String("p", "7777", "Port")
var iface = flag.String("i", "auto", "Interface. 'auto' to let echoip choose")
var help = flag.Bool("h", false, "Show usage")

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	ifaceName, addr := findBindAddress(*iface)

	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}

	udpaddr, err := net.ResolveUDPAddr("udp", addr+":"+*port)
	if err != nil {
		log.Fatal(err)
	}
	tcpaddr, err := net.ResolveTCPAddr("tcp", addr+":"+*port)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on %s %s (tcp, udp)\n", ifaceName, udpaddr)

	conn, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		log.Fatal(err)
	}

	for i := 1; i <= numWorkers; i++ {
		go udpworker(i, conn, db)
		go tcpworker(i, listener, db)
	}

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	<-ch

	conn.Close()
	log.Println("Bye")
}

func findBindAddress(ifaceName string) (string, string) {
	ifi := selectInterface(ifaceName)
	addrs, err := ifi.Addrs()
	if err != nil {
		log.Fatal(err)
	}
	addr := strings.Split(addrs[0].String(), "/")[0]
	return ifi.Name, addr
}

func selectInterface(ifaceName string) *net.Interface {
	var ifi *net.Interface
	var err error

	if ifaceName != "auto" {
		ifi, err = net.InterfaceByName(ifaceName)
		if err != nil {
			log.Fatal(err)
		}
		return ifi
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

top:
	for _, candidate := range interfaces {

		isLoopback := candidate.Flags&(1<<uint(net.FlagLoopback)) == 0
		isDown := candidate.Flags&net.FlagUp == 0
		if isLoopback || isDown {
			//log.Printf("skipping %s. isLoopback: %t, isDown: %t\n", candidate.Name, isLoopback, isDown)
			continue
		}
		for _, s := range skipPrefix {
			if strings.HasPrefix(candidate.Name, s) {
				//log.Printf("skipping %s, has skip prefix %s\n", candidate.Name, s)
				continue top
			}
		}
		ifi = &candidate
		break
	}
	return ifi
}

func tcpworker(id int, listener *net.TCPListener, db *geoip2.Reader) {
	var conn *net.TCPConn
	var err error
	var remoteAddr *net.TCPAddr
	var record *geoip2.City
	var outMsg string

	for {
		conn, err = listener.AcceptTCP()
		if err != nil {
			log.Println("ERROR", err)
			continue
		}
		remoteAddr = conn.RemoteAddr().(*net.TCPAddr)

		record, err = db.City(remoteAddr.IP)
		if err != nil {
			log.Println("ERROR", err)
			continue
		}

		outMsg = formatMsg(remoteAddr.IP.String(), record)
		log.Print("TCP", id, " ", outMsg)

		_, err = conn.Write([]byte(outMsg))
		if err != nil {
			log.Println("ERROR", err)
		}
		conn.Close() // Server close might cause lots of TIME_WAIT
	}
}

func udpworker(id int, conn *net.UDPConn, db *geoip2.Reader) {
	//log.Println("Worker", id, "started")

	var err error
	var remoteAddr *net.UDPAddr
	var record *geoip2.City
	var outMsg string
	buf := make([]byte, 1)

	for {
		_, remoteAddr, err = conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("ERROR", err)
			continue
		}

		record, err = db.City(remoteAddr.IP)
		if err != nil {
			log.Println("ERROR", err)
			continue
		}

		outMsg = formatMsg(remoteAddr.IP.String(), record)
		log.Print("UDP", id, " ", outMsg)

		_, err = conn.WriteToUDP([]byte(outMsg), remoteAddr)
		if err != nil {
			log.Println("ERROR", err)
		}
	}
	//log.Println("Worker", id, "end")
}

func formatMsg(remoteAddr string, record *geoip2.City) string {

	pieces := make([]string, 1, 4)
	pieces[0] = remoteAddr

	city := record.City.Names["en"]
	if city != "" {
		pieces = append(pieces, city)
	}

	// Subdivisions is State in US, Province in CA
	var subs []string
	for _, sub := range record.Subdivisions {
		subs = append(subs, sub.Names["en"])
	}
	if len(subs) > 0 {
		pieces = append(pieces, subs...)
	}

	country := record.Country.Names["en"]
	if country != "" {
		pieces = append(pieces, country)
	}

	return strings.Join(pieces, ",") + "\n"
}
