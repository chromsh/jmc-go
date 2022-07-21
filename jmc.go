package jmc

import (
	"encoding/json"
	"io/ioutil"

	"github.com/seancfoley/ipaddress-go/ipaddr"
)

var trie ipaddr.IPv4AddressTrie
var careerMap map[string]Career

func init() {
	type jsonMap struct {
		career Career
		path   string
	}
	careerMap = make(map[string]Career)
	careers := []jsonMap{
		{career: "docomo", path: "./japan-mobile-career-ip-address/data/docomo.json"},
		{career: "au", path: "./japan-mobile-career-ip-address/data/au.json"},
		{career: "softbank", path: "./japan-mobile-career-ip-address/data/softbank.json"},
		{career: "rakuten", path: "./japan-mobile-career-ip-address/data/rakuten.json"},
	}

	for _, c := range careers {
		d, err := ioutil.ReadFile(c.path)
		if err != nil {
			panic(err)
		}
		var jsonData Data
		err = json.Unmarshal(d, &jsonData)
		if err != nil {
			panic(err)
		}
		for _, ip := range jsonData.IPv4 {
			trie.Add(ipaddr.NewIPAddressString(ip).GetAddress().ToIPv4())
			careerMap[ip] = c.career
		}
	}
}

func DetectCareer(ip string) (Career, bool) {
	ipAddr := ipaddr.NewIPAddressString(ip).GetAddress()
	if !ipAddr.IsIPv4() {
		return CareerUnknown, false
	}
	res := trie.LongestPrefixMatch(ipAddr.ToIPv4())
	if res == nil {
		return CareerUnknown, false
	}
	career, ok := careerMap[res.String()]
	if !ok {
		return CareerUnknown, false
	}
	return career, true
}
