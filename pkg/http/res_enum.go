package http

var (

	//account
	EmailAlreadyRegistered = R(10001, "邮箱已经注册")
	UsernamePwdNotMatch    = R(10002, "用户名和密码不匹配")
	EmailNotFound          = R(10003, "邮箱不存在，无法发送重置密码邮件")
	EmailFreqLimit         = R(10004, "邮件发送过于频繁，请一分钟后再试")
	ResetPwdTokenExpired   = R(10005, "重置链接失效，请重新发送邮件")
	WrongOldPassword       = R(10006, "旧密码错误")
	EmailRequired          = R(10007, "邮箱不能为空")
	PasswordRequired       = R(10008, "密码不能为空")
	EmailNotSupport        = R(10009, "不支持的邮箱")
	PwdNotConsistent       = R(10010, "两次输入密码不一致")
	PwdFormatError         = R(10011, "密码格式不正确")
	MixinTransfersNotFound = R(10012, "mixin transfers not found")
	WrongOldEmail          = R(10013, "旧邮箱错误")

	NotAuthorized = R(40001, "Dapp Token invalid")
	ServerError   = R(40020, "服务器异常，请稍后再试~")
	ParamError    = R(40021, "http request param error")
)
