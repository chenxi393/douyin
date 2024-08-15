package response

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// 感觉最终可能还是用自己的json
// JSON 序列化的时候0 值会被忽略 FIXME
// 解决grpc返回成功状态码为0 会被忽略 字段为默认零值都会被忽略
// 24.4.2 弃用 rpc 返回的再封装打包 不直接序列化返回
func GrpcMarshal(v interface{}) ([]byte, error) {
	data, ok := v.(proto.Message)
	if ok {
		// FIXME protojson 会将64位的整数序列化成 字符串
		return protojson.MarshalOptions{EmitUnpopulated: true, UseProtoNames: true}.Marshal(data)
	}
	return json.Marshal(v)
}

type CommonResponse struct {
	StatusCode int         `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string      `json:"status_msg"`  // 返回状态信息
	Prompt     string      `json:"prompt"`      // 提示用户的
	Data       interface{} `json:"data,omitempty"`
}

func BuildStdResp(StatusCode int, StatusMsg string, data interface{}) *CommonResponse {
	return &CommonResponse{
		StatusCode: StatusCode,
		StatusMsg:  StatusMsg,
		Data:       data,
	}
}
