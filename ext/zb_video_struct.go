package ext

// uri
const (
	ZB_VIDEO_MAKE_QUERY = "/index.php?r=make/query"
)

const (
	ZB_VIDEO_SUCCESS = 1000

	ZB_VIDEO_NOT_FOUND = 1300
	ZB_VIDEO_MAKING    = 1100
	ZB_VIDEO_FAILED    = 1200
	ZB_VIDEO_ERROR     = 2000
)

var (
	ZB_VIDEO_RECODE_MAP = map[int]string{
		ZB_VIDEO_NOT_FOUND: "视频不存在",
		ZB_VIDEO_MAKING:    "视频制作中请稍后",
		ZB_VIDEO_FAILED:    "视频制作失败, 请重新制作",
		ZB_VIDEO_ERROR:     "发生其他错误, 请重新制作",
	}
)

type GetZbVideoReq struct {
	Code string `json:"Code"`
}

type ZbVideoInfo struct {
	Url  string `json:"url"`
	Path string `json:"path"`
}
type GetZbVideoRsp struct {
	State   int         `json:"state"`
	Message string      `json:"message"`
	Data    ZbVideoInfo `json:"data"`
}
