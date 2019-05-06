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