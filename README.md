# Nuxeo Go Client 1.0.2
This library is a Go HTTP client for Nuxeo Platform REST APIs.

This is compatible with All Nuxeo servers.

https://github1s.com/vpasquier/nuxeo-go-client

## Getting Started

### Server

- [Download a Nuxeo server](http://www.nuxeo.com/en/downloads) (the zip version)

- Unzip it

- Linux/Mac:
    - `NUXEO_HOME/bin/nuxeoctl start`
- Windows:
    - `NUXEO_HOME\bin\nuxeoctl.bat start`

- From your browser, go to `http://localhost:8080/nuxeo`

- Follow Nuxeo Wizard by clicking 'Next' buttons, re-start once completed

- Check Nuxeo correctly re-started `http://localhost:8080/nuxeo`
  - username: Administrator
  - password: Administrator

### Library Import

#### Import Nuxeo Go Client with:

```
require github.com/vpasquier/nuxeo-go-client
```

```
import "github.com/vpasquier/nuxeo-go-client"
```

### Usage

#### Creating a Client - Authentication

- Basic Auth:

```go
nuxeoClient := NuxeoClient().URL("https://demo.nuxeo.com/nuxeo").Username("Administrator").Password("Administrator").Debug(false).Build()
currentUser, err := nuxeoClient.Login()
log.println(currentUser.Username)
```

- Token:

```go
nuxeoClient := NuxeoClient().URL("http://localhost:8080/nuxeo").Token("XXXX").Build()
currentUser, err := nuxeoClient.Login()
log.println(currentUser.Username)
```

#### Options

- Headers:

```go
var headers map[string]string
headers["key"] = "value"

var cookies []*http.Cookie
var cookie *http.Cookie
cookie := &http.Cookie{
      Name:"go-resty",
      Value:"This is cookie value",
      Path: "/",
      Domain: "sample.com",
      MaxAge: 36000,
      HttpOnly: true,
      Secure: false,
    }

nuxeoClient := NuxeoClient().URL("https://demo.nuxeo.com/nuxeo").Timeout(1).Headers(headers).Cookies(cookies).Username("Administrator").Password("Administrator").Build()
```

- Schemas/Enrichers (schemas are by default empty, should be set to "*" to get all of them)

```go
nuxeoClient := NuxeoClient().URL("http://localhost:8080/nuxeo").Username("Administrator").Password("Administrator").Debug(false).Schemas([]string{"dublincore", "common"}).Enrichers("acls", "preview").Build()                       
```

- Debug (log request/response information - by default `false`):

```go
nuxeoClient := NuxeoClient().URL("https://demo.nuxeo.com/nuxeo").Username("Administrator").Password("Administrator").Debug(true).Build()
```

- More traces on the http calls and others: 

set env var `NUXEO_LOG_LEVEL` to `debug` (by default `info`)

#### APIs

#### Repository API

```go
// Here the document structure
type document struct {
	EntityType  string                 `json:"entity-type"`
	UID         string                 `json:"uid"`
	Path        string                 `json:"path"`
	Type        string                 `json:"type"`
	Name        string                 `json:"name"`
	Properties  map[string]interface{} `json:"properties"`
	nuxeoClient nuxeoClient
}
```

```go
// Fetch the root document
rootDocument, err := nuxeoClient.FetchDocumentRoot()
```

```go
// Fetch document by path
domain, err := nuxeoClient.FetchDocumentByPath("/default-domain")
```

```go
// Create a document
properties := map[string]interface{}{
	"dc:title": "New Document",
}

newDocument := document{
	EntityType: "document",
	Type:       "Workspaces",
	Name:       "new_file_with_go",
	Properties: properties,
}

newDocument, err = nuxeoClient.CreateDocument(domain.Path, newDocument)
```

```go
// Update a document
newDocument.Properties["dc:title"] = "Document Updated"
updatedDocument, err := nuxeoClient.UpdateDocument(newDocument)
```

```java
// Delete a document
err = nuxeoClient.DeleteDocument(updatedDocument)
```

```go
// Here the page provider result structure
type recordSet struct {
	Documents        []document `json:"entries"`
	TotalSize        int        `json:"totalSize"`
	CurrentPageIndex int        `json:"currentPageIndex"`
	NumberOfPages    int        `json:"numberOfPages"`
}
```


```go
// Fetch children
documents := domain.FetchChildren()
```

```go
// Get Blob
blob := document.FetchBlob("file:content)
```

```go
// Query
resultSet, err := nuxeoClient.Query("SELECT * FROM Domain")
assert.Equal(1, len(resultSet.Documents))
```

```go
// Directory represents a Nuxeo directory
type directory struct {
	EntityType    string                 `json:"entity-type"`
	DirectoryName string                 `json:"directoryName"`
	ID            string                 `json:"id"`
	Properties    map[string]interface{} `json:"properties"`
}

// DirectorySet represents a Nuxeo directory set
type directorySet struct {
	Entries []directory `json:"entries"`
}
```

```go
// Directories
directorySet, err := nuxeoClient.GetDirectory("continent")
assert.Equal(7, len(directorySet.Entries))
```

```go
// Create entry in directory
properties := make(map[string]interface{})
properties["id"] = "go"
properties["obsolete"] = "0"
properties["ordering"] = "10"
properties["label"] = "Go"

newDir := directory{
	EntityType:    "directoryEntry",
	DirectoryName: "continent",
	Properties:    properties,
}

returnedDir, err := nuxeoClient.CreateDirectory("continent", newDir)
```

```go
// Delete entry in directory
errDelete := nuxeoClient.DeleteDirectory("continent", "go")
```

```go
// Users API
returnedUser, err := nuxeoClient.GetUser("Administrator")

assert.Nil(err)
assert.Contains(returnedUser.Properties["groups"], "administrators")

properties := make(map[string]interface{})
properties["firstName"] = "Go"
properties["lastName"] = "Lang"
properties["group"] = [...]string{"administrators"}
properties["company"] = "nuxeo"
properties["email"] = "go@nuxeo.com"
properties["username"] = "go"

newUser := user{
	Username:   "go",
	EntityType: "user",
	Properties: properties,
}

returnedUser, err = nuxeoClient.CreateUser(newUser)

assert.Nil(err)

err = nuxeoClient.DeleteUser("go")
```

```go
// Async call
c := make(chan document, 1)

go nuxeoClient.AsyncFetchDocumentByPath("/default-domain", c)

select {
case rootDocument := <-c:
	assert.Equal("/default-domain", rootDocument.Path)
case <-time.After(1 * time.Second):
	assert.Fail("Result should have been received already")
```

#### Blobs

```go
// Attach document
params["document"] = "/default-domain/workspaces/workspace/file"
params["save"] = "true"
params["xpath"] = "file:content"

image, _ := ioutil.ReadFile("pink.jpg")

blob, blobError := nuxeoClient.Automation().Operation("Blob.AttachOnDocument").Parameters(params).Blob("pink.jpg", image).BlobExecute()
```

```go
// Fetch blob
file, err := nuxeoClient.FetchDocumentByPath("/default-domain/workspaces/workspace/file")
blob, blobError := file.FetchBlob("file:content")
assert.Equal(1025580, len(blob))
```

```go
// Async call for downloading a blob
file, err := nuxeoClient.FetchDocumentByPath("/default-domain/workspaces/workspace/file")

c := make(chan []byte, 1)

go file.AsyncFetchBlob("file:content", c)

select {
case blob := <-c:
	assert.Equal(1025580, len(blob))
case <-time.After(10 * time.Second):
	assert.Fail("Result should have been received already")
}
```

#### Automation/Operation API

```go
// Fetch document
params := make(map[string]string)
params["value"] = "/"
doc, err := nuxeoClient.Automation().Operation("Repository.GetDocument").Parameters(params).DocExecute()
```

```go
// Query
params["query"] = "SELECT * FROM Document"
records, err := nuxeoClient.Automation().Operation("Repository.Query").Parameters(params).DocListExecute()
```

## Missing Stuff

- Batch Upload (easy to do with https://github.com/go-resty/resty#using-file-directly-from-path)
- Automation has not been implemented/tested for all inputs/outputs (easy to enrich)
- Certainly other little gaps...

## Reporting Issues

We are glad to welcome new developers on this initiative, and even simple usage feedback is great.

- Ask your questions on [Nuxeo Answers](http://answers.nuxeo.com)
- Report issues on this GitHub repository (see [issues link](http://github.com/vpasquier/nuxeo-go-client/issues) on the right)
- Contribute: Send pull requests!

## About Nuxeo

Nuxeo dramatically improves how content-based applications are built, managed and deployed, making customers more agile, innovative and successful. Nuxeo provides a next generation, enterprise ready platform for building traditional and cutting-edge content oriented applications. Combining a powerful application development environment with SaaS-based tools and a modular architecture, the Nuxeo Platform and Products provide clear business value to some of the most recognizable brands including Verizon, Electronic Arts, Sharp, FICO, the U.S. Navy, and Boeing. Nuxeo is headquartered in New York and Paris. More information is available at [www.nuxeo.com](http://www.nuxeo.com/).
