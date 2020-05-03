package translate

import "github.com/JerryZhou343/cctool/internal/translate/common"

type Translate interface {
	Do(src string, channel common.Channel) (dst string, err error)
}

type translate struct{}
