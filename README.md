Instructions:

1. clone this outside of the GOPATH (so go modules work)
2. cd server
3. `go build && ./server`

In another shell:

```
$ dig  @127.0.0.1 -p 8053 A foo.bafybeib3u2yfpkoticclpzwcrwkjzc4hlqwpttohhqlmf55qdgk3hrutcm.foo.mkm.pub +short
10.20.30.40
$ dig  @127.0.0.1 -p 8053 A foo.bafybeidz2eomuhekgmhwnoxawruyrbrn6yg23p72qswsa5kicegoydzq4q.foo.mkm.pub +short
4.3.2.1
```

Main TODOs:

1. switch to /ipns. this now hosts only immutable /ipfs hashes (for speed of iterating locally)
2. make sure plugin chaining works so we can serve static NS records alongside the dynamic plugin
3. build docker image and deploy somewhere public (I'm preparing 0zone.mkm.pub for now)
