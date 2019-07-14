package main

import (
	"encoding/json"
	"fmt"
	"github.com/Tnze/CoolQ-Golang-SDK/cqp/util"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
)

func main() {}
func init() {
	cqp.AppID = "qqbot.pan-light.peterq.cn"
	cqp.Enable = onEnable
	cqp.Disable = onDisable
	cqp.GroupMsg = onGroupMsg
	cqp.PrivateMsg = onPrivateMsg
	cqp.GroupMemberDecrease = onMemberLeave
}

func onPrivateMsg(subType, msgID int32, fromQQ int64, msg string, font int32) int32 {
	defer handleErr()
	state.handlePrivateMsg(fromQQ, msg)
	return 0
}

func onMemberLeave(sub_type, send_time int32, from_group, from_qq, being_operate_qq int64) int32 {
	if from_group == qqGroup {
		state.handleUserMove(from_qq)
	}
	return 0
}

const qq = int64(1617085421)

//const qqGroup = int64(767598744) 测试群
const qqGroup = int64(438604465)
const stateFile = "state-%d-%d.json"

var state State

func onEnable() int32 {
	cqp.AddLog(cqp.Debug, "bot start", time.Now().String())
	state.QGroup = qqGroup
	state.QQ = qq
	state.File = fmt.Sprintf(stateFile, qq, qqGroup)
	state.Answer = "10"
	state.OnError = func(i ...interface{}) {
		cqp.AddLog(cqp.Error, "manager error", fmt.Sprint(i...))
	}
	state.Log = func(i ...interface{}) {
		cqp.AddLog(cqp.Info, "manager log", fmt.Sprint(i...))
	}
	state.init()
	return 0
}

func onDisable() int32 {
	defer handleErr()
	//将宠物们保存到pets.json
	cqp.AddLog(cqp.Debug, "bot disable", time.Now().String())
	state.saveTicker.Stop()
	state.save()
	return 0
}

func onGroupMsg(subType, msgID int32, fromGroup, fromQQ int64, fromAnonymous, msg string, font int32) int32 {
	defer handleErr()
	if fromGroup == qqGroup {
		state.handleGroupMsg(fromQQ, msg)
	}
	return 0
}

func handleErr() {
	if err := recover(); err != nil {
		cqp.AddLog(cqp.Fatal, "严重错误", fmt.Sprint(err))
	}
}

type State struct {
	File          string
	Answer        string
	QQ            int64
	QGroup        int64
	OnError       func(...interface{})
	Log           func(...interface{})
	Members       map[int64]*Member
	NextOutBatch  []*Member
	ConfirmString string
	saveTicker    *time.Ticker
}

type Member struct {
	Info     util.GroupMember
	Correct  bool
	BeenOut  bool
	TimeLeft int
}

func (s *State) init() {
	rand.Seed(time.Now().UnixNano())
	if s.saveTicker != nil {
		s.saveTicker.Stop()
	}
	s.saveTicker = time.NewTicker(10 * time.Second)
	go func() {
		for range s.saveTicker.C {
			s.save()
		}
	}()
	s.Members = map[int64]*Member{}
	bytes, e := ioutil.ReadFile(s.File)
	if e != nil {
		s.OnError("读取state失败")
		s.getMembers()
	} else {
		e = json.Unmarshal(bytes, &s.Members)
		if e != nil {
			s.OnError("恢复state失败")
			s.getMembers()
		}
	}
}

func (s *State) getMembers() {
	str := cqp.GetGroupMemberList(s.QGroup)
	members, e := util.UnpackGroupList(str)
	if e != nil {
		s.OnError("获取群成员错误", e)
		return
	}
	for _, m := range members {
		s.Members[m.QQ] = &Member{
			Info:     m,
			Correct:  false,
			BeenOut:  false,
			TimeLeft: 2,
		}
	}
	s.Log("members", members)
}

func (s *State) save() {
	bytes, e := json.Marshal(s.Members)
	if e != nil {
		return
	}
	ioutil.WriteFile(s.File, bytes, os.ModePerm)
}

