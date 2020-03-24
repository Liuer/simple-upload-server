# a deadly simple http upload server

## build
```
# build
go build -o simple-upload-server main.go
```

## run
```
# run default
./simple-upload-server
# then open browser input http://localhost:9090
# application will make a folder named 'upload' (if not exist) at current dir where the uploaded file were 

# run with specify port
./simple-upload-server -p 8086
# then open browser input http://localhost:8086

```