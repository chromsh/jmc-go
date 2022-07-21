# jmc-go
IP Address determination for japanese mobile careers


# usage

```go
import (
	"fmt"

	"github.com/chromsh/jmc-go"
)

func main() {
	career, ok := jmc.DetectCareer("203.138.180.0")
	if ok {
		fmt.Println(career)
	}
}
```

# benchmark

- DetectCareer
```
$ go test -benchmem -run=^$ -bench ^BenchmarkDetectCareer$ github.com/chromsh/jmc-go -count=1 -cpu=1
goos: linux
goarch: amd64
pkg: github.com/chromsh/jmc-go
cpu: AMD Ryzen 7 3700X 8-Core Processor
BenchmarkDetectCareer      93481             12863 ns/op            1601 B/op         40 allocs/op
PASS
ok      github.com/chromsh/jmc-go       1.371s
```