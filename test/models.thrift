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

const list<string> Departments = [
	"宣传部", "办公室", "工作室",
]

struct User {
	1: string ID
	2: required string UserName (label = "用户名", search = "Simple", requiredMsg = "请输入用户名")
	3: required string Password (label = "密码", widget = "password", requiredMsg = "请输入密码")
	4: string Name  (label = "姓名", search = "Name", search = "User-10")
	5: required string Email (label = "电邮", rule="^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$", ruleMsg="电邮格式不正确", requiredMsg = "请输入电邮")
	6: string Intro (label = "介绍", search = "User-2")
	7: string Picture
	8: string Remark
	9: bool IsAdmin (widget = "checkbox", label = "管理员")
	10: string UserGroupID (widget = "select", bindData = "UserGroup.Name")
	11: i32 Status (label = "状态", widget = "select", enums="UserStatus")
	12: string PubInfoID
	13: string OrganizationID
	14: string Department (label = "部门", widget = "selectPk", stringList="Departments")
}

