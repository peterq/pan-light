package downloader

// 下载片段
type segment struct {
	start  int64 // 片段起始地址
	len    int64 // 片段长度
	finish int64 // 片段实际下载长度, 需要确保只在下载该片段过程中使用该字段
}
