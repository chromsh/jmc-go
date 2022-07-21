package jmc

type Career string

const (
	CareerDocomo   Career = "docomo"
	CareerAU       Career = "au"
	CareerSoftbank Career = "softbank"
	CareerRakuten  Career = "rakuten"
	CareerUnknown  Career = "unknown"
)

func (c Career) String() string {
	return string(c)
}

type Data struct {
	IPv4 []string `json:"ipv4"`
	IPv6 []string `json:"ipv6"`
}

type Node4 struct {
	IP     int
	Career Career
	// ip range start?
	IsBegin bool
}
