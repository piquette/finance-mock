# finance-mock

finance-mock is a mock HTTP server that can be used in lieu of various remote financial data sources. The primary purpose of this server is for building test suites for piquette/finance-go that don't have to interact with the real financial data sources, making testing quicker and less prone to unpredictable upstream api errors/changes out of our control.


## Usage

Get it from the Homebrew tap or download it [from the releases page][releases]:

``` sh
brew install ''

# start a finance-mock service at login
brew services start finance-mock

# upgrade if you already have it
brew upgrade ''
```

Or if you have Go installed you can build it:

``` sh
go get -u github.com/piquette/finance-mock
```

Run it:

``` sh
finance-mock
```

Or with docker:
``` sh
# build
docker build . -t finance-mock
# run
docker run -p 12111:12111 finance-mock
```

Then from another terminal:

``` sh
curl -i http://localhost:12111/v7/finance/quote\?symbols\=GE
```

By default, finance-mock runs on port 12111, but is configurable with the
`-port` option.

## Development

### Testing

Run the test suite:

``` sh
go test ./...
```

## Release

Release builds are generated with [goreleaser]. Make sure you have the software
and a `GITHUB_TOKEN`: set in your env.

``` sh
go get -u github.com/goreleaser/goreleaser
export GITHUB_TOKEN=...
```

Commit changes and tag `HEAD`:

``` sh
git tag v0.1.1
git push origin --tags
```

Then run goreleaser and you're done! Check [releases] (it also pushes to the
Homebrew tap).

``` sh
goreleaser --rm-dist
```

[go-bindata]: https://github.com/jteeuwen/go-bindata
[goreleaser]: https://github.com/goreleaser/goreleaser
[releases]: https://github.com/piquette/finance-mock/releases