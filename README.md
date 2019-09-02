## Bluebell
> go version of my [api](https://api.dongfg.com)
-----
### Development
```
# enable vgo (will be default on after 1.13)
export GO111MODULE=on
# add package
go get xxx
```
config format:
```yaml
port: 9001
service:
  name: bluebell
  address: 10.64.60.246
  port: 9001
  check-url: http://10.64.60.246:9001/health
  check-interval: 100s
series:
  domain: www.zmz2019.com
```

### Release
```
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
goreleaser
```

### Changelog
#### v0.1.5
only keep production server
#### v0.1.4
enable cgo in docker
#### v0.1.3
call `packr2` before build
#### v0.1.2
fix gorelease Deprecation: [archive](https://goreleaser.com/deprecations/#archive)
#### v0.1.1
Package static file with `packr`
#### v0.1.0
First release