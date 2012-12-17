json
====
This project implements a JSON parser in Go

This file provides a custom JSON parser suitable for processing Twitter data.
Some differences from the standard Golang json package:
  * Does not use reflection, parses into standard map/slice/value structs.
  * Parses numbers into int64 where possible, float64 otherwise.
  * Faster!  Probably due to no reflection:

    BenchmarkParseTweet	   10000	    155714 ns/op
    BenchmarkCustomJSON	   20000	     85822 ns/op

Installing
----------
Run

    go get github.com/kurrik/json

Include in your source:

    import "github.com/kurrik/json"
