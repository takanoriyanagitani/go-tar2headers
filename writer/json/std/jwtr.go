package jwtr

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"iter"
	"os"

	th "github.com/takanoriyanagitani/go-tar2headers"
	. "github.com/takanoriyanagitani/go-tar2headers/util"
	tw "github.com/takanoriyanagitani/go-tar2headers/writer"
)

func WriterToHeadersWriter(wtr io.Writer) tw.HeadersWriter {
	return func(hdrs iter.Seq2[th.Header, error]) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			var bw *bufio.Writer = bufio.NewWriter(wtr)
			defer bw.Flush()

			var enc *json.Encoder = json.NewEncoder(bw)

			for hdr, e := range hdrs {
				select {
				case <-ctx.Done():
					return Empty, ctx.Err()
				default:
				}

				if nil != e {
					return Empty, e
				}

				e := enc.Encode(hdr)
				if nil != e {
					return Empty, e
				}
			}

			return Empty, nil
		}
	}
}

var HeadersWriterStdout tw.HeadersWriter = WriterToHeadersWriter(os.Stdout)
