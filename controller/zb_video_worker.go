package controller

import (
	"time"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-zbshipin/config"
	"github.com/reechou/robot-zbshipin/ext"
	"github.com/reechou/robot-zbshipin/robot_proto"
)

type ZbVideoWorker struct {
	cfg        *config.Config
	robotExt   *ext.RobotExt
	zbVideoExt *ext.ZbVideoExt

	MsgChan chan *robot_proto.ReceiveMsgInfo

	stop chan struct{}
	done chan struct{}
}

func NewZbVideoWorker(cfg *config.Config, robotExt *ext.RobotExt) *ZbVideoWorker {
	zbv := &ZbVideoWorker{
		cfg:      cfg,
		robotExt: robotExt,
		MsgChan:  make(chan *robot_proto.ReceiveMsgInfo, 1024),
		stop:     make(chan struct{}),
		done:     make(chan struct{}),
	}
	zbv.zbVideoExt = ext.NewZbVideoExt(cfg)
	go zbv.run()

	return zbv
}

func (self *ZbVideoWorker) run() {
	for {
		select {
		case msg := <-self.MsgChan:
			go self.runWorker(msg)
		case <-self.stop:
			close(self.done)
			return
		}
	}
}

func (self *ZbVideoWorker) Stop() {
	close(self.stop)
}

func (self *ZbVideoWorker) HandleMsg(msg *robot_proto.ReceiveMsgInfo) {
	select {
	case self.MsgChan <- msg:
	case <-self.stop:
		return
	}
}

func (self *ZbVideoWorker) runWorker(msg *robot_proto.ReceiveMsgInfo) {
	holmes.Debug("handle worker of msg: %s", msg.Msg)
	self.sendMsg(msg, robot_proto.RECEIVE_MSG_TYPE_TEXT, "正在查找视频中, 请稍后 ...")
	ok := self.check(msg)
	if ok {
		return
	}
	now := time.Now().Unix()
	for {
		select {
		case <-time.After(10 * time.Second):
			ok := self.check(msg)
			if ok {
				return
			}
			checkNow := time.Now().Unix()
			if checkNow-now > 120 {
				holmes.Error("msg[%v] cannot found zb video", msg)
				self.sendMsg(msg, robot_proto.RECEIVE_MSG_TYPE_TEXT, "未找到该视频哦")
				return
			}
		case <-self.stop:
			return
		}
	}
}

func (self *ZbVideoWorker) check(msg *robot_proto.ReceiveMsgInfo) bool {
	req := &ext.GetZbVideoReq{
		Code: msg.Msg,
	}
	rsp, err := self.zbVideoExt.GetZbVideo(req)
	if err != nil {
		holmes.Error("check zb video error: %v", err)
		return false
	}
	switch rsp.State {
	case ext.ZB_VIDEO_SUCCESS:
		holmes.Debug("get video success: %v", rsp)
		if rsp.Data.Path != "" {
			self.sendMsg(msg, robot_proto.RECEIVE_MSG_TYPE_VIDEO, rsp.Data.Path)
		} else if rsp.Data.Url != "" {
			self.sendMsg(msg, robot_proto.RECEIVE_MSG_TYPE_VIDEO, rsp.Data.Url)
		} else {
			holmes.Error("get zb video rsp[%v] not found video.", rsp)
		}
	case ext.ZB_VIDEO_FAILED, ext.ZB_VIDEO_NOT_FOUND, ext.ZB_VIDEO_ERROR:
		notifyMsg := ext.ZB_VIDEO_RECODE_MAP[rsp.State]
		holmes.Debug("video[%s] notify msg: %s", msg.Msg, notifyMsg)
		self.sendMsg(msg, robot_proto.RECEIVE_MSG_TYPE_TEXT, notifyMsg)
	case ext.ZB_VIDEO_MAKING:
		holmes.Debug("video[%s] is making.", msg.Msg)
		return false
	}

	return true
}

func (self *ZbVideoWorker) sendMsg(msg *robot_proto.ReceiveMsgInfo, msgType, retMSg string) error {
	var sendReq robot_proto.SendMsgInfo
	sendReq.SendMsgs = append(sendReq.SendMsgs, robot_proto.SendBaseInfo{
		WechatNick: msg.BaseInfo.WechatNick,
		ChatType:   robot_proto.FROM_TYPE_PEOPLE,
		UserName:   msg.BaseInfo.FromUserName,
		NickName:   msg.BaseInfo.FromNickName,
		MsgType:    msgType,
		Msg:        retMSg,
	})
	err := self.robotExt.SendMsgs("", &sendReq)
	if err != nil {
		holmes.Error("zb video worker send msg[%v] error: %v", sendReq, err)
		return err
	}
	return nil
}
