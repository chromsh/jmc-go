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

type data struct {
	IPv4 []string `json:"ipv4"`
	IPv6 []string `json:"ipv6"`
}
