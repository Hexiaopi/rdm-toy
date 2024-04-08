package retcode

var (
	//router
	Success                = NewCode("000000", "Success")
	RequestPathNotFound    = NewCode("000001", "Request Path Not Found")
	RequestMethodNotAllow  = NewCode("000002", "Request Method Not Allowed")
	RequestIllegal         = NewCode("000003", "Request Illegal")
	RequestUnMarshalError  = NewCode("000004", "Request Body UnMarshal Error")
	RequestNoPermission    = NewCode("000005", "Request No Permission")
	RequestAuthCheckFail   = NewCode("000006", "Request Auth Check Fail")
	GenerateAuthTokenFail  = NewCode("000007", "Generate Auth Token Fail")
	RequestTokenEmpty      = NewCode("000008", "Request Token Empty")
	RequestTokenWrong      = NewCode("000009", "Request Token Wrong")
	RequestTokenAuthFail   = NewCode("000010", "Request Token Auth Fail")
	RequestTokenAuthExpire = NewCode("000011", "Request Token Expire")
	RequestUserGetFail     = NewCode("000012", "Get User Fail")

	//module
	ListConnFail   = NewCode("101001", "List Redis Conn Fail")
	GetConnFail    = NewCode("101002", "Get Redis Conn Fail")
	CreateConnFail = NewCode("101003", "Create Redis Conn Fail")
	UpdateConnFail = NewCode("101004", "Update Redis Conn Fail")
	DeleteConnFail = NewCode("101005", "Delete Redis Conn Fail")

	ListDBFail = NewCode("101101", "List Redis DB Fail")
	GetDBFail  = NewCode("101102", "Get Redis DB Fail")

	ListKeyFail     = NewCode("101201", "List Redis Key Fail")
	GetKeyFail      = NewCode("101202", "Get Redis Key Fail")
	CreateKeyFail   = NewCode("101203", "Create Redis Key Fail")
	UpdateKeyFail   = NewCode("101204", "Update Redis Key Fail")
	DeleteKeyFail   = NewCode("101205", "Delete Redis Key Fail")
	ExistKeyAlready = NewCode("101206", "Redis Key Already Exist")

	UnknownError = NewCode("999999", "Unknown Error")
)
