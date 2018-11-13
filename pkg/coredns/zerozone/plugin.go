// package zerozone is a coredns plugin that serves Zero Zones.
package zerozone

import (
	"context"
	"net"
	"strings"

	"github.com/bitnami-labs/zerozone/pkg/store"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/mholt/caddy"
	"github.com/miekg/dns"
)

// ZeroZoneHandler implemens that coredns plugin handler.
type ZeroZoneHandler struct {
	Domain  string
	Fetcher store.Fetcher
	Next    plugin.Handler
}

func init() {
	caddy.RegisterPlugin("zerozone", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	c.Next()
	var ipfsNodeAddr string
	if !c.Args(&ipfsNodeAddr) {
		return plugin.Error("zerozone", c.ArgErr())
	}

	var fetcher store.Fetcher
	if strings.HasPrefix(ipfsNodeAddr, "http") {
		fetcher = store.NewIPNSGatewayFetcher(ipfsNodeAddr)
	} else {
		fetcher = store.NewIPNSFetcher(ipfsNodeAddr)
	}

	cfg := dnsserver.GetConfig(c)
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return &ZeroZoneHandler{
			Domain:  cfg.Zone,
			Fetcher: store.NewSingleFlightFetcher(store.NewCachingFetcher(fetcher)),
			Next:    next,
		}
	})

	return nil
}

func (h *ZeroZoneHandler) Name() string { return "zerozone" }

func parseQuery(qname, domain string) (hostname, zoneID string, ok bool) {
	comp := dns.SplitDomainName(strings.TrimSuffix(qname, domain))
	if len(comp) == 0 {
		return "", "", false
	}
	return strings.Join(comp[:len(comp)-1], "."), comp[len(comp)-1], true
}

func (h *ZeroZoneHandler) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r, Context: ctx}

	qname := state.Name()

	domain := plugin.Zones([]string{h.Domain}).Matches(qname)
	if domain == "" {
		return plugin.NextOrFailure(h.Name(), h.Next, ctx, w, r)
	}
	state.Zone = domain

	hostname, zoneID, ok := parseQuery(qname, domain)

	// allow falling through
	if !ok {
		return plugin.NextOrFailure(h.Name(), h.Next, ctx, w, r)
	}

	zone, err := h.Fetcher.FetchZone(zoneID)
	if err != nil {
		return dns.RcodeServerFailure, plugin.Error("zerozone", err)
	}

	m := new(dns.Msg)
	m.SetReply(r)

	for _, rr := range zone.Records {
		hdr := dns.RR_Header{Name: qname, Rrtype: state.QType(), Class: dns.ClassINET, Ttl: rr.TTL}

		if hostname == rr.Name && state.Type() == rr.Type {
			for _, d := range rr.RRDatas {
				var ans dns.RR
				switch t := state.QType(); t {
				case dns.TypeA:
					ans = &dns.A{
						Hdr: hdr,
						A:   net.ParseIP(d),
					}
				case dns.TypeAAAA:
					ans = &dns.AAAA{
						Hdr:  hdr,
						AAAA: net.ParseIP(d),
					}
				case dns.TypeCNAME:
					ans = &dns.CNAME{
						Hdr:    hdr,
						Target: d,
					}
				case dns.TypeTXT:
					ans = &dns.TXT{
						Hdr: hdr,
						Txt: split255(d),
					}
				default:
					log.Debugf("unhandled type %q", t)
				}
				m.Answer = append(m.Answer, ans)
			}

			break
		}
	}

	state.SizeAndDo(m)
	m = state.Scrub(m)
	w.WriteMsg(m)
	return dns.RcodeSuccess, nil
}

func split255(s string) []string {
	if len(s) < 255 {
		return []string{s}
	}
	sx := []string{}
	p, i := 0, 255
	for {
		if i <= len(s) {
			sx = append(sx, s[p:i])
		} else {
			sx = append(sx, s[p:])
			break

		}
		p, i = p+255, i+255
	}

	return sx
}
