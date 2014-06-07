namespace go models

struct UserGroup {
	1: string ID
	2: string Name
}

enum UserStatus {
	Banned (label="禁用")
	Offline (label="离线")
	Online (label="在线")
}

struct User {
	1: string ID
	2: required string UserName (label = "用户名", search = "Simple", requiredMsg = "请输入用户名")
	3: required string Password (label = "密码", widget = "password", requiredMsg = "请输入密码")
	4: string Name  (label = "姓名", search = "Name", search = "User-10")
	5: required string Email (label = "电邮", rule="", ruleMsg="电邮格式不正确", requiredMsg = "请输入电邮")
	6: string Intro (label = "介绍", search = "User-2")
	7: string Picture
	8: string Remark (widget = "richtext")
	9: bool IsAdmin (widget = "checkbox", label = "管理员")
	10: string UserGroupID
	11: i32 Status  (widget = "radio", label = "状态", dataSource="enum:UserStatus")
	12: list<string> Tags
}

