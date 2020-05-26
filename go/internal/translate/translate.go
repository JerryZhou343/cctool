package translate

type Translate interface {
	Do(src, from, to string) (dst string, err error)
}

type translate struct{}
