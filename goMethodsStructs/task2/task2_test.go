package task2

import "testing"

func TestRunExample(t *testing.T) {
	t.Run("Testing task2", func(t *testing.T) {
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
