# gosoundio
Go port of libsoundio

## Building on macOS

You need to have [libsoundio](http://libsound.io/) installed. You can install it with Homebrew. (If you don't have Homebrew, [install it first](https://brew.sh/).)

```
$ brew install libsoundio
```

Install this package using `go get`:

```
$ go get github.com/mojbro/gosoundio
```

Try it out using an example:

```
$ cd $GOPATH/src/github.com/mojbro/gosoundio/examples
$ go run experiment.go
```

## Building on Linux, Windows, etc

I haven't tried it yet. If you manage to build it on Linux or Windows, please contribute by updating this README.
