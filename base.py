# -*- coding: utf-8 -*-
import sys

from ptsd import ast
from ptsd.loader import Loader, Parser
from Cheetah.Template import Template
from os import path, mkdir
thrift_file = ""
type_ref = dict(
	string = "string",
	i32 = "int32",
	i64 = "int64",
	bool = "bool",
)

widget_types = set([
	"checkbox",
	"combobox",
	"date",
	"datetime",
	"droplist",
	"file",
	"files",
	"hidden",
	"opinion",
	"password",
	"photo",
	"photos",
	"radio",
	"relateSelect",
	"richtext",
	"select",
	"selectPk",
	"text",
	"textarea",
	"time",
	"userselect",
])

supported_annotations = set([
	"baseURL",
	"bindData",
	"disabled",
	"dm",
	"enums",
	"filterFields",
	"fk",
	"index",
	"label",
	"listedFields",
	"meta",
	"orderFields",
	"placeholder",
	"readonly",
	"relateData",
	"requiredMsg",
	"rule",
	"ruleMsg",
	"search",
	"stringList",
	"summary",
	"toList",
	"tplPackage",
	"viewUrl",
	"widget",

	#permission
	"Create", "Read", "Update", "Delete",
])

typedef = dict()

def type_translate(obj):
	if str(obj) in type_ref:
		return type_ref[str(obj)]
	if isinstance(obj, ast.List):
		return "[]%s" % type_translate(obj.value_type)
	if isinstance(obj, ast.Map):
		return "map[%s]%s" % (type_translate(obj.key_type), type_translate(obj.value_type))
	if str(obj.value) in typedef:
		return obj.value
	return "unknown(%s)" % obj

def add_properties(field, obj):
	field.label = field.name.value
	# todo: add field name checking
	for att in field.annotations:
		if att.name.value not in supported_annotations:
			raise Exception(thrift_file + " " + obj.name.value + " " + field.name.value +
			 " has invalid field annotation: " + att.name.value)
		setattr(field, att.name.value, att.value.value)

	if hasattr(field, "bindData"):
		col, label = field.bindData.split(".")
		tpl = open('tmpl/field_getBindData.tmpl', 'r').read()

		t = Template(tpl, searchList=[{"field": field, "col": col, "label": label}])
		field.bindData = str(t).strip()
		field.bindTable = col

	if hasattr(field, "meta"):
		if str(field.type) != "string":
			raise Exception(thrift_file + " " + obj.name.value +
				"meta data must should store in string type")

		tpl = open('tmpl/field_getMeta.tmpl', 'r').read()

		t = Template(tpl, searchList=[{"field": field, "table": field.meta}])
		field.metaFunc = str(t).strip()

	if not hasattr(field, "placeholder"):
		field.placeholder = ""

	if field.required and not hasattr(field, "requiredMsg"):
		field.requiredMsg = "请输入" + field.label

def get_widget_type(obj, field):
	for att in field.annotations:
		if att.name.value.lower() == "widget":
			if not att.value.value in widget_types:
				raise Exception(thrift_file + " " + obj.name.value + ":" +
					field.name.value + " has invalid widget type: " + att.value.value)
			return att.value.value
		if att.name.value.lower() == "dm":
			field.placeholder = att.value.value
			return "dm"
		if att.name.value.lower() == "meta":
			return "meta"
	return "text"

def transform_type(field_type, obj):
	if isinstance(field_type, (ast.Byte, ast.I16, ast.I32, ast.I64)):
		return 'Integer'
	elif isinstance(field_type, ast.Double):
		return 'Float'
	elif isinstance(field_type, ast.Bool):
		return 'Boolean'
	elif isinstance(field_type, (ast.Binary, ast.String)):
		return 'String'
	elif isinstance(field_type, ast.Identifier):
		return field_type.value
	elif isinstance(field_type, ast.Map):
		return 'Map(%s, %s)' % (
				transform_type(field_type.key_type), transform_type(field_type.value_type))
	elif isinstance(field_type, ast.List):
		return 'List(%s)' % transform_type(field_type.value_type)
	elif isinstance(field_type, ast.Set):
		# TODO(wickman) Support pystachio set?
		return 'List(%s)' % transform_type(field_type.value_type)
	raise ValueError(thrift_file + 'Unsupported conversion type: %s (base:%s)' % (field_type, type(field_type)))

def transform_field(field, indent=0):
	line = ' ' * indent
	line += '%s = ' % field.name.value
	if not field.required and not field.const_value:
		line += transform_type(field.type)
	elif field.required and not field.const_value:
		line += 'Required(%s)' % transform_type(field.type)
	elif not field.required and field.const_value:
		line += 'Default(%s, %s)' % (
				transform_type(field.type), ast.Const.render_value(field.const_value))
	else:
		raise ValueError('Cannot have required field with default values.')
	return line

class ListField:
	pass

def transform_const(obj):
	obj.go_type = type_translate(obj.type)
	obj.is_map = isinstance(obj.type, ast.Map)
	obj.is_list = isinstance(obj.type, ast.List)
	return obj

