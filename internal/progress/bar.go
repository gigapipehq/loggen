package progress

import (
	"fmt"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

type Bar struct {
	*progressbar.ProgressBar
}

func (b *Bar) Add(count int) {
	_ = b.ProgressBar.Add(count)
}

func NewBar(max int, description string) *Bar {
	return &Bar{
		progressbar.NewOptions(max,
			progressbar.OptionSetDescription(description),
			progressbar.OptionSetWriter(os.Stdout),
			progressbar.OptionSetWidth(10),
			progressbar.OptionThrottle(65*time.Millisecond),
			progressbar.OptionShowCount(),
			progressbar.OptionOnCompletion(func() {
				_, _ = fmt.Fprint(os.Stdout, "\nAll batches sent. Waiting for all requests to complete...")
			}),
			progressbar.OptionSpinnerType(14),
			progressbar.OptionFullWidth(),
			progressbar.OptionSetRenderBlankState(true),
		),
	}
}
