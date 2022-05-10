package utils

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

type DNS_PROVIDER struct {
	nameserver string
	fqdn       string
	class      dns.Class
}

var GoogleDns DNS_PROVIDER = DNS_PROVIDER{
	nameserver: "ns1.google.com:53",
	fqdn:       "o-o.myaddr.l.google.com.",
	class:      dns.ClassINET,
}

var CloudflareDns DNS_PROVIDER = DNS_PROVIDER{
	nameserver: "one.one.one.one:53",
	fqdn:       "whoami.cloudflare.",
	class:      dns.ClassCHAOS,
}

var (
	ErrNoTXTRecordFound  = errors.New("no TXT record found")
	ErrTooManyAnswers    = errors.New("too many answers")
	ErrInvalidAnswerType = errors.New("invalid answer type")
	ErrTooManyTXTRecords = errors.New("too many TXT records")
	ErrIPMalformed       = errors.New("IP address malformed")
)

func GetLocalIP() ([]string, error) {
	var ips []string
	netInterfaceAddresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, netInterfaceAddress := range netInterfaceAddresses {
		networkIP, ok := netInterfaceAddress.(*net.IPNet)
		if ok && !networkIP.IP.IsLoopback() && networkIP.IP.To4() != nil {
			ip := networkIP.IP.String()
			ips = append(ips, ip)
		}
	}
	return ips, nil
}

func GetExternalIP(dnsProvider DNS_PROVIDER) (string, error) {
	c := dns.Client{
		Timeout: 5 * time.Second,
	}

	m := dns.Msg{
		MsgHdr: dns.MsgHdr{
			Opcode: dns.OpcodeQuery,
		},
		Question: []dns.Question{
			{
				Name:   dnsProvider.fqdn,
				Qtype:  dns.TypeTXT,
				Qclass: uint16(dnsProvider.class),
			},
		},
	}

	r, _, err := c.Exchange(&m, dnsProvider.nameserver)
	if err != nil {
		return "", err
	}

	L := len(r.Answer)
	if L == 0 {
		return "", ErrNoTXTRecordFound
	} else if L > 1 {
		return "", fmt.Errorf("%w: %d instead of 1", ErrTooManyAnswers, L)
	}

	answer := r.Answer[0]
	txt, ok := answer.(*dns.TXT)
	if !ok {
		return "", fmt.Errorf("%w: %T instead of *dns.TXT",
			ErrInvalidAnswerType, answer)
	}

	L = len(txt.Txt)
	if L == 0 {
		return "", ErrNoTXTRecordFound
	} else if L > 1 {
		return "", fmt.Errorf("%w: %d instead of 1", ErrTooManyTXTRecords, L)
	}
	ipString := txt.Txt[0]

	publicIP := net.ParseIP(ipString)
	if publicIP == nil {
		return "", fmt.Errorf("%w: %q", ErrIPMalformed, ipString)
	}

	return publicIP.String(), nil
}