def init_Fields(obj):
	idField = obj.fields[0]
	obj.label = obj.name.value
	obj.relateObj = {}
	obj.fieldMap = {}
	for field in obj.fields:
		add_properties(field, obj)
		field.foreign = ""
		field.foreign_package = ""
		if field.name.value.endswith("ID"):
			field.foreign = field.name.value[:-2]
			field.foreign_type = field.foreign

		if hasattr(field, "fk"):
			if "." in field.fk:
				pos = field.fk.rindex(".")
				field.foreign = field.fk[pos+1:]
				field.foreign_package = field.fk[0:pos]

				if "/" in field.foreign_package:
					pos = field.foreign_package.rindex("/")
					field.foreign_type = field.foreign_package[pos+1:] + "." + field.foreign
				else:
					field.foreign_type = field.foreign_package + "." + field.foreign
			else:
				field.foreign_type = field.fk

		field.go_type = type_translate(field.type)
		field.type = str(field.type)
		field.widget_type = get_widget_type(obj, field)
		obj.fieldMap[field.name.value] = field

		if field.widget_type == "relateSelect":
			obj.relateObj[field.name.value] = field
			field.relateFields = []

	for field in obj.fields:
		if hasattr(field, "relateData"):
			col, label = field.relateData.split(".")
			if col not in obj.relateObj:
				raise Exception(thrift_file + " missing realte field " + col)
			obj.relateObj[col].relateFields.append((field.name.value, label))
			field.disabled = "True"
			delattr(field, "relateData")

	obj.toList = [i.name.value for i in obj.fields if hasattr(i, "toList")]
	if "ID" not in obj.toList:
		obj.toList.append("ID")

def init_ListedField(obj):
	idField = obj.fields[0]
	obj.listedFields = []
	obj.listedFieldStrings = []
	obj.listedFieldNames = []
	if hasattr(idField, "listedFields"):
		listedFields = [f.strip() for f in idField.listedFields.split(",")]
		obj.listedFieldNames = listedFields
		for fieldname in listedFields:
			if "." in fieldname:
				objName, objField = fieldname.split(".")
				for field in obj.fields:
					if field.name.value.endswith("ID") and field.name.value[:-2] == objName:
						f = ListField()
						f.label = field.label
						f.key = fieldname
						obj.listedFields.append(f)
						fs = ListField()
						fs.key = fieldname
						fs.objName = objName
						fs.objField = objField
						obj.listedFieldStrings.append(fs)
			else:
				for field in obj.fields:
					if field.name.value == fieldname:
						f = ListField()
						f.label = field.label
						f.key = fieldname
						obj.listedFields.append(f)

		if len(listedFields) > len(obj.listedFields):
			foundFields = [field.name.value for field in obj.listedFields if field.name.value in listedFields]
			missingFields = [field for field in listedFields if field not in foundFields]
			raise Exception(thrift_file + " " + obj.name.value +
				" missing listedFields: " + str(missingFields))

def init_OrderFields(obj):
	idField = obj.fields[0]
	if hasattr(idField, "orderFields"):
		orderFields = {}
		for field in map(lambda s:s.strip(), idField.orderFields.split(",")):
			if field.find(":") >= 0:
				field, order = field.split(":")[:2]
			else:
				order = ""
			orderFields[field] = order if order in ("asc", "desc") else "none"
		for field in obj.listedFields:
			if field.key in orderFields:
				field.order = orderFields[field.key]
				del orderFields[field.key]

		if len(filter(lambda s: hasattr(s, "order") and s.order!="none", obj.listedFields)) > 1:
			raise Exception(thrift_file + " " + obj.name.value +
				" too many default order index")
		if len(orderFields) > 0:
			raise Exception(thrift_file + " " + obj.name.value +
				" missing orderFields: " + ",".join(orderFields))

def init_FilterFields(obj):
	idField = obj.fields[0]
	obj.filterFields = []
	obj.termKeys = []
	obj.dateKeys = []
	if hasattr(idField, "filterFields"):
		filterFields = [f.strip() for f in idField.filterFields.split(",")]
		for fieldName in filterFields:
			for field in obj.fields:
				if field.name.value == fieldName:
					obj.filterFields.append(field)
					if field.widget_type in ["date", "time", "datetime"]:
						obj.dateKeys.append(fieldName + "Start")
						obj.dateKeys.append(fieldName + "End")
					elif field.type == "string":
						obj.termKeys.append(fieldName)
					elif field.type == "list<string>":
						obj.termKeys.append(fieldName)
					else:
						raise Exception("%s %s has invalid filterField: %s:%s" %
					(thrift_file, obj.name.value, fieldName, field.type))

		if len(filterFields) > len(obj.filterFields):
			foundFields = [field.name.value for field in obj.filterFields if field.name.value in filterFields]
			missingFields = [field for field in filterFields if field not in foundFields]
			raise Exception(thrift_file + " " + obj.name.value +
				" missing filterFields: " + str(missingFields))

def init_module(module):
	module.consts = []
	module.typedef = []
	module.enums = []
	module.structs = []

	for node in [i for i in module.values() if isinstance(i, ast.Typedef)]:
		node.go_type = type_translate(node.type)
		module.typedefs.append(node)
		typedef[str(node.name)] = node.type

	for node in module.values():
		if not isinstance(node, ast.Node):
			continue
		if isinstance(node, ast.Enum):
			module.enums.append(node)
		elif isinstance(node, ast.Const):
			module.consts.append(transform_const(node))
		elif isinstance(node, ast.Struct):
			module.structs.append(node)

	## process enum labels
	for obj in module.enums:
		obj.labels = {}
		for i in obj.values:
			label_anno = [j.value for j in i.annotations if j.name.value == "label"]
			if len(label_anno) == 0:
				obj.labels[i.tag] = '"%s"' % i.name
			else:
				obj.labels[i.tag] = label_anno[0]

	## init struct related
	for obj in module.structs:
		init_Fields(obj)
		init_ListedField(obj)
		init_OrderFields(obj)
		init_FilterFields(obj)

def load_thrift(thrift_idl):
	global thrift_file
	thrift_file = thrift_idl
	loader = Loader(thrift_idl, lambda x: x)
	if loader.namespace == "":
		print 'namespace go not found, please add `namespace go XXXX` to ' + thrift_file + " and retry"
		sys.exit(1)

	loader.namespace = str(loader.namespace)
	for module in loader.modules.values():
		init_module(module)
	return loader
