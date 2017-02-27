# uuid for golang

Generator of uuid of [RF4112](http://www.ietf.org/rfc/rfc4122.txt) uuid

Inpsired by https://github.com/kelektiv/node-uuid

Features :

- generate RFC4112 uuid version 1 and 4
- uses go crpyto/rand for random number generation
- Version 1 based on timestamp
- Version 4 based on random numbers

## Installation

`go get github.com/retiredbatman/gouuid`

## Example
```go
u1 := NewV1()
fmt.Printf("UUIDv1: %s\n", u1.Format())

u4 := NewV4()
fmt.Printf("UUIDv4: %s\n", u4.Format())
```

