package response

import "C"
import "reflect"

type PointerResponse struct {
	reflectType reflect.Type
	data        []byte
}

func (p *PointerResponse) GetTypeString() string {
	return p.reflectType.String()
}

func (p *PointerResponse) GetData() []byte {
	return p.data
}

func ClientError() *PointerResponse {
	msg := "could not get client"
	return &PointerResponse{
		reflectType: reflect.TypeOf(msg),
		data:        []byte(msg),
	}
}

func Error(err error) *PointerResponse {
	return &PointerResponse{
		reflectType: reflect.TypeOf(err),
		data:        []byte(err.Error()),
	}
}

func StatusResponse() *PointerResponse {
	msg := "result status not successful"
	return &PointerResponse{
		reflectType: reflect.TypeOf(msg),
		data:        []byte(msg),
	}
}

func New(responseType reflect.Type, msg []byte) *PointerResponse {
	return &PointerResponse{
		reflectType: responseType,
		data:        msg,
	}
}

func NewBoolean(value bool) *PointerResponse {
	if value {
		return &PointerResponse{
			reflectType: reflect.TypeOf(value),
			data:        []byte{1},
		}
	} else {
		return &PointerResponse{
			reflectType: reflect.TypeOf(value),
			data:        []byte{0},
		}
	}
}

type StringResponse struct {
	reflectType reflect.Type
	str         string
}

func (p *StringResponse) GetTypeString() string {
	return p.reflectType.String()
}

func (p *StringResponse) GetValue() string {
	return p.str
}

func StringError(err error) *StringResponse {
	return &StringResponse{
		reflectType: reflect.TypeOf(err),
		str:         err.Error(),
	}
}

func StringStatusResponse() *StringResponse {
	msg := "result status not successful"
	return &StringResponse{
		reflectType: reflect.TypeOf(msg),
		str:         msg,
	}
}

func NewString(responseType reflect.Type, msg string) *StringResponse {
	return &StringResponse{
		reflectType: responseType,
		str:         msg,
	}
}
