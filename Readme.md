# IrisDockerSample
Docker環境上で動くIrisのサンプル

# How Usage
開発環境
```
docker-compose.yml
services:
  web:
    ...
    ports:
      - "80:80"
    volumes:                                     # ADD
      - [This Project Directory]:/opt/iris_dev/  # ADD
```

```
// Install gom
go get github.com/mattn/gom
// Install package
gom install
// Docker Running
docker-compose up -d
```

