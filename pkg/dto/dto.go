package dto

type Response struct {
	Code     string      `json:"response_code"`
	Desc     string      `json:"response_desc"`
	Data     interface{} `json:"response_data"`
	HttpCode int         `json:"-"`
}

// Initialization of Response
func New() *Response {
	return &Response{
		Code: "XX",
		Desc: "General Error",
		Data: new(struct{}),
	}
}
func (r *Response) SetError(httpcode int, code, desc string, err error) {
	r.Code = code
	r.Desc = desc
	r.HttpCode = httpcode

	if err != nil {
		r.Data = map[string]interface{}{
			"error": err.Error(),
		}
	}
}

func (r *Response) SetSuccess(httpcode int, code, desc string, data interface{}) {
	r.Code = code
	r.Desc = desc
	r.Data = data
	r.HttpCode = httpcode
}

func NewError(httpcode int, code, desc string, err error) *Response {
	r := &Response{
		Code:     code,
		Desc:     desc,
		HttpCode: httpcode,
	}
	if err != nil {
		r.Data = map[string]interface{}{
			"error": err.Error(),
		}
	}

	return r

}
