package setup

import (
	"fmt"
	"strconv"

	"github.com/mholt/caddy/middleware"
	"github.com/mholt/caddy/middleware/status"
)

func Status(c *Controller) (middleware.Middleware, error) {
	configs, err := statusParse(c)
	if err != nil {
		return nil, err
	}

	status := status.Status{
		Configs: configs,
	}

	return func(next middleware.Handler) middleware.Handler {
		status.Next = next
		return status
	}, nil
}

func statusParse(c *Controller) ([]status.Config, error) {
	var configs []status.Config

	appendCfg := func(sc status.Config) error {
		for _, c := range configs {
			if c.PathScope == sc.PathScope {
				return fmt.Errorf("Duplicate browsing config for %s", c.PathScope)
			}
		}
		configs = append(configs, sc)
		return nil
	}

	for c.Next() {
		var sc status.Config

		// First argument is directory to allow browsing; default is site root
		if c.NextArg() {
			sc.PathScope = c.Val()
		} else {
			sc.PathScope = "/"
		}

		// Second argument would be the code
		if c.NextArg() {
			v := c.Val()
			if n, err := strconv.Atoi(v); err == nil {
				sc.Code = n
			} else {
				return nil, c.Errf("Invalid code '%s' specified for status middleware.", v)
			}
		} else {
			return nil, c.Errf("No code specified for status middleware.")
		}

		// Third argument would be the body
		if c.NextArg() {
			sc.Body = c.Val()
		} else {
			sc.Body = ""
		}

		// Save configuration
		if err := appendCfg(sc); err != nil {
			return configs, err
		}
	}

	return configs, nil
}
