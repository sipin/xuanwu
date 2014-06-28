package models

import (
	//Official libs
	"bytes"
	"fmt"
	"models/test"
	"regexp"
	"strconv"

	//3rd party libs
	"github.com/sipin/gothrift/thrift"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	//Own libs
	"db"
)

func init() {
	db.SetOnFinishInit(initUserIndex)
}

func initUserIndex() {
	session, collection := db.GetCol(UserTableName)
	defer session.Close()

reEnsureUserName:
	if err := collection.EnsureIndex(mgo.Index{
		Key:    []string{"UserName"},
		Unique: true,
		Sparse: true,
	}); err != nil {
		println("error ensureIndex User UserName", err.Error())
		err = collection.DropIndex("UserName")
		if err != nil {
			panic(err)
		}
		goto reEnsureUserName

	}

}

var UserTableName = "User"

type User struct {
	ID             bson.ObjectId `bson:"_id,omitempty" thrift:"ID,1"`
	UserName       string        `bson:"UserName" thrift:"UserName,2"`
	Password       string        `bson:"Password" thrift:"Password,3"`
	Name           string        `bson:"Name" thrift:"Name,4"`
	Email          string        `bson:"Email" thrift:"Email,5"`
	Intro          string        `bson:"Intro" thrift:"Intro,6"`
	Picture        string        `bson:"Picture" thrift:"Picture,7"`
	Remark         string        `bson:"Remark" thrift:"Remark,8"`
	IsAdmin        bool          `bson:"IsAdmin" thrift:"IsAdmin,9"`
	UserGroupID    string        `bson:"UserGroupID" thrift:"UserGroupID,10"`
	Status         int32         `bson:"Status" thrift:"Status,11"`
	PubInfoID      string        `bson:"PubInfoID" thrift:"PubInfoID,12"`
	OrganizationID string        `bson:"OrganizationID" thrift:"OrganizationID,13"`
	Gender         string        `bson:"Gender" thrift:"Gender,14"`
	Index          int32         `bson:"Index" thrift:"Index,15"`
	widgets        map[string]*Widget
}

func NewUser() *User {
	rval := new(User)
	rval.ID = bson.NewObjectId()
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
		case 14:
			if err := p.readField14(iprot); err != nil {
				return err
			}
		case 15:
			if err := p.readField15(iprot); err != nil {
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

func (p *User) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("ID", thrift.STRING, 1); err != nil {
		return fmt.Errorf("%T write field begin error 1:ID: %s", p, err)
	}
	if err := oprot.WriteString(string(p.ID)); err != nil {
		return fmt.Errorf("%T.ID (1) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:ID: %s", p, err)
	}
	return err
}

func (p *User) readField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 2: %s", err)
	} else {
		p.UserName = v
	}
	return nil
}

func (p *User) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("UserName", thrift.STRING, 2); err != nil {
		return fmt.Errorf("%T write field begin error 2:UserName: %s", p, err)
	}
	if err := oprot.WriteString(string(p.UserName)); err != nil {
		return fmt.Errorf("%T.UserName (2) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 2:UserName: %s", p, err)
	}
	return err
}

func (p *User) readField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 3: %s", err)
	} else {
		p.Password = v
	}
	return nil
}

func (p *User) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("Password", thrift.STRING, 3); err != nil {
		return fmt.Errorf("%T write field begin error 3:Password: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Password)); err != nil {
		return fmt.Errorf("%T.Password (3) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 3:Password: %s", p, err)
	}
	return err
}

func (p *User) readField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 4: %s", err)
	} else {
		p.Name = v
	}
	return nil
}

func (p *User) writeField4(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("Name", thrift.STRING, 4); err != nil {
		return fmt.Errorf("%T write field begin error 4:Name: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Name)); err != nil {
		return fmt.Errorf("%T.Name (4) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 4:Name: %s", p, err)
	}
	return err
}

func (p *User) readField5(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 5: %s", err)
	} else {
		p.Email = v
	}
	return nil
}

