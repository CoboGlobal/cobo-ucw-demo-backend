# Cobo UCW Backend Demo

## Automated Initialization (wire)
```
# install wire
go get github.com/google/wire/cmd/wire

# generate wire
cd cmd/cobo-ucw-backend
wire
```

## Generate Api
```shell
make api
```

## Generate Go Sdk
```shell
make sdk-go
```

## Docker
```bash
# build
docker build -t ucw-backend .

# run
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf ucw-backend
```
