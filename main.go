package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	parallel = flag.Int("parallel", 10, "number of parallel requests")
)

type siteSum struct {
	site string
	sum  string
	err  error
}

type hashFunc func(io.Reader) (string, error)

type App struct {
	sites    []string
	parallel int
	hash     hashFunc
	client   *http.Client
}

func NewApp(rawLinks []string, parallel int, hash hashFunc, timeout time.Duration) *App {
	var sites []string

	for _, s := range rawLinks {
		if link, err := url.Parse(s); err == nil {
			if link.Scheme == "" {
				link.Scheme = "http"
			}

			sites = append(sites, link.String())
		}
	}

	return &App{
		sites:    sites,
		parallel: parallel,
		hash:     hash,
		client: &http.Client{
			Timeout: timeout, // simple timeout
		},
	}
}

func (a *App) Run() ([]string, []error) {
	in, out := make(chan string), make(chan siteSum)

	var (
		sums []string
		errs []error
		wg   sync.WaitGroup
	)

	wg.Add(len(a.sites))

	go func() {
		for h := range out {
			if h.err != nil {
				errs = append(errs, h.err)
			} else {
				sums = append(sums, fmt.Sprintf("%v %v", h.site, h.sum))
			}

			wg.Done()
		}
	}()

	for i := 0; i < a.parallel; i++ {
		go a.fetchAndHash(in, out)
	}

	for _, s := range a.sites {
		in <- s
	}

	wg.Wait()

	close(in)
	close(out)

	return sums, errs
}

func (a *App) fetchAndHash(in chan string, out chan siteSum) {
	for s := range in {
		resp, err := a.client.Get(s)
		if err != nil {
			out <- siteSum{
				site: s,
				err:  err,
			}

			continue
		}
		sum, err := a.hash(resp.Body)
		if err != nil {
			out <- siteSum{
				site: s,
				err:  err,
			}

			continue
		}

		out <- siteSum{s, sum, err}
	}
}

var md5Hash = hashFunc(func(reader io.Reader) (string, error) {
	h := md5.New()
	_, err := io.Copy(h, reader)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
})

func main() {
	flag.Parse()

	app := NewApp(flag.Args(), *parallel, md5Hash, time.Minute)
	sums, _ := app.Run()

	for _, s := range sums {
		fmt.Println(s)
	}
}
