package version

import (
	"testing"
)

func TestPrintInfo(t *testing.T) {
	t.Run("success - print info", func(t *testing.T) {
		PrintInfo()
	})
}
