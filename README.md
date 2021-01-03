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

#### Creating a Client

```go
  nuxeoClient := NuxeoClient().URL("http://localhost:8080/nuxeo").Username("Administrator").Password("Administrator").Build()
  currentUser, err := nuxeoClient.Create()
  log.println(currentUser.Username)
```

### Authentication

Basic:

```go
  nuxeoClient := NuxeoClient().URL("http://localhost:8080/nuxeo").Username("Administrator").Password("Administrator").Build()
  currentUser, err := nuxeoClient.Create()
  log.println(currentUser.Username)
```

Token:

```go
  nuxeoClient := NuxeoClient().URL("http://localhost:8080/nuxeo").Token("XXXX").Build()
  currentUser, err := nuxeoClient.Create()
  log.println(currentUser.Username)
```

#### Options

```go
var headers map[string]string
headers["content-type"] = "application/json"

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
nuxeoClient := NuxeoClient().URL("http://localhost:8080/nuxeo").Timeout(10).Headers(headers).Cookies().Build()
```

```go
nuxeoClient := NuxeoClient().Schemas("dublincore", "common")
                         .Enrichers("acls", "preview")
```

```go
// To fetch all schemas
nuxeoClient := NuxeoClient().Schemas("*")
```

```go
// Log out
NuxeoClient().destroy()
```

#### APIs

General rule:

- When using `fetch` methods, `NuxeoClient` is making remote calls.
- When using `get` methods, objects are retrieved from memory.

#### Operation API

```go
// Fetch the root document
document := nuxeoClient.operation(Operations.REPOSITORY_GET_DOCUMENT).param("value", "/").execute();
```

```go
// Fetch the root document
documents := nuxeoClient.operation("Repository.Query")
                            .param("query", "SELECT * FROM Document")
                            .execute();
```

```go
// with blob
fileBlob := new FileBlob(File file);
nuxeoClient.operation(Operations.BLOB_ATTACH_ON_DOCUMENT)
           .voidOperation(true)
           .param("document", "/folder/file")
           .input(fileBlob)
           .execute();

// or with stream
streamBlob := new StreamBlob(InputStream stream, String filename);
nuxeoClient.operation(Operations.BLOB_ATTACH_ON_DOCUMENT)
           .voidOperation(true)
           .param("document", "/folder/file")
           .input(streamBlob)
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

// you need to close the stream or to get the file
blob := nuxeoClient.operation(Operations.DOCUMENT_GET_BLOB)
                       .input("folder/file")
                       .execute();
```

#### Repository API

```go
// Fetch the root document
rootDocument := nuxeoClient.repository().fetchDocumentRoot();
```

```go
// Fetch document by path
folder := nuxeoClient.repository().fetchDocumentByPath("/folder_2");
```

```go
// Create a document
document := Document.createWithName("file", "File");
document.setPropertyValue("dc:title", "new title");
result := nuxeoClient.repository().createDocumentByPath("/folder_1", document);
```

```go
// Update a document
document := nuxeoClient.repository().fetchDocumentByPath("/folder_1/note_0");
documentUpdated := Document.createWithId(document.getId(), "Note");
documentUpdated.setPropertyValue("dc:title", "note updated");
documentUpdated.setPropertyValue("dc:nature", "test");
documentUpdated := nuxeoClient.repository().updateDocument(documentUpdated);
```

```java
// Delete a document
Document documentToDelete = nuxeoClient.repository().fetchDocumentByPath("/folder_1/note_1");
nuxeoClient.repository().deleteDocument(documentToDelete);
```

```go
// Fetch children
folder := nuxeoClient.repository().fetchDocumentByPath("/folder_2");
children := folder.fetchChildren();
```

```go
// Fetch blob
file := nuxeoClient.repository().fetchDocumentByPath("/folder_2/file");
blob := file.fetchBlob();
```

```go
// Execute query
documents := nuxeoClient.repository().query("SELECT * From Note");
```

```go
// Fetch document asynchronously with callback
nuxeoClient.repository().fetchDocumentRoot(new Callback<Document>() {
            @Override
            public void onResponse(Call<Document> call, Response<Document>
                    response) {
                if (!response.isSuccessful()) {
                    ObjectMapper objectMapper = new ObjectMapper();
                    NuxeoClientException nuxeoClientException;
                    try {
                        nuxeoClientException = objectMapper.readValue(response.errorBody().string(),
                                NuxeoClientException.class);
                    } catch (IOException reason) {
                        throw new NuxeoClientException(reason);
                    }
                    fail(nuxeoClientException.getRemoteStackTrace());
                }
                Document folder = response.body();
                assertNotNull(folder);
                assertEquals("Folder", folder.getType());
                assertEquals("document", folder.getEntityType());
                assertEquals("/folder_2", folder.getPath());
                assertEquals("Folder 2", folder.getTitle());
            }

            @Override
            public void onFailure(Call<Document> call, Throwable t) {
                fail(t.getMessage());
            }
        });
```

#### Batch Upload

Batch uploads are executed through the `org.nuxeo.client.objects.upload.BatchUploadManager`.

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