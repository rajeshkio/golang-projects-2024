package task1

import (
	"testing"
)

func TestRunExample(t *testing.T) {
	t.Run("RunExample does not panic", func(t *testing.T) {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("RunExample panicked: %v", r)
				}
			}()
			RunExample()
		}()
	})
}
