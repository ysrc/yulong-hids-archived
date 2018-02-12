package monitor

type filterInfo struct {
	Port    []int
	Process []string
	File    []string
}

var filter filterInfo

const (
	fileSize int64 = 20480000
	UDP      uint8 = 17
	TCP      uint8 = 6
)

func init() {
	// 硬编码白名单
	filter.Port = []int{137, 139, 445}
	filter.File = []string{`c:\windows\temp`}
}
