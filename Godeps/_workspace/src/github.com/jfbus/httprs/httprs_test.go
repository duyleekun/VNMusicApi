package httprs

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type fakeResponseWriter struct {
	code int
	h    http.Header
	tmp  *os.File
}

func (f *fakeResponseWriter) Header() http.Header {
	return f.h
}

func (f *fakeResponseWriter) Write(b []byte) (int, error) {
	return f.tmp.Write(b)
}

func (f *fakeResponseWriter) Close(b []byte) error {
	return f.tmp.Close()
}

func (f *fakeResponseWriter) WriteHeader(code int) {
	f.code = code
}

func (f *fakeResponseWriter) Response() *http.Response {
	f.tmp.Seek(0, os.SEEK_SET)
	return &http.Response{Body: f.tmp, StatusCode: f.code, Header: f.h}
}

type fakeRoundTripper struct {
	src *os.File
}

func (f *fakeRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	fw := &fakeResponseWriter{h: http.Header{}}
	var err error
	fw.tmp, err = ioutil.TempFile(os.TempDir(), "httprs")
	if err != nil {
		return nil, err
	}
	http.ServeContent(fw, r, "temp.txt", time.Now(), f.src)

	return fw.Response(), nil
}

const SZ = 1024

func newRS() *HttpReadSeeker {
	tmp, err := ioutil.TempFile(os.TempDir(), "httprs")
	if err != nil {
		return nil
	}
	for i := 0; i < SZ; i++ {
		tmp.WriteString(fmt.Sprintf("%04d", i))
	}

	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		return nil
	}
	res := &http.Response{
		Header: http.Header{
			"Accept-Ranges": []string{"bytes"},
		},
		Request:       req,
		ContentLength: SZ * 4,
	}
	return NewHttpReadSeeker(res, &http.Client{Transport: &fakeRoundTripper{tmp}})
}

func TestHttpReaderSeeker(t *testing.T) {
	Convey("Scenario: testing HttpReaderSeeker", t, func() {

		Convey("Read should start at the beginning", func() {
			r := newRS()
			So(r, ShouldNotBeNil)
			defer r.Close()
			buf := make([]byte, 4)
			n, err := io.ReadFull(r, buf)
			So(n, ShouldEqual, 4)
			So(err, ShouldBeNil)
			So(string(buf), ShouldEqual, "0000")
		})

		Convey("Seek w SEEK_SET should seek to right offset", func() {
			r := newRS()
			So(r, ShouldNotBeNil)
			defer r.Close()
			s, err := r.Seek(4*64, os.SEEK_SET)
			So(s, ShouldEqual, 4*64)
			So(err, ShouldBeNil)
			buf := make([]byte, 4)
			n, err := io.ReadFull(r, buf)
			So(n, ShouldEqual, 4)
			So(err, ShouldBeNil)
			So(string(buf), ShouldEqual, "0064")
		})

		Convey("Read + Seek w SEEK_CUR should seek to right offset", func() {
			r := newRS()
			So(r, ShouldNotBeNil)
			defer r.Close()
			buf := make([]byte, 4)
			io.ReadFull(r, buf)
			s, err := r.Seek(4*64, os.SEEK_CUR)
			So(s, ShouldEqual, 4*64+4)
			So(err, ShouldBeNil)
			n, err := io.ReadFull(r, buf)
			So(n, ShouldEqual, 4)
			So(err, ShouldBeNil)
			So(string(buf), ShouldEqual, "0065")
		})

		Convey("Seek w SEEK_END should seek to right offset", func() {
			r := newRS()
			So(r, ShouldNotBeNil)
			defer r.Close()
			buf := make([]byte, 4)
			io.ReadFull(r, buf)
			s, err := r.Seek(4, os.SEEK_END)
			So(s, ShouldEqual, SZ*4-4)
			So(err, ShouldBeNil)
			n, err := io.ReadFull(r, buf)
			So(n, ShouldEqual, 4)
			So(err, ShouldBeNil)
			So(string(buf), ShouldEqual, fmt.Sprintf("%04d", SZ-1))
		})
	})
}
