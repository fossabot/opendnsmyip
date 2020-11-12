/*
Package publicip returns the public facing IPv4 address of the requesting client by querying servers
at OpenDNS.

Example:

	package main

	import (
	"fmt"
	"github.com/polera/publicip"
	)

	func main() {

	myIpAddr, err := publicip.GetIP()
	if err != nil {
		fmt.Printf("Error getting IP address: %s\n", err)
	} else {
		fmt.Printf("Public IP address is: %s\n", myIpAddr)
	}

	}
*/
package publicip

import (
	"fmt"

	"github.com/miekg/dns"
)

/*
GetIP returns the public facing IPv4 address of the requesting client by querying servers
at OpenDNS.
*/
func GetIP() (string, error) {
	config := dns.ClientConfig{Servers: []string{"208.67.220.220", "208.67.222.222"}, Port: "53"}
	dnsClient := new(dns.Client)
	message := new(dns.Msg)
	message.SetQuestion("myip.opendns.com.", dns.TypeA)
	message.RecursionDesired = false
	return doDNSLookup(config, dnsClient, message)
}

func doDNSLookup(config dns.ClientConfig, client *dns.Client, message *dns.Msg) (string, error) {
	err := fmt.Errorf("Error querying servers at OpenDNS")
	for _, server := range config.Servers {
		serverAddr := fmt.Sprintf("%s:%s", server, config.Port)
		response, _, cliErr := client.Exchange(message, serverAddr)
		if cliErr != nil {
			return "", fmt.Errorf("Error on DNS lookup: %w", cliErr)
		}
		if response.Rcode != dns.RcodeSuccess {
			err = fmt.Errorf("DNS call not successful. Response code: %d", response.Rcode)
		} else {
			for _, answer := range response.Answer {
				if aRecord, ok := answer.(*dns.A); ok {
					return aRecord.A.String(), nil
				}
			}
		}
	}
	return "", err
}
