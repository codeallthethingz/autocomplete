package disksort

import (
	"context"
	"io"

	"github.com/lanrat/extsort"
)

func Sort(input chan string, w io.Writer, dedup bool) error {
	sort, outChan, errChan := extsort.Strings(input, extsort.DefaultConfig())
	sort.Sort(context.Background())
	previousLine := ""
	for line := range outChan {
		if dedup && line == previousLine {
			continue
		}
		previousLine = line
		w.Write([]byte(line + "\n"))
	}
	if err := <-errChan; err != nil {
		return err
	}
	return nil
}
