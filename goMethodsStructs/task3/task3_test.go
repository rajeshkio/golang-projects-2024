package task3

import (
	"testing"
)

func TestRunExample(t *testing.T) {
	t.Run("Test task3 pointers and receiver", func(t *testing.T) {
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
