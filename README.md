## Bluebell
> go version of my [api](https://api.dongfg.com)
-----
### Usage
Download latest release and run:
```shell script
./bluebell configs/config-example.yml
```

### Development
```
# enable vgo (will be default on after 1.13)
export GO111MODULE=on
# add package
go get xxx
```
config format:
```yaml
port: 20000
series:
  domain: www.rrys2019.com
```

### Release
```
goreleaser -f build/.goreleaser.yml
```