package utils

import "time"

func Retry(f func() error) (err error) {
	if err = f(); err == nil {
		return
	}

	m := 2 * time.Second
	t := m

	for i := 1; i < 3; i++ {
		t *= m
		time.Sleep(t)

		if err = f(); err == nil {
			return
		}
	}

	return
}
