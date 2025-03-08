package test

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

var httpClient *http.Client

func HxClient() *http.Client {
	if httpClient != nil {
		return httpClient
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	httpClient = &http.Client{Jar: jar}
	return httpClient
}

func HxRequest(t *testing.T, method string, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	assert.NoError(t, err)
	req.Header.Add("Hx-Request", "true")
	return req
}

func HxGet(t *testing.T, url string) (*http.Response, *goquery.Document) {
	req := HxRequest(t, http.MethodGet, url, nil)
	res, err := HxClient().Do(req)
	assert.NoError(t, err)
	return res, HxDoc(t, res)
}

func HxPost(t *testing.T, url string, data map[string]string) (*http.Response, *goquery.Document) {
	req := HxRequest(t, http.MethodPost, url, formDataReader(data))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := HxClient().Do(req)
	assert.NoError(t, err)
	doc := HxDoc(t, res)
	return res, doc
}

func HxDoc(t *testing.T, res *http.Response) *goquery.Document {
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

func AssertNoErrorMsg(t *testing.T, doc *goquery.Document) bool {
	msg := FindErrorMsg(doc)
	return assert.Empty(t, msg, "%s", msg)
}

func assertCookie(t *testing.T, res *http.Response, name string) bool {
	cookies := res.Cookies()
	assert.Len(t, cookies, 1, "No cookies")
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
