# Zero Zone

> This project is still in preparation. I'm moving this from 0zone.mkm.pub to a non-personal domain, and switching to IPNS

Zero Zone is a Zero Conf public domain registrar. 

With Zero Zone, anybody can create a zone. Zero Zone is automation friendly; no registrations, no captchas.
A zone is just a JSON file available via IPFS.

```
$ dig A foo.bafybeib3u2yfpkoticclpzwcrwkjzc4hlqwpttohhqlmf55qdgk3hrutcm.0zone.mkm.pub +short
10.20.30.40
```

Zero zones are not good-looking. If you want a better looking name, or one you can choose yourself this tool is not for you.
Whether zero zones are temporary is up to you: as long as the IPFS file gets served the zone is alive.

The zone is then mapped as subdomain of 0zone.mkm.pub, using the base32 encoding
of the IPFS file address.

The JSON zone format reuses Google's Cloud DNS [format](https://cloud.google.com/dns/records/json-record),
for no other reason except that the authors didn't yet find a better standard (suggestions accepted) and this was easy to parse.

```
$ ipfs cat /ipfs/bafybeib3u2yfpkoticclpzwcrwkjzc4hlqwpttohhqlmf55qdgk3hrutcm
{
   "records" : [
      {
         "rrdatas" : [
            "10.20.30.40"
         ],
         "type" : "A",
         "ttl" : 600,
         "name" : "foo"
      },
      {
         "type" : "TXT",
         "rrdatas" : [
            "yadda yadda"
         ],
         "name" : "_acme-challenge.foo",
         "ttl" : 600
      }
   ]
}

```

## Background

When automating deployments that require TLS certificates, you often also need to automate DNS,
which in turn requires the automation to have API access to a DNS service. This raises the barrier to entry
for those who just want deploy something anywhere, without having to have access to a cloud DNS service.

An example use case for this tool is streamlining installation procedures for [Bitnami Kubernetes Production Runtime](https://github.com/bitnami/kube-prod-runtime/tree/master/kubeprod)

## Install

These instructions assume a IPFS node running on localhost.

1. clone this outside of the GOPATH (so go modules work)
2. `cd cmd/server`
3. `go build && ./server`

You can try it out in another shell:

```
$ dig  @127.0.0.1 -p 8053 A foo.bafybeib3u2yfpkoticclpzwcrwkjzc4hlqwpttohhqlmf55qdgk3hrutcm.0zone.mkm.pub +short
10.20.30.40
$ dig  @127.0.0.1 -p 8053 A foo.bafybeidz2eomuhekgmhwnoxawruyrbrn6yg23p72qswsa5kicegoydzq4q.0zone.mkm.pub +short
4.3.2.1
```

## Contributing

PRs accepted.

Main TODOs:

1. switch to /ipns. this now hosts only immutable /ipfs hashes (for speed of iterating locally)
2. make sure plugin chaining works so we can serve static NS records alongside the dynamic plugin
3. build docker image and deploy somewhere public (I'm preparing 0zone.mkm.pub for now)
