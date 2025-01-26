package main

import (
	"bufio"
	"context"
	"iter"
	"log"
	"os"

	th "github.com/takanoriyanagitani/go-tar2headers"
	. "github.com/takanoriyanagitani/go-tar2headers/util"
	tw "github.com/takanoriyanagitani/go-tar2headers/writer"
	js "github.com/takanoriyanagitani/go-tar2headers/writer/json/std"
)

var stdin2headers IO[iter.Seq2[th.Header, error]] = Of(
	th.ReaderToHeaders(bufio.NewReader(os.Stdin)),
)

var headersWriter tw.HeadersWriter = js.HeadersWriterStdout

var stdin2headers2stdout IO[Void] = Bind(
	stdin2headers,
	headersWriter,
)

var sub IO[Void] = func(ctx context.Context) (Void, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return stdin2headers2stdout(ctx)
}

func main() {
	_, e := sub(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
