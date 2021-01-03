package nuxeoclient

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	// DefaultURL is the client url if none has been set
	DefaultURL = "http://localhost:8080/nuxeo"
)

// Client interface
type Client interface {
	Create() (User, error)
}

// ClientBuilder interface
type ClientBuilder interface {
	URL(string) ClientBuilder
	Username(string) ClientBuilder
	Password(string) ClientBuilder
	Token(string) ClientBuilder
	Debug(bool) ClientBuilder
	Timeout(int) ClientBuilder
	Headers(map[string]string) ClientBuilder
	Cookies([]*http.Cookie) ClientBuilder
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
	headers     map[string]string
	cookies     []*http.Cookie
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
	headers     map[string]string
	cookies     []*http.Cookie
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

func (cb *clientBuilder) Timeout(timeout int) ClientBuilder {
	cb.timeout = timeout
	return cb
}

func (cb *clientBuilder) Headers(headers map[string]string) ClientBuilder {
	cb.headers = headers
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
	return &nuxeoClient{
		url:      cb.url,
		username: cb.username,
		password: cb.password,
		debug:    cb.debug,
		timeout:  cb.timeout,
		headers:  cb.headers,
		cookies:  cb.cookies,
	}
}

// NuxeoClient is the Nuxeo client builder
func NuxeoClient() ClientBuilder {
	return &clientBuilder{}
}

// Create the client after applying configuration
func (nuxeoClient *nuxeoClient) Create() (User, error) {

	log.Println("Creating Nuxeo Client...")

	client := resty.New()

	client.SetCookies(nuxeoClient.cookies)
	client.SetHeaders(nuxeoClient.headers)
	client.SetDebug(nuxeoClient.debug)
	client.SetTimeout(time.Duration(nuxeoClient.timeout))

	if nuxeoClient.token == "" {
		client.SetBasicAuth(nuxeoClient.username, nuxeoClient.password)
	} else {
		client.SetAuthToken(nuxeoClient.token)
	}

	url := nuxeoClient.url
	if url == "" {
		url = DefaultURL
	}

	url += "/api/v1/automation/login"

	resp, err := client.R().EnableTrace().Post(url)

	if err != nil {
		log.Panic(err)
		return User{}, errors.New("Client cannot be created")
	}

	data := resp.Body()

	if !json.Valid(data) {
		log.Panicf("The response is not json validated")
		return User{}, errors.New("Unmarshalling issue with the current user response")
	}

	var currentUser User
	jsonErr := json.Unmarshal(resp.Body(), &currentUser)

	if jsonErr != nil {
		log.Panicf("Can't create Nuxeo Client cause %d", jsonErr)
		return User{}, errors.New("Unmarshalling issue with the current user response")
	}

	return currentUser, nil
}
