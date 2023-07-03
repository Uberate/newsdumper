# The news getter

建议时间间隔大于 30 分钟。

# Quick start with default config

```shell
go run cmd/bin/main.go -c ./config/config.yaml -o output/ -Ft
```

# For docker

```shell
docker run -itd --name news-getter 
```

## Linux

### Source code

```shell
make build
make linux-install
```

如果需要进行商业使用，烦请联系:

- [ubserate@gmail.com](ubserate@gmail.com)

# 文档与说明

## 编码说明

[编码](docs/codes/README.md)
提供了一些工具方法，以快速接入某些接口。
