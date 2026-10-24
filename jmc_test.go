package jmc

import (
	"encoding/binary"
	"encoding/json"
	"net"
	"testing"
)

func TestDetectCareer(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name  string
		args  args
		want  Career
		want1 bool
	}{
		{"not exist", args{"0.0.0.1"}, CareerUnknown, false},
		{"docomo", args{"1.66.96.1"}, CareerDocomo, true},
		{"au", args{"106.132.0.0"}, CareerAU, true},
		{"softbank", args{"126.32.81.255"}, CareerSoftbank, true},
		{"rakuten", args{"101.102.61.234"}, CareerRakuten, true},
		{"network address", args{"1.66.96.0"}, CareerDocomo, true},
		{"broadcast address", args{"1.66.103.255"}, CareerDocomo, true},
		{"network address -1", args{"1.66.95.255"}, CareerUnknown, false},
		{"adjacent block", args{"1.66.104.0"}, CareerDocomo, true},
		{"broadcast address +1", args{"1.66.106.0"}, CareerUnknown, false},
		{"docomo 2025 addition", args{"49.109.224.0"}, CareerDocomo, true},
		{"docomo 2025 adjacent additions", args{"1.73.243.255"}, CareerDocomo, true},
		{"docomo 2026 February addition", args{"49.109.68.0"}, CareerDocomo, true},
		{"docomo 2026 March addition", args{"49.109.240.0"}, CareerDocomo, true},
		{"docomo 2026 June addition", args{"49.96.64.0"}, CareerDocomo, true},
		{"docomo 2026 June last address", args{"49.96.87.255"}, CareerDocomo, true},
		{"docomo last address +1", args{"49.96.88.0"}, CareerUnknown, false},
		{"docomo mail destination", args{"203.138.180.0"}, CareerUnknown, false},
		{"docomo mail source", args{"203.138.203.0"}, CareerUnknown, false},
		{"docomo shared mail server", args{"203.138.204.160"}, CareerUnknown, false},
		{"docomo new mail source", args{"49.102.168.0"}, CareerUnknown, false},
		{"docomo unconfirmed scheduled addition", args{"49.109.80.0"}, CareerUnknown, false},
		{"docomo 2027 January plan", args{"27.230.0.0"}, CareerUnknown, false},
		{"docomo 2027 January small plan", args{"27.230.160.0"}, CareerUnknown, false},
		{"docomo 2027 February plan", args{"49.105.192.0"}, CareerUnknown, false},
		{"au before exclusions", args{"106.133.255.255"}, CareerAU, true},
		{"au first excluded", args{"106.134.0.0"}, CareerUnknown, false},
		{"au second excluded", args{"106.135.0.0"}, CareerUnknown, false},
		{"au last excluded", args{"106.135.255.255"}, CareerUnknown, false},
		{"au after exclusions", args{"106.136.0.0"}, CareerUnknown, false},
		{"empty input", args{""}, CareerUnknown, false},
		{"invalid input", args{"not-an-ip"}, CareerUnknown, false},
		{"invalid octet", args{"256.1.2.3"}, CareerUnknown, false},
		{"invalid punctuation", args{"1.2.3.4.5"}, CareerUnknown, false},
		{"host and port", args{"1.66.96.1:443"}, CareerUnknown, false},
		{"IPv6", args{"240a:61::1"}, CareerUnknown, false},
		{"mapped IPv6", args{"::ffff:1.66.96.1"}, CareerUnknown, false},
		{"private IPv4", args{"192.168.1.1"}, CareerUnknown, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := DetectCareer(tt.args.ip)
			if got != tt.want {
				t.Errorf("DetectCareer() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DetectCareer() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

// Compare the trie with an independent standard-library lookup at both ends
// of every embedded IPv4 block and immediately outside each block.
func TestEmbeddedIPv4Boundaries(t *testing.T) {
	type entry struct {
		career Career
		block  *net.IPNet
		bits   int
	}
	var entries []entry
	for _, career := range []Career{CareerDocomo, CareerAU, CareerSoftbank, CareerRakuten} {
		content, err := files.ReadFile("data/" + career.String() + ".json")
		if err != nil {
			t.Fatal(err)
		}
		var d data
		if err := json.Unmarshal(content, &d); err != nil {
			t.Fatal(err)
		}
		for _, cidr := range d.IPv4 {
			ip, block, err := net.ParseCIDR(cidr)
			if err != nil || ip.To4() == nil || block.String() != cidr {
				t.Fatalf("invalid or noncanonical IPv4 CIDR %q: %v", cidr, err)
			}
			bits, _ := block.Mask.Size()
			for _, other := range entries {
				if block.Contains(other.block.IP) || other.block.Contains(block.IP) {
					t.Fatalf("overlapping ranges: %s/%s and %s/%s", career, block, other.career, other.block)
				}
			}
			entries = append(entries, entry{career, block, bits})
		}
	}
	for _, e := range entries {
		first := binary.BigEndian.Uint32(e.block.IP.To4())
		last := first | ^binary.BigEndian.Uint32(net.IP(e.block.Mask).To4())
		points := []uint32{first, last}
		if first > 0 {
			points = append(points, first-1)
		}
		if last < ^uint32(0) {
			points = append(points, last+1)
		}
		for _, point := range points {
			ip := make(net.IP, net.IPv4len)
			binary.BigEndian.PutUint32(ip, point)
			want, found, best := CareerUnknown, false, -1
			for _, other := range entries {
				if other.block.Contains(ip) && other.bits > best {
					want, found, best = other.career, true, other.bits
				}
			}
			got, ok := DetectCareer(ip.String())
			if got != want || ok != found {
				t.Errorf("DetectCareer(%q) = (%s, %v), want (%s, %v)", ip, got, ok, want, found)
			}
		}
	}
}

func BenchmarkDetectCareer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DetectCareer("1.66.96.1")
	}
}
