package models

import (
	//Official libs
	"bytes"
	"fmt"

	//3rd party libs
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"github.com/sipin/gothrift/thrift"

	//Own libs
	"db"
)

type UserGroup struct {
	ID         bson.ObjectId `bson:"_id" thrift:"ID,1"`
	Name  string `bson:"Name" thrift:"Name,2"`
	widgets map[string]*Widget
}


func NewUserGroup() *UserGroup {
	rval := new(UserGroup)
	rval.ID = bson.NewObjectId()
	rval.widgets = make(map[string]*Widget, 2)
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
func (p *UserGroup) readField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 2: %s", err)
	} else {
		p.Name = v
	}
	return nil
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

func (p *UserGroup) writeField1(oprot thrift.TProtocol) (err error) {
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
func (p *UserGroup) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("Name", thrift.STRING, 2); err != nil {
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
	session, col := db.GetCol("UserGroup")
	defer session.Close()

	if o.ID == "" {
		o.ID = bson.NewObjectId()
	}



	return col.UpsertId(o.ID, o)
}

//Form methods

func (o *UserGroup) ReadForm(params map[string]string) {
	if val, ok := params["Name"]; ok {
		o.Name = val
	}
}

func (o *UserGroup) NameWidget() *Widget {
	name := "Name"
	ret, ok := o.widgets[name]
	if !ok || ret==nil {
		ret = &Widget{
			Label: "Name",
			Value : o.Name,
			Name: "Name",
			PlaceHolder: "",
			Type: "string",
		}
		o.widgets[name] = ret
	}

	return ret
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

	if sortFields == nil {
		q.Sort("-_id")
	} else {
		q.Sort(sortFields...)
	}

	err = q.One(&result)
	return
}

func UserGroupFind(query interface{}, limit int, offset int, sortFields ...string) (result []*UserGroup, err error) {
	session, col := db.GetCol("UserGroup")
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

func UserGroupFindAll(query interface{}, sortFields ...string) (result []*UserGroup, err error) {
	session, col := db.GetCol("UserGroup")
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
		err = ErrInvalidObjectId
		return
	}
	err = col.FindId(bson.ObjectIdHex(id)).One(&result)
	return
}

func UserGroupRemoveByID(id string) (result *UserGroup, err error) {
	session, col := db.GetCol("UserGroup")
	defer session.Close()

	err = col.RemoveId(bson.ObjectIdHex(id))
	return
}

// Search
