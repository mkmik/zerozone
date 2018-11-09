package zerozone

import (
	"context"
	"encoding/json"
	"net"
	"strings"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/mholt/caddy"
	"github.com/miekg/dns"
)

type ZeroZoneHandler struct {
	Domain string
	Shell  *shell.Shell
	Next   plugin.Handler
}

type Zone struct {
	Records []Record `json:"records"`
}

type Record struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	TTL     uint32   `json:"ttl"`
	RRDatas []string `json:"rrdatas"`
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

	cfg := dnsserver.GetConfig(c)
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return &ZeroZoneHandler{
			Domain: cfg.Zone,
			Shell:  shell.NewShell(ipfsNodeAddr),
			Next:   next,
		}
	})

	return nil
}

func (h *ZeroZoneHandler) Name() string { return "zerozone" }

func (h *ZeroZoneHandler) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r, Context: ctx}

	qname := state.Name()

	domain := plugin.Zones([]string{h.Domain}).Matches(qname)
	if domain == "" {
		return plugin.NextOrFailure(h.Name(), h.Next, ctx, w, r)
	}
	state.Zone = domain

	comp := dns.SplitDomainName(strings.TrimSuffix(qname, domain))

	// allow falling through
	if len(comp) == 0 {
		return plugin.NextOrFailure(h.Name(), h.Next, ctx, w, r)
	}

	hash := comp[len(comp)-1]
	rs, err := h.Shell.Cat(hash)
	if err != nil {
		return dns.RcodeServerFailure, err
	}
	defer rs.Close()

	var zone Zone
	if err := json.NewDecoder(rs).Decode(&zone); err != nil {
		return 0, err
	}
	key := strings.Join(comp[:len(comp)-1], ".")

	m := new(dns.Msg)
	m.SetReply(r)

	for _, rr := range zone.Records {
		hdr := dns.RR_Header{Name: qname, Rrtype: state.QType(), Class: dns.ClassINET, Ttl: rr.TTL}

		if key == rr.Name && state.Type() == rr.Type {
			for _, d := range rr.RRDatas {
				var ans dns.RR
				switch state.QType() {
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
				case dns.TypeTXT:
					ans = &dns.TXT{
						Hdr: hdr,
						Txt: split255(d),
					}
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