func (p *User) writeField5(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("Email", thrift.STRING, 5); err != nil {
		return fmt.Errorf("%T write field begin error 5:Email: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Email)); err != nil {
		return fmt.Errorf("%T.Email (5) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 5:Email: %s", p, err)
	}
	return err
}

func (p *User) readField6(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 6: %s", err)
	} else {
		p.Intro = v
	}
	return nil
}

func (p *User) writeField6(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("Intro", thrift.STRING, 6); err != nil {
		return fmt.Errorf("%T write field begin error 6:Intro: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Intro)); err != nil {
		return fmt.Errorf("%T.Intro (6) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 6:Intro: %s", p, err)
	}
	return err
}

func (p *User) readField7(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 7: %s", err)
	} else {
		p.Picture = v
	}
	return nil
}

func (p *User) writeField7(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("Picture", thrift.STRING, 7); err != nil {
		return fmt.Errorf("%T write field begin error 7:Picture: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Picture)); err != nil {
		return fmt.Errorf("%T.Picture (7) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 7:Picture: %s", p, err)
	}
	return err
}

func (p *User) readField8(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 8: %s", err)
	} else {
		p.Remark = v
	}
	return nil
}

func (p *User) writeField8(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("Remark", thrift.STRING, 8); err != nil {
		return fmt.Errorf("%T write field begin error 8:Remark: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Remark)); err != nil {
		return fmt.Errorf("%T.Remark (8) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 8:Remark: %s", p, err)
	}
	return err
}

func (p *User) readField9(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBool(); err != nil {
		return fmt.Errorf("error reading field 9: %s", err)
	} else {
		p.IsAdmin = v
	}
	return nil
}

func (p *User) writeField9(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("IsAdmin", thrift.BOOL, 9); err != nil {
		return fmt.Errorf("%T write field begin error 9:IsAdmin: %s", p, err)
	}
	if err := oprot.WriteBool(bool(p.IsAdmin)); err != nil {
		return fmt.Errorf("%T.IsAdmin (9) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 9:IsAdmin: %s", p, err)
	}
	return err
}

func (p *User) readField10(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 10: %s", err)
	} else {
		p.UserGroupID = v
	}
	return nil
}

func (p *User) writeField10(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("UserGroupID", thrift.STRING, 10); err != nil {
		return fmt.Errorf("%T write field begin error 10:UserGroupID: %s", p, err)
	}
	if err := oprot.WriteString(string(p.UserGroupID)); err != nil {
		return fmt.Errorf("%T.UserGroupID (10) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 10:UserGroupID: %s", p, err)
	}
	return err
}

func (p *User) readField11(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return fmt.Errorf("error reading field 11: %s", err)
	} else {
		p.Status = v
	}
	return nil
}

func (p *User) writeField11(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("Status", thrift.I32, 11); err != nil {
		return fmt.Errorf("%T write field begin error 11:Status: %s", p, err)
	}
	if err := oprot.WriteI32(int32(p.Status)); err != nil {
		return fmt.Errorf("%T.Status (11) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 11:Status: %s", p, err)
	}
	return err
}

func (p *User) readField12(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 12: %s", err)
	} else {
		p.PubInfoID = v
	}
	return nil
}

func (p *User) writeField12(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("PubInfoID", thrift.STRING, 12); err != nil {
		return fmt.Errorf("%T write field begin error 12:PubInfoID: %s", p, err)
	}
	if err := oprot.WriteString(string(p.PubInfoID)); err != nil {
		return fmt.Errorf("%T.PubInfoID (12) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 12:PubInfoID: %s", p, err)
	}
	return err
}

func (p *User) readField13(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 13: %s", err)
	} else {
		p.OrganizationID = v
	}
	return nil
}

func (p *User) writeField13(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("OrganizationID", thrift.STRING, 13); err != nil {
		return fmt.Errorf("%T write field begin error 13:OrganizationID: %s", p, err)
	}
	if err := oprot.WriteString(string(p.OrganizationID)); err != nil {
		return fmt.Errorf("%T.OrganizationID (13) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 13:OrganizationID: %s", p, err)
	}
	return err
}

