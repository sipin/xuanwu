package models

import (
	//Official libs
	"bytes"
	"fmt"
	"strconv"

	//3rd party libs
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"github.com/sipin/gothrift/thrift"
	"github.com/mattbaird/elastigo/core"

	//Own libs
	"db"
)

type User struct {
	ID         bson.ObjectId `bson:"_id" thrift:"ID,1"`
	UserName  string `bson:"UserName" thrift:"UserName,2"`
	Password  string `bson:"Password" thrift:"Password,3"`
	Name  string `bson:"Name" thrift:"Name,4"`
	Email  string `bson:"Email" thrift:"Email,5"`
	Intro  string `bson:"Intro" thrift:"Intro,6"`
	Picture  string `bson:"Picture" thrift:"Picture,7"`
	Remark  string `bson:"Remark" thrift:"Remark,8"`
	IsAdmin  bool `bson:"IsAdmin" thrift:"IsAdmin,9"`
	UserGroupID  string `bson:"UserGroupID" thrift:"UserGroupID,10"`
	Status  int32 `bson:"Status" thrift:"Status,11"`
	PubInfoID  string `bson:"PubInfoID" thrift:"PubInfoID,12"`
	OrganizationID  string `bson:"OrganizationID" thrift:"OrganizationID,13"`
	widgets map[string]*Widget
}

type UserSearchSimpleObj struct {
	UserName string `json:"UserName"`
}
type UserSearchNameObj struct {
	Name string `json:"Name"`
}
type UserSearchUserObj struct {
	Name string `json:"Name"`
	Intro string `json:"Intro"`
}

func NewUser() *User {
	rval := new(User)
	rval.ID = bson.NewObjectId()
	rval.widgets = make(map[string]*Widget, 13)
	return rval
}

func NewUserFromBytes(data []byte) *User {
	o := new(User)
	transport := thrift.NewStreamTransportR(bytes.NewBuffer(data))
	protocol := thrift.NewTCompactProtocol(transport)
	o.Read(protocol)
	return o
}

func NewUserWithParams(params map[string]string) *User {
	o := new(User)
	o.ReadForm(params)
	return o
}

//Thrift Methods

func (p *User) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error: %s", p, err)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.readField2(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.readField3(iprot); err != nil {
				return err
			}
		case 4:
			if err := p.readField4(iprot); err != nil {
				return err
			}
		case 5:
			if err := p.readField5(iprot); err != nil {
				return err
			}
		case 6:
			if err := p.readField6(iprot); err != nil {
				return err
			}
		case 7:
			if err := p.readField7(iprot); err != nil {
				return err
			}
		case 8:
			if err := p.readField8(iprot); err != nil {
				return err
			}
		case 9:
			if err := p.readField9(iprot); err != nil {
				return err
			}
		case 10:
			if err := p.readField10(iprot); err != nil {
				return err
			}
		case 11:
			if err := p.readField11(iprot); err != nil {
				return err
			}
		case 12:
			if err := p.readField12(iprot); err != nil {
				return err
			}
		case 13:
			if err := p.readField13(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *User) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 1: %s", err)
	} else {
		p.ID = bson.ObjectId(v)
	}
	return nil
}
func (p *User) readField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 2: %s", err)
	} else {
		p.UserName = v
	}
	return nil
}
func (p *User) readField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 3: %s", err)
	} else {
		p.Password = v
	}
	return nil
}
func (p *User) readField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 4: %s", err)
	} else {
		p.Name = v
	}
	return nil
}
func (p *User) readField5(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 5: %s", err)
	} else {
		p.Email = v
	}
	return nil
}
func (p *User) readField6(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 6: %s", err)
	} else {
		p.Intro = v
	}
	return nil
}
func (p *User) readField7(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 7: %s", err)
	} else {
		p.Picture = v
	}
	return nil
}
func (p *User) readField8(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 8: %s", err)
	} else {
		p.Remark = v
	}
	return nil
}
func (p *User) readField9(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBool(); err != nil {
		return fmt.Errorf("error reading field 9: %s", err)
	} else {
		p.IsAdmin = v
	}
	return nil
}
func (p *User) readField10(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 10: %s", err)
	} else {
		p.UserGroupID = v
	}
	return nil
}
func (p *User) readField11(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return fmt.Errorf("error reading field 11: %s", err)
	} else {
		p.Status = v
	}
	return nil
}
func (p *User) readField12(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 12: %s", err)
	} else {
		p.PubInfoID = v
	}
	return nil
}
func (p *User) readField13(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 13: %s", err)
	} else {
		p.OrganizationID = v
	}
	return nil
}

