package pagination

import (
	"errors"
)

var ErrOffsetRange = errors.New("offset is out of range")

func Pagination[T any](messages []T, limit, offset int) ([]T, error) {
	if offset >= len(messages) {
		return nil, ErrOffsetRange
	}

	end := offset + limit

	if end > len(messages) {
		end = len(messages)
	}

	return messages[offset:end], nil
}
