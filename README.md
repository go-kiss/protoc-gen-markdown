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

## demo

You could check the [hello.md](hello.md) to see the result.

## doc

I have write a [Chinese article](https://taoshu.in/go/create-protoc-plugin.html) to introduce how it works.
