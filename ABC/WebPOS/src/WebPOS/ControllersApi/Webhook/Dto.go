package Webhook

type ResponseDto struct {
	Success bool
	Code    int
	Msg     string
}

type RequestDto struct {
	ApiKey			string `form:"api_key"`
	Chain			string `form:"chain"`
	Group			string `form:"group"`
	Detail			string `form:"detail"`
}