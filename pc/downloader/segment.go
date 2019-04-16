package downloader

type segmentState int

const (
	segmentWait segmentState = iota
	segmentDownloading
	segmentFinished
)

// 下载片段
type segment struct {
	start  int64 // 片段起始地址
	len    int64 // 片段长度
	finish int64 // 片段实际下载长度
	state  segmentState
}
