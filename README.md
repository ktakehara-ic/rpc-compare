# rpc-comapre

```sh
❯ go test . -run='^&' -bench=. -benchmem
goos: darwin
goarch: arm64
pkg: github.com/ktakehara-icd/rpc-compare
BenchmarkAll/small/gRPC-gRPC-8             28088             43998 ns/op            9323 B/op        185 allocs/op
BenchmarkAll/small/gRPC-connect-8          20739             58860 ns/op           13776 B/op        189 allocs/op
BenchmarkAll/small/connect-gRPC-8          18327             64908 ns/op           12595 B/op        193 allocs/op
BenchmarkAll/small/connect-connect-8       12538             95072 ns/op           76187 B/op        151 allocs/op
BenchmarkAll/small/connect-twirp(json)-8                   10000            103084 ns/op          112050 B/op        222 allocs/op
BenchmarkAll/small/connect-twirp(proto)-8                  12950             89666 ns/op          108785 B/op        144 allocs/op
BenchmarkAll/small/twirp-twirp(json)-8                     22948             52083 ns/op           13173 B/op        217 allocs/op
BenchmarkAll/small/twirp-twirp(proto)-8                    28729             41663 ns/op            9755 B/op        127 allocs/op
BenchmarkAll/small/twirp-connect(json)-8                   18488             65097 ns/op           46071 B/op        222 allocs/op
BenchmarkAll/small/twirp-connect(proto)-8                  22747             53289 ns/op           43495 B/op        143 allocs/op
BenchmarkAll/small/REST-REST-8                             33710             40307 ns/op            8445 B/op        108 allocs/op
BenchmarkAll/big/gRPC-gRPC-8                                1702            701185 ns/op         1647469 B/op      10250 allocs/op
BenchmarkAll/big/gRPC-connect-8                             1592            848463 ns/op         1559028 B/op      10241 allocs/op
BenchmarkAll/big/connect-gRPC-8                             1297            845921 ns/op         2069412 B/op      10386 allocs/op
BenchmarkAll/big/connect-connect-8                          1539            800393 ns/op         1850423 B/op      10207 allocs/op
BenchmarkAll/big/connect-twirp(json)-8                       336           3592098 ns/op         4165508 B/op      50318 allocs/op
BenchmarkAll/big/connect-twirp(proto)-8                     1550            759838 ns/op         1959226 B/op      10187 allocs/op
BenchmarkAll/big/twirp-twirp(json)-8                         237           5013279 ns/op         4677845 B/op      50311 allocs/op
BenchmarkAll/big/twirp-twirp(proto)-8                       1401            871558 ns/op         2892000 B/op      10190 allocs/op
BenchmarkAll/big/twirp-connect(json)-8                       235           5125497 ns/op         4512303 B/op      50328 allocs/op
BenchmarkAll/big/twirp-connect(proto)-8                     1362            900644 ns/op         2837845 B/op      10216 allocs/op
BenchmarkAll/big/REST-REST-8                                 469           2519901 ns/op         3148794 B/op      10174 allocs/op
PASS
ok      github.com/ktakehara-icd/rpc-compare    35.329s
```
