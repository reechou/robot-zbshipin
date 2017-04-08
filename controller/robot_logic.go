package controller

import (
	"github.com/reechou/holmes"
	"github.com/reechou/robot-zbshipin/robot_proto"
)

func (self *Logic) HandleReceiveMsg(msg *robot_proto.ReceiveMsgInfo) {
	holmes.Debug("receive robot msg: %v", msg)
	switch msg.BaseInfo.ReceiveEvent {
	case robot_proto.RECEIVE_EVENT_MSG:
		self.handleMsg(msg)
	}
}

func (self *Logic) handleMsg(msg *robot_proto.ReceiveMsgInfo) {
	switch msg.MsgType {
	case robot_proto.RECEIVE_MSG_TYPE_TEXT:
		self.zbw.HandleMsg(msg)
	}
}
