package pair

type Pair[L any, R any] struct {
	Left  L
	Right R
}

func Map[R any, T any](p Pair[error, R], mapper func(R) T) Pair[error, T] {
	if nil != p.Left {
		return Pair[error, T]{Left: p.Left}
	}
	var t T = mapper(p.Right)
	return Pair[error, T]{Right: t}
}

func PairNew[L, R any](right R, left L) Pair[L, R] {
	return Pair[L, R]{
		Left:  left,
		Right: right,
	}
}

func AndThen[R any, T any](
	p Pair[error, R],
	mapper func(R) Pair[error, T],
) Pair[error, T] {
	if nil != p.Left {
		return Pair[error, T]{Left: p.Left}
	}
	return mapper(p.Right)
}

func Left[L, R any](left L) Pair[L, R]   { return Pair[L, R]{Left: left} }
func Right[L, R any](right R) Pair[L, R] { return Pair[L, R]{Right: right} }

func UnwrapOr[R any](
	p Pair[error, R],
	alt R,
) R {
	switch p.Left {
	case nil:
		return p.Right
	default:
		return alt
	}
}
