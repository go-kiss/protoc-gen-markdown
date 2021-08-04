# protoc-gen-markdown

## install

```bash
go install github.com/lvht/protoc-gen-markdown
```

## generate markdown

```bash
protoc --markdown_out=. hello.proto
# set path prefix to /api
protoc --markdown_out=prefix=/api:. hello.proto
```