func (p *User) readField14(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 14: %s", err)
	} else {
		p.Gender = v
	}
	return nil
}

func (p *User) writeField14(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("Gender", thrift.STRING, 14); err != nil {
		return fmt.Errorf("%T write field begin error 14:Gender: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Gender)); err != nil {
		return fmt.Errorf("%T.Gender (14) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 14:Gender: %s", p, err)
	}
	return err
}

func (p *User) readField15(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return fmt.Errorf("error reading field 15: %s", err)
	} else {
		p.Index = v
	}
	return nil
}

func (p *User) writeField15(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("Index", thrift.I32, 15); err != nil {
		return fmt.Errorf("%T write field begin error 15:Index: %s", p, err)
	}
	if err := oprot.WriteI32(int32(p.Index)); err != nil {
		return fmt.Errorf("%T.Index (15) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 15:Index: %s", p, err)
	}
	return err
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
	if err := p.writeField14(oprot); err != nil {
		return err
	}
	if err := p.writeField15(oprot); err != nil {
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
	session, col := UserCol()
	defer session.Close()

	if o.ID == "" {
		o.ID = bson.NewObjectId()
	}
	return col.UpsertId(o.ID, o)
}

func (o *User) Sync() (err error) {
	session, col := UserCol()
	defer session.Close()

	_, err = col.Find(o).Apply(mgo.Change{
		Update:    o,
		Upsert:    true,
		ReturnNew: true,
	}, o)
	return
}

func UserCol() (session *mgo.Session, col *mgo.Collection) {
	return db.GetCol(UserTableName)
}

//Form methods

func (w *User) initWidget() {
	w.widgets = make(map[string]*Widget, 15)
}

func (o *User) ReadForm(params map[string]string) (hasError bool) {
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
	if val, ok := params["Gender"]; ok {
		o.Gender = val
	}
	if val, ok := params["Index"]; ok {
		intVal, _ := strconv.Atoi(val)
		o.Index = int32(intVal)
	}
	return o.ValidateData()
}

func (o *User) ValidateData() (hasError bool) {
	if o.UserName == "" {
		o.UserNameWidget().ErrorMsg = "请输入用户名"
		hasError = true
	}

	if o.Password == "" {
		o.PasswordWidget().ErrorMsg = "请输入密码"
		hasError = true
	}

	if o.Email == "" {
		hasError = true
		o.EmailWidget().ErrorMsg = "请输入电邮"
	} else {
		if !regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).Match([]byte(o.Email)) {
			hasError = true
			o.EmailWidget().ErrorMsg = "电邮格式不正确"
		}
	}
	return
}

