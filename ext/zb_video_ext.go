package ext

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	
	"github.com/reechou/holmes"
	"github.com/reechou/robot-zbshipin/config"
)

type ZbVideoExt struct {
	key    string
	client *http.Client
	cfg    *config.Config
}

func NewZbVideoExt(cfg *config.Config) *ZbVideoExt {
	zbVideo := &ZbVideoExt{
		client: &http.Client{},
		cfg:    cfg,
	}
	
	return zbVideo
}

func (self *ZbVideoExt) GetZbVideo(request *GetZbVideoReq) (*GetZbVideoRsp, error) {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		holmes.Error("json encode error: %v", err)
		return nil, err
	}
	holmes.Debug("get zb video request: %s", string(reqBytes))
	
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s%s", self.cfg.ZBVideoHost.Host, ZB_VIDEO_MAKE_QUERY), bytes.NewBuffer(reqBytes))
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return nil, err
	}
	var response GetZbVideoRsp
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		holmes.Error("json decode error: %v [%s]", err, string(rspBody))
		return nil, err
	}
	
	return &response, nil
}
