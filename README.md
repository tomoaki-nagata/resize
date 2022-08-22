# resize

## Development

```
docker run --rm -it -v $(pwd):/go/src/app -w /go/src/app golang:1.19.0-bullseye sh -c 'go mod tidy && gofmt -l -s -w . && go build -o resize+in+1000x1000 && ./resize+in+1000x1000'
```

## Build for Intel Mac

```
docker run --rm -it -v $(pwd):/go/src/app -w /go/src/app -e GOOS=darwin -e GOARCH=amd64 golang:1.19.0-bullseye sh -c 'go mod tidy && gofmt -l -s -w . && go build -o resizeamd+in+1000x1000'
```

## Build for Apple Silicon

```
docker run --rm -it -v $(pwd):/go/src/app -w /go/src/app -e GOOS=darwin -e GOARCH=arm64 golang:1.19.0-bullseye sh -c 'go mod tidy && gofmt -l -s -w . && go build -o resizearm+in+1000x1000'
```

## Note

| architecture         | setup | dev | cpu | sec | plat | note      |
| -------------------- | :---: | :-: | :-: | :-: | :--: | :-------: |
| Golang (imaging)     | o     | o   | x   | o   | o    |           |
| Golang (image)       | o     | x   | x   | o   | o    |           |
| Golang (bimg)        | x     | o   | o   | o   | o    |           |
| ImageMagick (shell)  | x     | o   | x   | x   | o    |           |
| libvips (shell)      | x     | ?   | o   | o   | x    |           |
| ImageMagick (docker) | x     | o   | o   | x   | o    |           |
| libvips (docker)     | x     | ?   | o   | o   | o    |           |
| Automator            | o     | o   | x   | o   | x    | x quality |
