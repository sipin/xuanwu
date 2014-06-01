namespace go models

struct UserGroup {
	1: string ID
	2: string Name
}

struct User {
	1: string ID
	2: string UserName (label = "用户名", search = "Simple")
	3: string Password (label = "密码")
	4: string Name  (label = "姓名", search = "Name", search = "User-10")
	5: string Email (label = "电邮")
	6: string Intro (label = "介绍", search = "User-2")
	7: string Picture
	8: string Remark
	9: bool IsAdmin
	10: string UserGroupID
	11: i32 Status				//用户状态，0:禁用; 1:离线; 2:在线
	12: string PubInfoID
	13: string OrganizationID
}

