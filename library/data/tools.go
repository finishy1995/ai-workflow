package data

import (
	"time"
)

func retry3Times(f func() error) error {
	var err error

	for i := 0; i < 3; i++ {
		err = f()
		if err == nil {
			return nil
		}
		time.Sleep(time.Millisecond * 100)
	}
	return err
}
