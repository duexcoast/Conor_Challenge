# SED Coding Challenge

Conor Ney

I decided to contrast two approaches to the problem. One using the standard 
library, and the other using a more expressive regex engine.

Go's stdlib `regexp` package is powered by the [RE2](https://github.com/google/re2/wiki/Syntax)
regex engine. RE2 is much more limited in it's expressive power than most other
regex engines; namely, it does not support backreferences or lookarounds. The 
reasoning for this absence is to guarantee match time linear to the length of 
the string using any arbitrary regular expression from users. 

In order to solve the problem using simply a regular expression you would need
to utilize those missing features. In lieu of that, with the `regexp` package
solution, I just implemented a short loop to satisfy that final condition.

Alternatively, I also found a more expressive regex package, `regexp2`, which 
fully implements PCRE. This solution was drastically shorter(a 6 line function), 
although after benchmarking the two approaches, the solution using the standard 
library was approximately 15x faster.

### Benchmarks

| Benchmark        | Time per operation     | Bytes allocated per operation | Allocations per operation |
|-------------------------------|-----------------------:|------------------:|------------------:|
| BenchmarkValidateRE2-8        |          167.4 ns/op   |          0 B/op   |       0 allocs/op | 
| BenchmarkValidatePCRE-8       |         2553 ns/op     |         80 B/op   |       1 allocs/op |
| BenchmarkValidateInputRE2-8   |       120468 ns/op     |       4098 B/op   |       1 allocs/op |
| BenchmarkValidateInputPCRE-8  |      2086367 ns/op     |      73531 B/op   |    1001 allocs/op |

It's unclear to me without looking further into the matter whether this is an
inherent property of the complex regular expression, which would be surprising 
to me given the extremely small size of the input strings. I think it's more 
likely that the `regexp2` package is simply very unoptimized.

To that extent, I ran a benchmark comparing the performance of the stdlib `regexp`
package with the `regexp2` package, performing the exact same operation:

| Benchmark                       | Time per operation     | Bytes allocated per operation | Allocations per operation |
|---------------------------------|-----------------------:|------------------------------:|--------------------------:|
| BenchmarkValidateRE2-8          |          168.1 ns/op   |                      0 B/op   |               0 allocs/op | 
| BenchmarkValidatePCRESideBySide-8   |  531.2 ns/op       |                     80 B/op   |               1 allocs/op |

We're a lot closer here, but it's still 3x as slow, and looking at further 
comparisons [here](https://itnext.io/best-regexp-alternative-for-go-be42abdc1fbb)
(wish I had seen this before I picked it) it's evident that the `regexp2` package 
is among the poorer performing libraries available.

All in all, when matching a string the size of a credit card number, a better 
optimized regex engine should be able to perform excellently.

## Links

[Best Regexp Alternative For Go](https://itnext.io/best-regexp-alternative-for-go-be42abdc1fbb)

Discusses the performance and capabilities of various regex engines for Go. 

[Regular Expression Matching in the Wild](https://swtch.com/~rsc/regexp/regexp3.html)

Part 3 of a famous series of articles by Russ Cox detailing the implementation
of RE2, the engine behind Go's stdlib regex package. 
