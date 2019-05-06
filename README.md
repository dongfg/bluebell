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

### Release
```
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
goreleaser
```

### Changelog
#### v0.1.3
call `packr2` before build
#### v0.1.2
fix gorelease Deprecation: [archive](https://goreleaser.com/deprecations/#archive)
#### v0.1.1
Package static file with `packr`
#### v0.1.0
First release