package jmc

import "testing"

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
		{"network address", args{"203.138.180.0"}, CareerDocomo, true},
		{"broadcast address", args{"203.138.180.255"}, CareerDocomo, true},
		{"network addres -1", args{"203.138.179.255"}, CareerUnknown, false},
		// 203.138.181.0/24が存在するので、その+1でテスト
		{"broadcast address +1", args{"203.138.182.0"}, CareerUnknown, false},
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

func BenchmarkDetectCareer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DetectCareer("203.138.180.0")
	}
}
