package storage

import (
	"github.com/golang/protobuf/proto"
	"github.com/peterq/pan-light/pc/dep"
	"github.com/peterq/pan-light/pc/util"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

var dataFilePath = dep.DataPath("pan-light.data")

var Global State

func init() {
	dep.OnInit(doInit)

	dep.OnClose(func() {
		log.Println("退出, 数据存盘...")
	})
}

func doInit() {
	defer writeToDiskTask()
	var last State
	debug := false
	if util.First(util.PathExists(dataFilePath)).(bool) {
		//解码数据
		var err error
		if debug {
			err = proto.Unmarshal(util.First(ioutil.ReadFile(dataFilePath)).([]byte), &last)
		} else {
			err = proto.Unmarshal(util.First(ioutil.ReadFile(dataFilePath)).([]byte), &Global)
		}
		if err != nil {
			dep.Fatal(err.Error())
		}
		if !debug {
			prepareState()
			return
		}
	}
	// 设置默认值
	Global.UserStorageMap = map[string]*UserStorage{}
	Global.DownloadList = map[string]*DownloadState{}
	Global.Settings = &StateSetting{
		DownloadSegSize:   1024 * 1024 * 2, // 2M
		DownloadCoroutine: 8,               // 8协程下载
	}
	// 初始化一个用户存储
	ChangeUserStorage()
	if debug {
		Global.UserStorageMap[Global.UserStorageId] = last.UserStorageMap[last.UserStorageId]
	}
	//log.Println(string(util.First(json.Marshal(Global)).([]byte)))
}

func writeToDiskTask() {
	fn := func() {
		bin, err := proto.Marshal(&Global)
		if err != nil {
			dep.Fatal(err.Error())
		}
		err = ioutil.WriteFile(dataFilePath, bin, 0655)
		if err != nil {
			dep.Fatal(err.Error())
		}
	}
	dep.OnClose(fn)
	go func() {
		for range time.Tick(10 * time.Second) {
			fn()
		}
	}()
}

// 下载数据锁
var downloadLockMap = map[string]*sync.Mutex{}

func LockForDownload(fid string) *sync.Mutex {
	return downloadLockMap[fid]
}

func NewTask(fid string, task *DownloadState) {
	Global.DownloadList[fid] = task
	downloadLockMap[fid] = new(sync.Mutex)
}

// 用户数据锁
var userStorageLockMap = map[string]*sync.Mutex{}

func LockForUserStorage(id string) *sync.Mutex {
	return userStorageLockMap[id]
}

// 从磁盘恢复状态时需要做一些初始化操作
func prepareState() {
	// 下载内容设置锁, 清空分配状态
	for fid, st := range Global.DownloadList {
		downloadLockMap[fid] = new(sync.Mutex)
		st.Downloading = false
		st.Stopping = false
		for _, seg := range st.Seg {
			seg.Distributed = false
		}
	}
	// 用户锁初始化
	for id := range Global.UserStorageMap {
		userStorageLockMap[id] = new(sync.Mutex)
	}
}

// 改变用户存储, 可以用来做账号切换
func ChangeUserStorage() {
	id := util.UniqId()
	userStorageLockMap[id] = new(sync.Mutex)
	us := new(UserStorage)
	Global.UserStorageMap[id] = us
	us.DataBucket = map[string]string{}
	Global.UserStorageId = id
}

func UserStorageSet(k, v string) {
	LockForUserStorage(Global.UserStorageId).Lock()
	defer LockForUserStorage(Global.UserStorageId).Unlock()
	Global.UserStorageMap[Global.UserStorageId].DataBucket[k] = v
}

func UserStorageGet(k string) string {
	LockForUserStorage(Global.UserStorageId).Lock()
	defer LockForUserStorage(Global.UserStorageId).Unlock()
	v, ok := Global.UserStorageMap[Global.UserStorageId].DataBucket[k]
	if ok {
		return v
	}
	return ""
}