func (o *User) IDWidget() *Widget {
	name := "ID"
	ret, ok := o.widgets[name]
	if !ok || ret == nil {
		ret = &Widget{
			Value: o.Id(),
		}
		if o.widgets == nil {
			o.initWidget()
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) UserNameWidget() *Widget {
	name := "UserName"
	ret, ok := o.widgets[name]
	if !ok || ret == nil {
		ret = &Widget{
			Label:       "用户名",
			Value:       o.UserName,
			Name:        "UserName",
			PlaceHolder: "",
			Type:        "text",
			Required:    true,
		}
		if o.widgets == nil {
			o.initWidget()
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) PasswordWidget() *Widget {
	name := "Password"
	ret, ok := o.widgets[name]
	if !ok || ret == nil {
		ret = &Widget{
			Label:       "密码",
			Value:       o.Password,
			Name:        "Password",
			PlaceHolder: "",
			Type:        "password",
			Required:    true,
		}
		if o.widgets == nil {
			o.initWidget()
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) NameWidget() *Widget {
	name := "Name"
	ret, ok := o.widgets[name]
	if !ok || ret == nil {
		ret = &Widget{
			Label:       "姓名",
			Value:       o.Name,
			Name:        "Name",
			PlaceHolder: "",
			Type:        "text",
		}
		if o.widgets == nil {
			o.initWidget()
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) EmailWidget() *Widget {
	name := "Email"
	ret, ok := o.widgets[name]
	if !ok || ret == nil {
		ret = &Widget{
			Label:       "电邮",
			Value:       o.Email,
			Name:        "Email",
			PlaceHolder: "",
			Type:        "text",
			Required:    true,
		}
		if o.widgets == nil {
			o.initWidget()
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) IntroWidget() *Widget {
	name := "Intro"
	ret, ok := o.widgets[name]
	if !ok || ret == nil {
		ret = &Widget{
			Label:       "介绍",
			Value:       o.Intro,
			Name:        "Intro",
			PlaceHolder: "",
			Type:        "text",
		}
		if o.widgets == nil {
			o.initWidget()
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) PictureWidget() *Widget {
	name := "Picture"
	ret, ok := o.widgets[name]
	if !ok || ret == nil {
		ret = &Widget{
			Label:       "Picture",
			Value:       o.Picture,
			Name:        "Picture",
			PlaceHolder: "",
			Type:        "text",
		}
		if o.widgets == nil {
			o.initWidget()
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) RemarkWidget() *Widget {
	name := "Remark"
	ret, ok := o.widgets[name]
	if !ok || ret == nil {
		ret = &Widget{
			Label:       "Remark",
			Value:       o.Remark,
			Name:        "Remark",
			PlaceHolder: "",
			Type:        "text",
		}
		if o.widgets == nil {
			o.initWidget()
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) IsAdminWidget() *Widget {
	name := "IsAdmin"
	ret, ok := o.widgets[name]
	if !ok || ret == nil {
		ret = &Widget{
			Label:       "管理员",
			Value:       strconv.FormatBool(o.IsAdmin),
			Name:        "IsAdmin",
			PlaceHolder: "",
			Type:        "checkbox",
		}
		if o.widgets == nil {
			o.initWidget()
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) UserGroupIDWidget() *Widget {
	name := "UserGroupID"
	ret, ok := o.widgets[name]
	if !ok || ret == nil {
		ret = &Widget{
			Label:       "UserGroupID",
			Value:       o.UserGroupID,
			Name:        "UserGroupID",
			PlaceHolder: "",
			Type:        "select",
			GetBindData: func() (data []*IDLabelPair) {
				objs, _ := UserGroupFindAll(nil)
				length := len(objs)
				length += 1
				data = make([]*IDLabelPair, 0, length)
				data = append(data, &IDLabelPair{"", ""})
				for _, obj := range objs {
					data = append(data, &IDLabelPair{
						ID:    obj.ID.Hex(),
						Label: obj.Name,
					})
				}

				return
			},
		}
		if o.widgets == nil {
			o.initWidget()
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) StatusWidget() *Widget {
	name := "Status"
	ret, ok := o.widgets[name]
	if !ok || ret == nil {
		ret = &Widget{
			Label:       "状态",
			Value:       strconv.FormatInt(int64(o.Status), 10),
			Name:        "Status",
			PlaceHolder: "",
			Type:        "select",
			EnumKey:     test.UserStatusKey,
			EnumData:    test.UserStatusLabel,
		}
		if o.widgets == nil {
			o.initWidget()
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) PubInfoIDWidget() *Widget {
	name := "PubInfoID"
	ret, ok := o.widgets[name]
	if !ok || ret == nil {
		ret = &Widget{
			Label:       "PubInfoID",
			Value:       o.PubInfoID,
			Name:        "PubInfoID",
			PlaceHolder: "",
			Type:        "text",
		}
		if o.widgets == nil {
			o.initWidget()
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) OrganizationIDWidget() *Widget {
	name := "OrganizationID"
	ret, ok := o.widgets[name]
	if !ok || ret == nil {
		ret = &Widget{
			Label:       "OrganizationID",
			Value:       o.OrganizationID,
			Name:        "OrganizationID",
			PlaceHolder: "",
			Type:        "text",
		}
		if o.widgets == nil {
			o.initWidget()
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) GenderWidget() *Widget {
	name := "Gender"
	ret, ok := o.widgets[name]
	if !ok || ret == nil {
		ret = &Widget{
			Label:       "部门",
			Value:       o.Gender,
			Name:        "Gender",
			PlaceHolder: "",
			Type:        "selectPk",
			StringList:  test.Genders,
		}
		if o.widgets == nil {
			o.initWidget()
		}
		o.widgets[name] = ret
	}

	return ret
}
func (o *User) IndexWidget() *Widget {
	name := "Index"
	ret, ok := o.widgets[name]
	if !ok || ret == nil {
		ret = &Widget{
			Label:       "Index",
			Value:       strconv.FormatInt(int64(o.Index), 10),
			Name:        "Index",
			PlaceHolder: "",
			Type:        "text",
		}
		if o.widgets == nil {
			o.initWidget()
		}
		o.widgets[name] = ret
	}

	return ret
}

func (o *User) GetListedLabels() []*IDLabelPair {
	return []*IDLabelPair{}
}

func (o *User) Id() string {
	return o.ID.Hex()
}

func (o *User) GetLabel() string {
	return "ID"
}

func (o *User) GetFieldAsString(fieldKey string) (Value string) {
	switch fieldKey {
	case "ID":
		Value = o.ID.Hex()
	case "UserName":
		widget := o.UserNameWidget()
		if widget.Type == "select" {
			Value = widget.String()
		} else {
			Value = o.UserName
		}
	case "Password":
		widget := o.PasswordWidget()
		if widget.Type == "select" {
			Value = widget.String()
		} else {
			Value = o.Password
		}
	case "Name":
		widget := o.NameWidget()
		if widget.Type == "select" {
			Value = widget.String()
		} else {
			Value = o.Name
		}
	case "Email":
		widget := o.EmailWidget()
		if widget.Type == "select" {
			Value = widget.String()
		} else {
			Value = o.Email
		}
	case "Intro":
		widget := o.IntroWidget()
		if widget.Type == "select" {
			Value = widget.String()
		} else {
			Value = o.Intro
		}
	case "Picture":
		widget := o.PictureWidget()
		if widget.Type == "select" {
			Value = widget.String()
		} else {
			Value = o.Picture
		}
	case "Remark":
		widget := o.RemarkWidget()
		if widget.Type == "select" {
			Value = widget.String()
		} else {
			Value = o.Remark
		}
	case "IsAdmin":
		Value = strconv.FormatBool(o.IsAdmin)
	case "UserGroupID":
		widget := o.UserGroupIDWidget()
		if widget.Type == "select" {
			Value = widget.String()
		} else {
			Value = o.UserGroupID
		}
	case "Status":
		widget := o.StatusWidget()
		if widget.Type == "select" {
			idx := int32(o.Status)
			Value = widget.StringList[idx]
		} else {
			Value = strconv.FormatInt(int64(o.Status), 10)
		}
	case "PubInfoID":
		widget := o.PubInfoIDWidget()
		if widget.Type == "select" {
			Value = widget.String()
		} else {
			Value = o.PubInfoID
		}
	case "OrganizationID":
		widget := o.OrganizationIDWidget()
		if widget.Type == "select" {
			Value = widget.String()
		} else {
			Value = o.OrganizationID
		}
	case "Gender":
		widget := o.GenderWidget()
		if widget.Type == "select" {
			Value = widget.String()
		} else {
			Value = o.Gender
		}
	case "Index":
		widget := o.IndexWidget()
		if widget.Type == "select" {
			idx := int32(o.Index)
			Value = widget.StringList[idx]
		} else {
			Value = strconv.FormatInt(int64(o.Index), 10)
		}
	}
	return
}

type UserWidget struct {
	ID             *Widget
	UserName       *Widget
	Password       *Widget
	Name           *Widget
	Email          *Widget
	Intro          *Widget
	Picture        *Widget
	Remark         *Widget
	IsAdmin        *Widget
	UserGroupID    *Widget
	Status         *Widget
	PubInfoID      *Widget
	OrganizationID *Widget
	Gender         *Widget
	Index          *Widget
}

func (o *User) WidgetStruct() *UserWidget {
	return &UserWidget{
		ID:             o.IDWidget(),
		UserName:       o.UserNameWidget(),
		Password:       o.PasswordWidget(),
		Name:           o.NameWidget(),
		Email:          o.EmailWidget(),
		Intro:          o.IntroWidget(),
		Picture:        o.PictureWidget(),
		Remark:         o.RemarkWidget(),
		IsAdmin:        o.IsAdminWidget(),
		UserGroupID:    o.UserGroupIDWidget(),
		Status:         o.StatusWidget(),
		PubInfoID:      o.PubInfoIDWidget(),
		OrganizationID: o.OrganizationIDWidget(),
		Gender:         o.GenderWidget(),
		Index:          o.IndexWidget(),
	}
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
		o.GenderWidget(),
		o.IndexWidget(),
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

	userSort(q, sortFields)

	err = q.One(&result)
	return
}

func userSort(q *mgo.Query, sortFields []string) {
	if len(sortFields) > 0 {
		q.Sort(sortFields...)
		return
	}

	q.Sort("Index", "-_id")
}

func UserFind(query interface{}, limit int, offset int, sortFields ...string) (result []*User, err error) {
	session, col := db.GetCol("User")
	defer session.Close()

	q := col.Find(query).Limit(limit).Skip(offset)

	userSort(q, sortFields)

	err = q.All(&result)
	return
}

func UserFindAll(query interface{}, sortFields ...string) (result []*User, err error) {
	session, col := db.GetCol("User")
	defer session.Close()

	q := col.Find(query)

	userSort(q, sortFields)

	err = q.All(&result)
	return
}

func UserCount(query interface{}) (result int) {
	session, col := db.GetCol("User")
	defer session.Close()

	result, _ = col.Find(query).Count()
	return
}

func UserFindByIDs(id []string) (result []*User, err error) {
	session, col := db.GetCol("User")
	defer session.Close()

	ids := make([]bson.ObjectId, 0, len(id))
	for _, i := range id {
		if bson.IsObjectIdHex(i) {
			ids = append(ids, bson.ObjectIdHex(i))
		}
	}
	err = col.Find(db.M{"_id": db.M{"$in": ids}}).All(&result)
	return
}

func UserToIDList(ms []*User) []string {
	ret := make([]string, len(ms))
	for idx, m := range ms {
		ret[idx] = m.GetFieldAsString("ID")
	}
	return ret
}

func UserToMap(ms []*User) map[string]*User {
	ret := make(map[string]*User, len(ms))
	for _, m := range ms {
		ret[m.Id()] = m
	}
	return ret
}

func UserFindByID(id string) (result *User, err error) {
	session, col := db.GetCol("User")
	defer session.Close()

	if !bson.IsObjectIdHex(id) {
		err = mgo.ErrNotFound
		return
	}
	err = col.FindId(bson.ObjectIdHex(id)).One(&result)
	return
}

func UserRemoveAll(query interface{}) (info *mgo.ChangeInfo, err error) {
	session, col := db.GetCol("User")
	defer session.Close()

	return col.RemoveAll(query)
}

func UserRemoveByID(id string) (result *User, err error) {
	session, col := db.GetCol("User")
	defer session.Close()

	if !bson.IsObjectIdHex(id) {
		err = mgo.ErrNotFound
		return
	}
	err = col.RemoveId(bson.ObjectIdHex(id))
	return
}

//Search

func (o *User) IsSearchEnabled() bool {
	return false
}

//end search

func (o *User) ViewUrl(id string) string {
	return "/admin/user/" + id
}

func (o *User) GetSummary() string {
	return o.Content
}

func (o *User) TypeName() string {
	return UserTableName
}
