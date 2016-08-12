# IrisDockerSample
Docker環境上で動くIrisのサンプル

# How Usage
開発環境

```
// Install gom
go get github.com/mattn/gom
// Install package
gom install
// running DB server
docker-compose up -d db
// running server
go run main.go
```

