package zerozone

import (
	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("zerozone", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	c.Next() // Ignore "zerozone" and give us the next token.
	return nil
}
