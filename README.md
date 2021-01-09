# Nuxeo Go Client

This library is a Go HTTP client for Nuxeo Platform REST APIs.

This is compatible with All Nuxeo servers.

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
require github.com/vpasquier/nuxeo-go-client v1.0.0
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
// Fetch children
documents := domain.FetchChildren()
```

```go
// Get Blob
blob := document.FetchBlob()
blob := document.FetchCustomBlob("customfile:content")
```

For information:

```go
type Blob struct {
  filename string
  size int
  file File?
}
```

```go
// Query
resultSet, err := nuxeoClient.Query("SELECT * FROM Domain")
assert.Equal(len(resultSet.Documents), 1)
```

for information:

```go
type RecordSet struct {
	Documents        []document `json:"entries"`
	TotalSize        int        `json:"totalSize"`
	CurrentPageIndex int        `json:"currentPageIndex"`
	NumberOfPages    int        `json:"numberOfPages"`
}
```

```go
// Async call for downloading a blob

```

```go
// Directories
resultSet, err := nuxeoClient.Directory("continent")
```

For information:

```go
// Directory represents a Nuxeo directory
type Directory struct {
	directoryName string                 `json:"directoryName"`
	id            string                 `json:"id"`
	properties    map[string]interface{} `json:"properties"`
}

// DirectorySet represents a Nuxeo directory set
type DirectorySet struct {
	entries []Directory `json:"entries"`
}
```

#### Operation API

```go
// Fetch document
document := nuxeoClient.operation(Operations.REPOSITORY_GET_DOCUMENT).param("value", "/").execute();
```

```go
// Query
documents := nuxeoClient.operation("Repository.Query")
                            .param("query", "SELECT * FROM Document")
                            .execute();
```

```go
// Attach blobs
fileBlob := new FileBlob(File file);
nuxeoClient.operation(Operations.BLOB_ATTACH_ON_DOCUMENT)
           .voidOperation(true)
           .param("document", "/folder/file")
           .input(fileBlob)
           .execute();

inputBlobs := new []Blobs();
inputBlobs.add(File file1);
inputBlobs.add(new StreamBlob(InputStream stream, String filename2));
Blobs blobs = nuxeoClient.operation(Operations.BLOB_ATTACH_ON_DOCUMENT)
                         .voidOperation(true)
                         .param("xpath", "files:files")
                         .param("document", "/folder/file")
                         .input(inputBlobs)
                         .execute();
```

```go
// Fetch blob
file := nuxeoClient.repository().fetchDocumentByPath("/folder_2/file");
blob := file.fetchBlob();
```

#### Batch Upload

```java
// Batch Upload Manager
BatchUploadManager batchUploadManager = nuxeoClient.uploadManager();
BatchUpload batchUpload = batchUploadManager.createBatch();
```

```java
// Upload File
File file = FileUtils.getResourceFileFromContext("sample.jpg");
batchUpload = batchUpload.upload("1", file);

// Fetch/Refresh the batch file information from server
batchUpload = batchUpload.fetchBatchUpload("1");

// Directly from the manager
batchUpload = batchUpload.fetchBatchUpload(batchUpload.getBatchId(), "1");

// Upload another file and check files
file = FileUtils.getResourceFileFromContext("blob.json");
batchUpload.upload("2", file);
List<BatchUpload> batchFiles = batchUpload.fetchBatchUploads();
```
Batch upload can be executed in a [chunk mode](https://doc.nuxeo.com/display/NXDOC/Blob+Upload+for+Batch+Processing?src=search#BlobUploadforBatchProcessing-UploadingaFilebyChunksUploadingaFilebyChunks).

```java
// Upload file chunks
BatchUploadManager batchUploadManager = nuxeoClient.uploadManager();
BatchUpload batchUpload = batchUploadManager.createBatch();
batchUpload.enableChunk();
File file = FileUtils.getResourceFileFromContext("sample.jpg");
batchUpload = batchUpload.upload("1", file);
```

Chunk size is by default 1MB (int 1024*1024). You can update this value with:

```java
batchUpload.chunkSize(1024);
```

Attach batch to a document:

```java
Document doc = new Document("file", "File");
doc.set("dc:title", "new title");
doc = nuxeoClient.repository().createDocumentByPath("/folder_1", doc);
doc.set("file:content", batchUpload.getBatchBlob());
doc = doc.updateDocument();
```

or with operation:

```java
Document doc = new Document("file", "File");
doc.set("dc:title", "new title");
doc = nuxeoClient.repository().createDocumentByPath("/folder_1", doc);
Blob blob = batchUpload.operation(Operations.BLOB_ATTACH_ON_DOCUMENT).param("document", doc).execute();
```

## Reporting Issues

We are glad to welcome new developers on this initiative, and even simple usage feedback is great.

- Ask your questions on [Nuxeo Answers](http://answers.nuxeo.com)
- Report issues on this GitHub repository (see [issues link](http://github.com/vpasquier/nuxeo-go-client/issues) on the right)
- Contribute: Send pull requests!

## About Nuxeo

Nuxeo dramatically improves how content-based applications are built, managed and deployed, making customers more agile, innovative and successful. Nuxeo provides a next generation, enterprise ready platform for building traditional and cutting-edge content oriented applications. Combining a powerful application development environment with SaaS-based tools and a modular architecture, the Nuxeo Platform and Products provide clear business value to some of the most recognizable brands including Verizon, Electronic Arts, Sharp, FICO, the U.S. Navy, and Boeing. Nuxeo is headquartered in New York and Paris. More information is available at [www.nuxeo.com](http://www.nuxeo.com/).