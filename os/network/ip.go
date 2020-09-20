package network

import (
	"context"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/sjengpho/tin/tin"
)

var readAll = ioutil.ReadAll

// ErrTimeout means the IP lookup has timed out.
var ErrTimeout = errors.New("ip lookup: timed out")

// NewPublicIPLookup returns a tin.PublicIPLookup.
func NewPublicIPLookup() tin.PublicIPLookup {
	return &ipLookupper{
		timeout: 10 * time.Second,
		sources: []string{
			"https://ifconfig.me/ip",
			"https://ifconfig.co/ip",
			"https://api.ipify.org",
		},
	}
}

// ipLookupper implements tin.PublicIPLookup.
type ipLookupper struct {
	timeout time.Duration
	sources []string
}

// Lookup returns a tin.PublicIP.
func (i *ipLookupper) Lookup() (tin.PublicIP, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, i.timeout)
	defer cancel()

	// Requesting the IP address from the sources.
	var wg sync.WaitGroup
	ch := make(chan net.IP, 1)
	for _, s := range i.sources {
		wg.Add(1)
		go i.lookup(ctx, s, ch, &wg)
	}

	// Cleaning up channel.
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Cancelling the other requests when one of the sources
	// has returned a result. When the deadline has exceeded and none of the
	// sources has returned a result an error is returned.
	select {
	case i := <-ch:
		return tin.PublicIP{IP: i}, nil
	case <-ctx.Done():
		return tin.PublicIP{}, ErrTimeout
	}
}

// lookup fetches the IP address and sends the result to the channel.
func (i *ipLookupper) lookup(ctx context.Context, url string, ch chan net.IP, wg *sync.WaitGroup) {
	defer wg.Done()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return
	}

	defer resp.Body.Close()
	data, err := readAll(resp.Body)
	if err != nil {
		return
	}

	ip := net.ParseIP(string(data))
	if ip == nil {
		return
	}

	ch <- ip
}
