# jmc-go
IP Address determination for japanese mobile careers

IPv4 data reviewed on **2026-10-24**. Docomo and au were updated from their
official public lists. SoftBank could not be fully retrieved; Rakuten's existing
ASN-derived ranges remain unverified for mobile-only use. See
[SOURCES.md](SOURCES.md) for sources, changes, deferred ranges, and limitations.
The API and JSON schema are unchanged; IPv6 data was left unchanged.

## Data synchronization and tests

The source repository is `chromsh/japan-mobile-career-ip-address`. JSON files are
copied into `data/` and embedded at build time. For a local sibling checkout:

```sh
make sync-data DATA_SOURCE=../japan-mobile-career-ip-address
make check-data DATA_SOURCE=../japan-mobile-career-ip-address
go test ./...
go test -race ./...
go vet ./...
```

The existing `make embed` fetches the remote submodule before copying. The local
synchronization targets do not fetch or change its pinned commit. Pin the
submodule to the source commit being reviewed, then copy and verify that data.

Docomo mail-only server ranges now return `CareerUnknown, false`. They are not
internet-access source ranges. The example below uses a published access range.

# usage

```go
import (
	"fmt"

	"github.com/chromsh/jmc-go"
)

func main() {
	career, ok := jmc.DetectCareer("1.66.96.1")
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
