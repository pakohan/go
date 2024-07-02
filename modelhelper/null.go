package modelhelper

func Coalesce[T any](n *T) T {
	if n == nil {
		var n T
		return n
	}

	return *n
}
