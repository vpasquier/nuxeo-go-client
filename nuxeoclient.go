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
	Create() (user, error)
	FetchDocumentRoot() (document, error)
	FetchDocumentByPath(path string) (document, error)
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

// Create the client after applying configuration
func (nuxeoClient *nuxeoClient) Create() (user, error) {

	url := nuxeoClient.url + "/api/v1/automation/login"

	resp, err := nuxeoClient.client.R().EnableTrace().Post(url)

	if err != nil {
		log.Printf("%v", err)
		return user{}, errors.New("Client cannot be created")
	}

	data := resp.Body()

	if !json.Valid(data) {
		log.Panicf("The response is not json validated")
		return user{}, errors.New("Unmarshalling issue with the current user response")
	}

	var currentUser user
	jsonErr := json.Unmarshal(resp.Body(), &currentUser)

	if jsonErr != nil {
		log.Panicf("Can't create Nuxeo Client cause %v", jsonErr)
		return user{}, errors.New("Unmarshalling issue with the current user response")
	}

	log.Println("Nuxeo Client Initialized")

	return currentUser, nil
}

func (nuxeoClient *nuxeoClient) FetchDocumentRoot() (document, error) {

	url := nuxeoClient.url + "/api/v1/path//"

	resp, err := nuxeoClient.client.R().EnableTrace().Get(url)

	if err != nil {
		log.Printf("%v", err)
		return document{}, errors.New("Error while fetching document")
	}

	data := resp.Body()

	if !json.Valid(data) {
		log.Panicf("The response is not json validated")
		return document{}, errors.New("Unmarshalling issue with the current document response")
	}

	var currentDoc document
	jsonErr := json.Unmarshal(resp.Body(), &currentDoc)

	if jsonErr != nil {
		log.Panicf("Can't create Nuxeo Client cause %v", jsonErr)
		return document{}, errors.New("Unmarshalling issue with the current user response")
	}

	currentDoc.nuxeoClient = *nuxeoClient

	return currentDoc, nil
}

func (nuxeoClient *nuxeoClient) FetchDocumentByPath(path string) (document, error) {

	url := nuxeoClient.url + "/api/v1/path" + path

	resp, err := nuxeoClient.client.R().EnableTrace().Get(url)

	if err != nil {
		log.Printf("%v", err)
		return document{}, errors.New("Error while fetching document")
	}

	data := resp.Body()

	if !json.Valid(data) {
		log.Panicf("The response is not json validated")
		return document{}, errors.New("Unmarshalling issue with the current document response")
	}

	var currentDoc document
	jsonErr := json.Unmarshal(resp.Body(), &currentDoc)

	if jsonErr != nil {
		log.Panicf("Can't create Nuxeo Client cause %v", jsonErr)
		return document{}, errors.New("Unmarshalling issue with the current user response")
	}

	currentDoc.nuxeoClient = *nuxeoClient

	return currentDoc, nil
}
