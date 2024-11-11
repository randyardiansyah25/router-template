# Update go mod go version

```
go mod edit -go <go version>
go mod tidy
```

## BUILD WITH DOCKER

### BUILD

```docker
docker build -t ussiteam/router-template:1.0.0 . 
```

### PUSH

```docker
docker push ussiteam/router-template:1.0.0
```

## BUILD IN WINDOWS

### FOR LINUX

```shell
$env:GOOS = "linux"; $env:GOARCH="amd64"; go build -a -v -tags netgo -ldflags '-w' -o bin\router-template
```

### FOR WINDOWS

```shell
$env:GOOS = "windows"; go build -a -tags netgo -ldflags '-w' -o .\bin\router-template.exe
```
