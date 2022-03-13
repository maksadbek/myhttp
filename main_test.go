package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestParseURL(t *testing.T) {
	rawLinks := []string{
		"abc",
		"abc.com",
		"goo.gle",
		"www.google.com",
		"example.com/a/b/c/d/e",
	}

	app := NewApp(rawLinks, 10, nil, time.Second)

	expectedSites := []string{
		"http://abc.com",
		"http://goo.gle",
		"http://www.google.com",
		"http://example.com/a/b/c/d/e",
	}

	assertEqualSlice(t, expectedSites, app.sites)
}

func TestAppRunSuccess(t *testing.T) {
	var (
		expectedErrs []error
		expectedSums []string
	)

	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "Hello")
	}))
	defer server1.Close()

	expectedSums = append(expectedSums, server1.URL+" 8b1a9953c4611296a827abf8c47804d7") // md5 of "Hello"

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "Hello world")
	}))
	defer server2.Close()

	expectedSums = append(expectedSums, server2.URL+" 3e25960a79dbc69b674cd4ec67a72c62") // md5 of "Hello world"

	app := NewApp([]string{server1.URL, server2.URL}, 10, md5Hash, time.Minute)
	sums, errs := app.Run()

	assertEqualErrors(t, expectedErrs, errs)
	assertEqualSlice(t, expectedSums, sums)
}

func TestAppRunServerTimeout(t *testing.T) {
	var (
		expectedErrs []error
		expectedSums []string
	)

	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second)
		_, _ = fmt.Fprint(w, "Hello")
	}))
	defer server1.Close()

	expectedErrs = append(expectedErrs, fmt.Errorf(
		"Get \"%v\": context deadline exceeded (Client.Timeout exceeded while awaiting headers)",
		server1.URL))

	app := NewApp([]string{server1.URL, server1.URL}, 10, md5Hash, time.Millisecond)
	sums, errs := app.Run()

	assertEqualErrors(t, expectedErrs, errs)
	assertEqualSlice(t, expectedSums, sums)
}

func TestAppRunHashError(t *testing.T) {
	var (
		expectedErrs []error
		expectedSums []string
	)

	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "Hello")
	}))
	defer server1.Close()

	expectedErrs = append(expectedErrs, errors.New("surprise"))

	erronousHash := func(r io.Reader) (string, error) {
		return "", errors.New("surprise")
	}

	app := NewApp([]string{server1.URL}, 10, erronousHash, time.Minute)
	sums, errs := app.Run()

	assertEqualErrors(t, expectedErrs, errs)
	assertEqualSlice(t, expectedSums, sums)
}

func assertEqualSlice(t *testing.T, want, got []string) {
	set := map[string]struct{}{}
	for _, g := range got {
		set[g] = struct{}{}
	}

	for _, w := range want {
		if _, ok := set[w]; !ok {
			t.Errorf("expected %v, got nil", w)
		}
	}
}

func assertEqualErrors(t *testing.T, want, got []error) {
	set := map[string]struct{}{}

	for _, g := range got {
		set[g.Error()] = struct{}{}
	}

	for _, w := range want {
		if _, ok := set[w.Error()]; !ok {
			t.Errorf("expected %v, got nil", w)
		}
	}
}
