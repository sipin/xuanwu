package models

import (
	//Official libs
	"bytes"
	"fmt"

	//3rd party libs
	"github.com/sipin/gothrift/thrift"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	//Own libs
	"db"
)

func init() {
}

var UserGroupTableName = "UserGroup"

type UserGroup struct {
	ID      bson.ObjectId `bson:"_id,omitempty" thrift:"ID,1"`
	Name    string        `bson:"Name" thrift:"Name,2"`
	widgets map[string]*Widget
}

func NewUserGroup() *UserGroup {
	rval := new(UserGroup)
	rval.ID = bson.NewObjectId()
	return rval
}

func NewUserGroupFromBytes(data []byte) *UserGroup {
	o := new(UserGroup)
	transport := thrift.NewStreamTransportR(bytes.NewBuffer(data))
	protocol := thrift.NewTCompactProtocol(transport)
	o.Read(protocol)
	return o
}

func NewUserGroupWithParams(params map[string]string) *UserGroup {
	o := new(UserGroup)
	o.ReadForm(params)
	return o
}

//Thrift Methods

func (p *UserGroup) Read(iprot thrift.TProtocol) error {
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

func (p *UserGroup) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 1: %s", err)
	} else {
		p.ID = bson.ObjectId(v)
	}
	return nil
}

func (p *UserGroup) writeField1(oprot thrift.TProtocol) (err error) {
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

func (p *UserGroup) readField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 2: %s", err)
	} else {
		p.Name = v
	}
	return nil
}

func (p *UserGroup) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("Name", thrift.STRING, 2); err != nil {
		return fmt.Errorf("%T write field begin error 2:Name: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Name)); err != nil {
		return fmt.Errorf("%T.Name (2) field write error: %s", p, err)
	}

	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 2:Name: %s", p, err)
	}
	return err
}

func (p *UserGroup) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("UserGroup"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
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

func (p *UserGroup) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("UserGroup(%+v)", *p)
}

func (p *UserGroup) ToBytes() []byte {
	transport := thrift.NewTMemoryBuffer()
	protocol := thrift.NewTCompactProtocol(transport)
	p.Write(protocol)
	protocol.Flush()

	return transport.Bytes()
}

//mongo methods

func (o *UserGroup) Save() (info *mgo.ChangeInfo, err error) {
	session, col := UserGroupCol()
	defer session.Close()

	if o.ID == "" {
		o.ID = bson.NewObjectId()
	}

	return col.UpsertId(o.ID, o)
}

func (o *UserGroup) Sync() (err error) {
	session, col := UserGroupCol()
	defer session.Close()

	_, err = col.Find(o).Apply(mgo.Change{
		Update:    o,
		Upsert:    true,
		ReturnNew: true,
	}, o)
	return
}

func UserGroupCol() (session *mgo.Session, col *mgo.Collection) {
	return db.GetCol(UserGroupTableName)
}

//Form methods

func (w *UserGroup) initWidget() {
	w.widgets = make(map[string]*Widget, 2)
}

func (o *UserGroup) ReadForm(params map[string]string) (hasError bool) {
	if val, ok := params["Name"]; ok {
		o.Name = val
	}
	return o.ValidateData()
}

func (o *UserGroup) ValidateData() (hasError bool) {
	return
}

func (o *UserGroup) NameWidget() *Widget {
	name := "Name"
	ret, ok := o.widgets[name]
	if !ok || ret == nil {
		ret = &Widget{
			Label:       "Name",
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

func (o *UserGroup) GetListedLabels() []*IDLabelPair {
	return []*IDLabelPair{}
}

func (o *UserGroup) Id() string {
	return o.ID.Hex()
}

func (o *UserGroup) GetLabel() string {
	return "ID"
}

func (o *UserGroup) GetFieldAsString(fieldKey string) (Value string) {
	switch fieldKey {
	case "ID":
		Value = o.ID.Hex()
	case "Name":
		Value = o.Name
	}
	return
}

type UserGroupWidget struct {
	Name *Widget
}

func (o *UserGroup) WidgetStruct() *UserGroupWidget {
	return &UserGroupWidget{
		Name: o.NameWidget(),
	}
}

func (o *UserGroup) Widgets() []*Widget {
	return []*Widget{
		o.NameWidget(),
	}
}

//foreigh keys

//Collection Manage methods

func UserGroupFindOne(query interface{}, sortFields ...string) (result *UserGroup, err error) {
	session, col := db.GetCol("UserGroup")
	defer session.Close()

	q := col.Find(query)

	usergroupSort(q, sortFields)

	err = q.One(&result)
	return
}

func usergroupSort(q *mgo.Query, sortFields []string) {
	if len(sortFields) > 0 {
		q.Sort(sortFields...)
		return
	}

	q.Sort("-_id")
}

func UserGroupFind(query interface{}, limit int, offset int, sortFields ...string) (result []*UserGroup, err error) {
	session, col := db.GetCol("UserGroup")
	defer session.Close()

	q := col.Find(query).Limit(limit).Skip(offset)

	usergroupSort(q, sortFields)

	err = q.All(&result)
	return
}

func UserGroupFindAll(query interface{}, sortFields ...string) (result []*UserGroup, err error) {
	session, col := db.GetCol("UserGroup")
	defer session.Close()

	q := col.Find(query)

	usergroupSort(q, sortFields)

	err = q.All(&result)
	return
}

func UserGroupCount(query interface{}) (result int) {
	session, col := db.GetCol("UserGroup")
	defer session.Close()

	result, _ = col.Find(query).Count()
	return
}

func UserGroupFindByID(id string) (result *UserGroup, err error) {
	session, col := db.GetCol("UserGroup")
	defer session.Close()

	if !bson.IsObjectIdHex(id) {
		err = mgo.ErrNotFound
		return
	}
	err = col.FindId(bson.ObjectIdHex(id)).One(&result)
	return
}

func UserGroupRemoveAll(query interface{}) (info *mgo.ChangeInfo, err error) {
	session, col := db.GetCol("UserGroup")
	defer session.Close()

	return col.RemoveAll(query)
}

func UserGroupRemoveByID(id string) (result *UserGroup, err error) {
	session, col := db.GetCol("UserGroup")
	defer session.Close()

	if !bson.IsObjectIdHex(id) {
		err = mgo.ErrNotFound
		return
	}
	err = col.RemoveId(bson.ObjectIdHex(id))
	return
}

// Search

func (o *UserGroup) HasSimpleSearch() bool {
	return false
}