func (s *State) handleGroupMsg(fromQQ int64, msg string) {
	msg = strings.Trim(msg, " ")
	member, ok := state.Members[fromQQ]
	if !ok {
		return
	}
	if member.Info.Auth == 1 { // 普通成员不处理
		return
	}
	cmd := strings.Split(msg, " ")
	switch cmd[0] {
	case "踢人":
		s.cmdOut(cmd)
	case "确认":
		s.cmdConfirm(cmd)
	}
}

type MemberSlice []*Member

func (s MemberSlice) Len() int {
	return len(s)
}

func (s MemberSlice) Less(i, j int) bool {
	if s[i].TimeLeft != s[i].TimeLeft {
		return s[i].TimeLeft > s[j].TimeLeft
	}
	return s[i].Info.LastChat.Unix() < s[j].Info.LastChat.Unix()
}

func (s MemberSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *State) cmdOut(cmd []string) {
	if len(cmd) != 2 {
		cqp.SendGroupMsg(s.QGroup, "踢人参数不正确. eg: [踢人 3] 剔除 3 个回答错误的群员")
		return
	}
	n, e := strconv.Atoi(cmd[1])
	if e != nil {
		cqp.SendGroupMsg(s.QGroup, "踢人参数不正确. eg: [踢人 3] 剔除 3 个回答错误的群员")
		return
	}
	msg := "以下群员将被剔除:"

	var sorted = make(MemberSlice, len(state.Members))
	i := 0
	for _, m := range state.Members {
		sorted[i] = m
		i++
	}
	sort.Sort(sorted)
	var members []*Member
	for _, m := range sorted {
		if len(members) >= n {
			break
		}
		if m.Info.Auth > 1 || m.Correct || m.BeenOut {
			continue
		}
		if m.Info.QQ == 295851641 {
			continue
		}
		members = append(members, m)
		reason := "回答错误"
		if m.TimeLeft == 2 {
			reason = "没有回答"
		}
		msg += fmt.Sprint(m.Info.QQ, " ", reason, "\n")
	}
	if len(members) == 0 {
		cqp.SendGroupMsg(s.QGroup, "没有可移除的群员")
		return
	}
	state.NextOutBatch = members
	r := rand.Intn(10000-1000) + 1000
	state.ConfirmString = strconv.Itoa(r)
	msg += "请管理员回复[确认 " + state.ConfirmString + "] 开始剔除以上成员"
	cqp.SendGroupMsg(s.QGroup, msg)
}

func (s *State) cmdConfirm(cmd []string) {
	if len(cmd) != 2 || cmd[1] == "" {
		return
	}
	if cmd[1] != state.ConfirmString {
		return
	}
	cqp.SendGroupMsg(s.QGroup, "管理员已确认, 开始踢人")
	s.Log("开始踢人")
	s.ConfirmString = ""
	var members = s.NextOutBatch
	go func() {
		i := 0
		for _, m := range members {
			m.BeenOut = true
			if cqp.SetGroupKick(s.QGroup, m.Info.QQ, false) == 0 {
				i++
			}
		}
		cqp.SendGroupMsg(s.QGroup, fmt.Sprintf("已剔除%d名成员", i))
	}()
}

func (s *State) handlePrivateMsg(fromQQ int64, msg string) {
	msg = strings.Trim(msg, " ")
	member, ok := state.Members[fromQQ]
	if !ok {
		return
	}
	if member.TimeLeft == 0 || member.Correct {
		return
	}
	if member.Info.Auth != 1 {
		cqp.SendPrivateMsg(fromQQ, "管理员无需回答")
		return
	}
	member.TimeLeft--
	if msg == s.Answer {
		member.Correct = true
		cqp.SendPrivateMsg(fromQQ, "回答正确, 您将继续留在开发群")
	} else {
		if member.TimeLeft == 0 {
			cqp.SendPrivateMsg(fromQQ, "回答错误, 您将被移除群聊")
		} else {
			cqp.SendPrivateMsg(fromQQ, fmt.Sprintf("回答错误, 您还剩%d次机会", member.TimeLeft))
		}
	}
}

func (s *State) handleUserMove(fromQQ int64) {
	member, ok := state.Members[fromQQ]
	if !ok {
		return
	}
	member.BeenOut = true
}