func (p *User) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("User"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := p.writeField4(oprot); err != nil {
		return err
	}
	if err := p.writeField5(oprot); err != nil {
		return err
	}
	if err := p.writeField6(oprot); err != nil {
		return err
	}
	if err := p.writeField7(oprot); err != nil {
		return err
	}
	if err := p.writeField8(oprot); err != nil {
		return err
	}
	if err := p.writeField9(oprot); err != nil {
		return err
	}
	if err := p.writeField10(oprot); err != nil {
		return err
	}
	if err := p.writeField11(oprot); err != nil {
		return err
	}
	if err := p.writeField12(oprot); err != nil {
		return err
	}
	if err := p.writeField13(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("write struct stop error: %s", err)
	}
	return nil
}

func (p *User) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("ID", thrift.STRING, 1); err != nil {
		return fmt.Errorf("%T write field begin error string:ID: %s", p, err)
	}
	if err := oprot.WriteString(string(p.ID)); err != nil {
		return fmt.Errorf("%T.ID (string) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:ID: %s", p, err)
	}
	return err
}
func (p *User) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("UserName", thrift.STRING, 2); err != nil {
		return fmt.Errorf("%T write field begin error string:UserName: %s", p, err)
	}
	if err := oprot.WriteString(string(p.UserName)); err != nil {
		return fmt.Errorf("%T.UserName (string) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:ID: %s", p, err)
	}
	return err
}
func (p *User) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("Password", thrift.STRING, 3); err != nil {
		return fmt.Errorf("%T write field begin error string:Password: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Password)); err != nil {
		return fmt.Errorf("%T.Password (string) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:ID: %s", p, err)
	}
	return err
}
func (p *User) writeField4(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("Name", thrift.STRING, 4); err != nil {
		return fmt.Errorf("%T write field begin error string:Name: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Name)); err != nil {
		return fmt.Errorf("%T.Name (string) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:ID: %s", p, err)
	}
	return err
}
func (p *User) writeField5(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("Email", thrift.STRING, 5); err != nil {
		return fmt.Errorf("%T write field begin error string:Email: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Email)); err != nil {
		return fmt.Errorf("%T.Email (string) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:ID: %s", p, err)
	}
	return err
}
func (p *User) writeField6(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("Intro", thrift.STRING, 6); err != nil {
		return fmt.Errorf("%T write field begin error string:Intro: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Intro)); err != nil {
		return fmt.Errorf("%T.Intro (string) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:ID: %s", p, err)
	}
	return err
}
func (p *User) writeField7(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("Picture", thrift.STRING, 7); err != nil {
		return fmt.Errorf("%T write field begin error string:Picture: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Picture)); err != nil {
		return fmt.Errorf("%T.Picture (string) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:ID: %s", p, err)
	}
	return err
}
func (p *User) writeField8(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("Remark", thrift.STRING, 8); err != nil {
		return fmt.Errorf("%T write field begin error string:Remark: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Remark)); err != nil {
		return fmt.Errorf("%T.Remark (string) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:ID: %s", p, err)
	}
	return err
}
func (p *User) writeField9(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("IsAdmin", thrift.BOOL, 9); err != nil {
		return fmt.Errorf("%T write field begin error 9:IsAdmin: %s", p, err)
	}
	if err := oprot.WriteBool(bool(p.IsAdmin)); err != nil {
		return fmt.Errorf("%T.IsAdmin (bool) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:ID: %s", p, err)
	}
	return err
}
func (p *User) writeField10(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("UserGroupID", thrift.STRING, 10); err != nil {
		return fmt.Errorf("%T write field begin error string:UserGroupID: %s", p, err)
	}
	if err := oprot.WriteString(string(p.UserGroupID)); err != nil {
		return fmt.Errorf("%T.UserGroupID (string) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:ID: %s", p, err)
	}
	return err
}
func (p *User) writeField11(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("Status", thrift.I32, 11); err != nil {
		return fmt.Errorf("%T write field begin error i32:Status: %s", p, err)
	}
	if err := oprot.WriteI32(int32(p.Status)); err != nil {
		return fmt.Errorf("%T.Status (i32) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:ID: %s", p, err)
	}
	return err
}
func (p *User) writeField12(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("PubInfoID", thrift.STRING, 12); err != nil {
		return fmt.Errorf("%T write field begin error string:PubInfoID: %s", p, err)
	}
	if err := oprot.WriteString(string(p.PubInfoID)); err != nil {
		return fmt.Errorf("%T.PubInfoID (string) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:ID: %s", p, err)
	}
	return err
}
func (p *User) writeField13(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("OrganizationID", thrift.STRING, 13); err != nil {
		return fmt.Errorf("%T write field begin error string:OrganizationID: %s", p, err)
	}
	if err := oprot.WriteString(string(p.OrganizationID)); err != nil {
		return fmt.Errorf("%T.OrganizationID (string) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:ID: %s", p, err)
	}
	return err
}

func (p *User) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("User(%+v)", *p)
}

func (p *User) ToBytes() []byte {
	transport := thrift.NewTMemoryBuffer()
	protocol := thrift.NewTCompactProtocol(transport)
	p.Write(protocol)
	protocol.Flush()

	return transport.Bytes()
}

//mongo methods

func (o *User) Save() (info *mgo.ChangeInfo, err error) {
	session, col := db.GetCol("User")
	defer session.Close()

	if o.ID == "" {
		o.ID = bson.NewObjectId()
	}

	core.Index("user", "simple", o.ID.Hex(), nil, UserSearchSimpleObj{
		o.UserName,
	})
	core.Index("user", "name", o.ID.Hex(), nil, UserSearchNameObj{
		o.Name,
	})
	core.Index("user", "user", o.ID.Hex(), nil, UserSearchUserObj{
		o.Name,
		o.Intro,
	})


	return col.UpsertId(o.ID, o)
}

//Form methods

func (o *User) ReadForm(params map[string]string) {
	if val, ok := params["UserName"]; ok {
		o.UserName = val
	}
	if val, ok := params["Password"]; ok {
		o.Password = val
	}
	if val, ok := params["Name"]; ok {
		o.Name = val
	}
	if val, ok := params["Email"]; ok {
		o.Email = val
	}
	if val, ok := params["Intro"]; ok {
		o.Intro = val
	}
	if val, ok := params["Picture"]; ok {
		o.Picture = val
	}
	if val, ok := params["Remark"]; ok {
		o.Remark = val
	}
	if val, ok := params["IsAdmin"]; ok {
		o.IsAdmin = (val != "")
	} else {
		o.IsAdmin = false	
	}
	if val, ok := params["UserGroupID"]; ok {
		o.UserGroupID = val
	}
	if val, ok := params["Status"]; ok {
		intVal, _ := strconv.Atoi(val)
		o.Status = int32(intVal)
	}
	if val, ok := params["PubInfoID"]; ok {
		o.PubInfoID = val
	}
	if val, ok := params["OrganizationID"]; ok {
		o.OrganizationID = val
	}
}

func (o *User) UserNameWidget() *Widget {
	name := "UserName"
	ret, ok := o.widgets[name]
	if !ok || ret==nil {
		ret = &Widget{
			Label: "用户名",
			Value : o.UserName,
			Name: "UserName",
			PlaceHolder: "",
			Type: "string",
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) PasswordWidget() *Widget {
	name := "Password"
	ret, ok := o.widgets[name]
	if !ok || ret==nil {
		ret = &Widget{
			Label: "密码",
			Value : o.Password,
			Name: "Password",
			PlaceHolder: "",
			Type: "string",
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) NameWidget() *Widget {
	name := "Name"
	ret, ok := o.widgets[name]
	if !ok || ret==nil {
		ret = &Widget{
			Label: "姓名",
			Value : o.Name,
			Name: "Name",
			PlaceHolder: "",
			Type: "string",
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) EmailWidget() *Widget {
	name := "Email"
	ret, ok := o.widgets[name]
	if !ok || ret==nil {
		ret = &Widget{
			Label: "电邮",
			Value : o.Email,
			Name: "Email",
			PlaceHolder: "",
			Type: "string",
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) IntroWidget() *Widget {
	name := "Intro"
	ret, ok := o.widgets[name]
	if !ok || ret==nil {
		ret = &Widget{
			Label: "介绍",
			Value : o.Intro,
			Name: "Intro",
			PlaceHolder: "",
			Type: "string",
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) PictureWidget() *Widget {
	name := "Picture"
	ret, ok := o.widgets[name]
	if !ok || ret==nil {
		ret = &Widget{
			Label: "Picture",
			Value : o.Picture,
			Name: "Picture",
			PlaceHolder: "",
			Type: "string",
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) RemarkWidget() *Widget {
	name := "Remark"
	ret, ok := o.widgets[name]
	if !ok || ret==nil {
		ret = &Widget{
			Label: "Remark",
			Value : o.Remark,
			Name: "Remark",
			PlaceHolder: "",
			Type: "string",
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) IsAdminWidget() *Widget {
	name := "IsAdmin"
	ret, ok := o.widgets[name]
	if !ok || ret==nil {
		ret = &Widget{
			Label: "IsAdmin",
			Value : strconv.FormatBool(o.IsAdmin),
			Name: "IsAdmin",
			PlaceHolder: "",
			Type: "bool",
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) UserGroupIDWidget() *Widget {
	name := "UserGroupID"
	ret, ok := o.widgets[name]
	if !ok || ret==nil {
		ret = &Widget{
			Label: "UserGroupID",
			Value : o.UserGroupID,
			Name: "UserGroupID",
			PlaceHolder: "",
			Type: "string",
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) StatusWidget() *Widget {
	name := "Status"
	ret, ok := o.widgets[name]
	if !ok || ret==nil {
		ret = &Widget{
			Label: "Status",
			Value : strconv.FormatInt(int64(o.Status), 10),
			Name: "Status",
			PlaceHolder: "",
			Type: "i32",
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) PubInfoIDWidget() *Widget {
	name := "PubInfoID"
	ret, ok := o.widgets[name]
	if !ok || ret==nil {
		ret = &Widget{
			Label: "PubInfoID",
			Value : o.PubInfoID,
			Name: "PubInfoID",
			PlaceHolder: "",
			Type: "string",
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) OrganizationIDWidget() *Widget {
	name := "OrganizationID"
	ret, ok := o.widgets[name]
	if !ok || ret==nil {
		ret = &Widget{
			Label: "OrganizationID",
			Value : o.OrganizationID,
			Name: "OrganizationID",
			PlaceHolder: "",
			Type: "string",
		}
		o.widgets[name] = ret
	}

	return ret
}

func (o *User) Widgets() []*Widget {
	return []*Widget{
		o.UserNameWidget(),
		o.PasswordWidget(),
		o.NameWidget(),
		o.EmailWidget(),
		o.IntroWidget(),
		o.PictureWidget(),
		o.RemarkWidget(),
		o.IsAdminWidget(),
		o.UserGroupIDWidget(),
		o.StatusWidget(),
		o.PubInfoIDWidget(),
		o.OrganizationIDWidget(),
	}
}

//foreigh keys

func (p *User) GetUserGroup() (result *UserGroup, err error) {
	return UserGroupFindByID(p.UserGroupID)
}

func (p *User) SetUserGroup(obj *UserGroup) {
	p.UserGroupID = obj.ID.Hex()
}

func (o *UserGroup) GetAllUser() (result []*User, err error) {
	query := bson.M{"UserGroupID": o.ID.Hex()}
	return UserFindAll(query)
}
func (p *User) GetPubInfo() (result *PubInfo, err error) {
	return PubInfoFindByID(p.PubInfoID)
}

func (p *User) SetPubInfo(obj *PubInfo) {
	p.PubInfoID = obj.ID.Hex()
}

func (o *PubInfo) GetAllUser() (result []*User, err error) {
	query := bson.M{"PubInfoID": o.ID.Hex()}
	return UserFindAll(query)
}
func (p *User) GetOrganization() (result *Organization, err error) {
	return OrganizationFindByID(p.OrganizationID)
}

func (p *User) SetOrganization(obj *Organization) {
	p.OrganizationID = obj.ID.Hex()
}

func (o *Organization) GetAllUser() (result []*User, err error) {
	query := bson.M{"OrganizationID": o.ID.Hex()}
	return UserFindAll(query)
}


//Collection Manage methods

func UserFindOne(query interface{}, sortFields ...string) (result *User, err error) {
	session, col := db.GetCol("User")
	defer session.Close()

	q := col.Find(query)

	if sortFields == nil {
		q.Sort("-_id")
	} else {
		q.Sort(sortFields...)
	}

	err = q.One(&result)
	return
}

func UserFind(query interface{}, limit int, offset int, sortFields ...string) (result []*User, err error) {
	session, col := db.GetCol("User")
	defer session.Close()

	q := col.Find(query).Limit(limit).Skip(offset)

	if sortFields == nil {
		q.Sort("-_id")
	} else {
		q.Sort(sortFields...)
	}

	err = q.All(&result)
	return
}

func UserFindAll(query interface{}, sortFields ...string) (result []*User, err error) {
	session, col := db.GetCol("User")
	defer session.Close()

	q := col.Find(query)

	if sortFields == nil {
		q.Sort("-_id")
	} else {
		q.Sort(sortFields...)
	}

	err = q.All(&result)
	return
}

func UserCount(query interface{}) (result int) {
	session, col := db.GetCol("User")
	defer session.Close()

	result, _ = col.Find(query).Count()
	return
}

func UserFindByID(id string) (result *User, err error) {
	session, col := db.GetCol("User")
	defer session.Close()

	if !bson.IsObjectIdHex(id) {
		err = ErrInvalidObjectId
		return
	}
	err = col.FindId(bson.ObjectIdHex(id)).One(&result)
	return
}

func UserRemoveByID(id string) (result *User, err error) {
	session, col := db.GetCol("User")
	defer session.Close()

	err = col.RemoveId(bson.ObjectIdHex(id))
	core.Delete("user", "simple", id, nil)
	core.Delete("user", "name", id, nil)
	core.Delete("user", "user", id, nil)
	return
}

// Search
func UserSearchSimple(word string, limit int, offset int) (core.SearchResult, error) {
	searchJson := `{
    "query" : {
        "query_string" :  {
	      "default_operator": "OR",
	      "fields": ` + `["UserName"]` + `,
	      "query": "` + word + `"
	    }
    }
}`
	args := map[string]interface{} {"from" : offset, "size": limit}
	return core.SearchRequest("user", "simple", args, searchJson)
}
func UserSearchName(word string, limit int, offset int) (core.SearchResult, error) {
	searchJson := `{
    "query" : {
        "query_string" :  {
	      "default_operator": "OR",
	      "fields": ` + `["Name"]` + `,
	      "query": "` + word + `"
	    }
    }
}`
	args := map[string]interface{} {"from" : offset, "size": limit}
	return core.SearchRequest("user", "name", args, searchJson)
}
func UserSearchUser(word string, limit int, offset int) (core.SearchResult, error) {
	searchJson := `{
    "query" : {
        "query_string" :  {
	      "default_operator": "OR",
	      "fields": ` + `["Name", "Intro"]` + `,
	      "query": "` + word + `"
	    }
    }
}`
	args := map[string]interface{} {"from" : offset, "size": limit}
	return core.SearchRequest("user", "user", args, searchJson)
}
