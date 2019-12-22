# protoc-gen-markdown

## install

```bash
go get github.com/lvht/protoc-gen-markdown
```

## generate markdown

```bash
protoc --markdown_out=. hello.proto
# set path prefix to /api
protoc --markdown_out=path_prefix=/api:. hello.proto
```

## todo
- [ ] support enum
- [ ] support nested message
