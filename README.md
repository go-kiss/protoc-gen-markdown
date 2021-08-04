# protoc-gen-markdown

## install

```bash
go install github.com/go-kiss/protoc-gen-markdown
```

## generate markdown

```bash
protoc --markdown_out=Mhello.proto=./:. ./hello.proto
# set path prefix to /api
protoc --markdown_out=Mhello.proto=./,prefix=/api:. ./hello.proto
```
