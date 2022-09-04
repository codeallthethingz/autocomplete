package disksort

import (
	"context"
	"fmt"
	"io"

	"github.com/lanrat/extsort"
)

func Sort(input chan string, w io.Writer) error {
	sort, outChan, errChan := extsort.Strings(input, extsort.DefaultConfig())
	sort.Sort(context.Background())
	previousLine := ""
	first := true
	previousCount := 1
	for line := range outChan {
		if !first {
			if line == previousLine {
				previousCount++
			} else {
				_, err := w.Write([]byte(fmt.Sprintf("%d %s\n", previousCount, previousLine)))
				if err != nil {
					return err
				}
				previousLine = line
				previousCount = 1
			}
		} else {
			previousLine = line
			first = false
		}
	}
	_, err := w.Write([]byte(fmt.Sprintf("%d %s\n", previousCount, previousLine)))
	if err != nil {
		return err
	}
	if err := <-errChan; err != nil {
		return err
	}
	return nil
}
