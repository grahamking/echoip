package main

import (
	"log"
	"net"
	"strings"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

const numWorkers = 4
const port = "7777"

var skipPrefix = []string{"vm", "vbox", "tun"}

func main() {
	log.Println("Start")

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	var addr string
top:
	for _, ifi := range interfaces {

		if ifi.Flags&(1<<uint(net.FlagLoopback)) == 0 ||
			ifi.Flags&net.FlagUp == 0 {
			continue
		}
		for _, s := range skipPrefix {
			if strings.HasPrefix(ifi.Name, s) {
				continue top
			}
		}

		addrs, err := ifi.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		addr = strings.Split(addrs[0].String(), "/")[0]
		break
	}

	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}

	udpaddr, err := net.ResolveUDPAddr("udp", addr+":"+port)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listening on", udpaddr)

	conn, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 1; i <= numWorkers; i++ {
		go worker(i, conn, db, wg)
	}

	wg.Wait()

	conn.Close()
	log.Println("End")
}

func worker(id int, conn *net.UDPConn, db *geoip2.Reader, wg sync.WaitGroup) {
	log.Println("Worker", id, "started")

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

		outMsg = formatMsg(remoteAddr, record)
		log.Print(id, " ", outMsg)

		_, err = conn.WriteToUDP([]byte(outMsg), remoteAddr)
		if err != nil {
			log.Println("ERROR", err)
			continue
		}
	}

	log.Println("Worker", id, "end")
	wg.Done()
}

func formatMsg(remoteAddr *net.UDPAddr, record *geoip2.City) string {

	pieces := make([]string, 1, 4)
	pieces[0] = remoteAddr.IP.String()

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
