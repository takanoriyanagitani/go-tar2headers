package wtr

import (
	"iter"

	th "github.com/takanoriyanagitani/go-tar2headers"
	. "github.com/takanoriyanagitani/go-tar2headers/util"
)

type HeaderWriter func(th.Header) IO[Void]

type HeadersWriter func(iter.Seq2[th.Header, error]) IO[Void]
