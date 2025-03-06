package test

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func hxClient() *http.Client {
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func hxRequest(t *testing.T, method string, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	assert.NoError(t, err)
	req.Header.Add("Hx-Request", "true")
	return req
}

func HxGet(t *testing.T, url string) (*http.Response, *goquery.Document) {
	req := hxRequest(t, http.MethodGet, url, nil)
	res, err := hxClient().Do(req)
	assert.NoError(t, err)
	return res, hxDoc(t, res)
}

func HxPost(t *testing.T, url string, data map[string]string) (*http.Response, *goquery.Document) {
	req := hxRequest(t, http.MethodPost, url, formDataReader(data))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := hxClient().Do(req)
	assert.NoError(t, err)
	doc := hxDoc(t, res)
	return res, doc
}

func hxDoc(t *testing.T, res *http.Response) *goquery.Document {
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	assert.NoError(t, err)
	return doc
}

func formDataReader(data map[string]string) io.Reader {
	values := url.Values{}
	for k, v := range data {
		values.Set(k, v)
	}
	return strings.NewReader(values.Encode())
}

func FindErrorMsg(doc *goquery.Document) string {
	return doc.Find("#error").Text()
}

func AssertOk(t *testing.T, res *http.Response) bool {
	return assert.Equal(t, http.StatusOK, res.StatusCode)
}

func AssertRedirect(t *testing.T, location string, res *http.Response) bool {
	return assert.Equal(t, http.StatusFound, res.StatusCode) &&
		assert.Equal(t, location, res.Header.Get("Location"))
}

func AssertNoErrorMsg(t *testing.T, doc *goquery.Document) bool {
	msg := FindErrorMsg(doc)
	return assert.Empty(t, msg, "%s", msg)
}

func assertCookie(t *testing.T, res *http.Response, name string) bool {
	cookies := res.Cookies()
	assert.Len(t, cookies, 1)
	cookie := cookies[0]
	assert.Equal(t, name, cookie.Name)
	return assert.NotEmpty(t, cookie.Value)
}

func AssertSession(t *testing.T, res *http.Response) bool {
	return assertCookie(t, res, "session_id")
}

func HxPushUrl(res *http.Response) string {
	return res.Header.Get("HX-Push-Url")
}
