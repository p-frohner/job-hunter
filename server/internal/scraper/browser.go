package scraper

import (
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

const UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"

// NewBrowser launches a single Chrome instance shared across all scrapers.
// Sharing one browser reduces memory usage and prevents the process from crashing
// under concurrent scraping load.
func NewBrowser(headless bool) (*rod.Browser, error) {
	l := launcher.New().
		Headless(headless).
		Set("disable-blink-features", "AutomationControlled").
		Set("user-agent", UserAgent).
		Set("no-sandbox", "").            // required: running as root in Docker
		Set("disable-dev-shm-usage", ""). // required: /dev/shm is 64MB in Docker
		Set("use-gl", "swiftshader").      // software GL: satisfies canvas fingerprinting without a real GPU
		Set("ignore-certificate-errors", "") // required: VPN/proxy SSL inspection

	if bin := os.Getenv("CHROME_BIN"); bin != "" {
		l = l.Bin(bin)
	}

	u, err := l.Launch()
	if err != nil {
		return nil, err
	}

	browser := rod.New().ControlURL(u)
	if err := browser.Connect(); err != nil {
		return nil, err
	}

	return browser, nil
}
