package nuxeoclient

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-resty/resty/v2"
)

// ClientBuilder interface
type ClientBuilder interface {
	URL(string) ClientBuilder
	Username(string) ClientBuilder
	Password(string) ClientBuilder
	Token(string) ClientBuilder
	Debug(bool) ClientBuilder
	Timeout(int) ClientBuilder
	Schemas([]string) ClientBuilder
	Enrichers([]string) ClientBuilder
	Headers(map[string]string) ClientBuilder
	Cookies([]*http.Cookie) ClientBuilder
	Repository(string) ClientBuilder
	Build() Client
}

// Mutable
type clientBuilder struct {
	url         string
	username    string
	password    string
	token       string
	debug       bool
	enableTrace bool
	timeout     int
	schemas     []string
	enrichers   []string
	headers     map[string]string
	cookies     []*http.Cookie
	repository  string
}

// Immutable
type nuxeoClient struct {
	url         string
	username    string
	password    string
	token       string
	debug       bool
	enableTrace bool
	timeout     int
	schemas     []string
	enrichers   []string
	headers     map[string]string
	cookies     []*http.Cookie
	repository  string
	client      *resty.Client
}

func (cb *clientBuilder) URL(url string) ClientBuilder {
	cb.url = url
	return cb
}

func (cb *clientBuilder) Username(username string) ClientBuilder {
	cb.username = username
	return cb
}

func (cb *clientBuilder) Password(password string) ClientBuilder {
	cb.password = password
	return cb
}

func (cb *clientBuilder) Token(token string) ClientBuilder {
	cb.token = token
	return cb
}

func (cb *clientBuilder) Repository(repository string) ClientBuilder {
	cb.repository = repository
	return cb
}

func (cb *clientBuilder) Timeout(timeout int) ClientBuilder {
	cb.timeout = timeout
	return cb
}

func (cb *clientBuilder) Headers(headers map[string]string) ClientBuilder {
	cb.headers = headers
	return cb
}

func (cb *clientBuilder) Schemas(schemas []string) ClientBuilder {

	cb.schemas = schemas
	return cb
}

func (cb *clientBuilder) Enrichers(enrichers []string) ClientBuilder {
	cb.enrichers = enrichers
	return cb
}

func (cb *clientBuilder) Cookies(cookies []*http.Cookie) ClientBuilder {
	cb.cookies = cookies
	return cb
}

func (cb *clientBuilder) Debug(debug bool) ClientBuilder {
	cb.debug = debug
	return cb
}

func (cb *clientBuilder) Build() Client {

	log.Println("Creating Nuxeo Client...")

	client := resty.New()

	if val, ok := cb.headers["Content-Type"]; !ok {
		log.Debug(val)
		if cb.headers == nil {
			cb.headers = make(map[string]string)
		}
		cb.headers["Content-Type"] = "application/json"
	}

	for i := 0; i < len(cb.schemas); i++ {
		cb.headers["properties"] = cb.schemas[i]
	}

	for i := 0; i < len(cb.enrichers); i++ {
		cb.headers["enrichers.document"] = cb.enrichers[i]
	}

	if cb.repository != "" {
		cb.headers["Repository"] = cb.repository
	}

	client.SetCookies(cb.cookies)
	client.SetHeaders(cb.headers)
	client.SetDebug(cb.debug)
	client.SetTimeout(time.Duration(cb.timeout) * time.Minute)

	if cb.token == "" {
		client.SetBasicAuth(cb.username, cb.password)
	} else {
		client.SetAuthToken(cb.token)
	}

	url := cb.url
	if url == "" {
		url = DefaultURL
	}
	cb.url = url

	return &nuxeoClient{
		url:        cb.url,
		username:   cb.username,
		password:   cb.password,
		debug:      cb.debug,
		timeout:    cb.timeout,
		headers:    cb.headers,
		cookies:    cb.cookies,
		repository: cb.repository,
		client:     client,
	}
}

// NuxeoClient is the Nuxeo client builder
func NuxeoClient() ClientBuilder {
	return &clientBuilder{}
}
