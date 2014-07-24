# -*- coding: utf-8 -*-
import sys

from ptsd import ast
from ptsd.loader import Loader, Parser
from Cheetah.Template import Template
from os import path, mkdir

type_ref = dict(
	string = "string",
	i32 = "int32",
	i64 = "int64",
	bool = "bool",
)

widget_types = set([
	"text",
	"textarea",
	"opinion",
	"richtext",
	"password",
	"radio",
	"checkbox",
	"select",
	"selectPk",
	"date",
	"datetime",
	"photo",
	"photos",
	"file",
	"files",
	"droplist",
	"combobox",
	"hidden",
	"time",
	"userselect",
	"relateSelect"
])

supported_annotations = set([
	"label",
	"baseURL",
	"listedFields",
	"orderFields",
	"search",
	"dm",
	"meta",
	"widget",
	"bindData",
	"relateData",
	"requiredMsg",
	"rule",
	"ruleMsg",
	"filterFields",
	"placeholder",
	"stringList",
	"readonly",
	"disabled",
	"index",
	"enums",
	"toList",
	"summary",
	"viewUrl",
	"tplPackage",

	#permission
	"Create", "Read", "Update", "Delete",
])

typedef = dict()
typedef_tpl = '''
package $namespace

#for v in $objs
type $v.name $v.go_type
#end for
'''.strip()

if len(sys.argv) != 3:
	print "usage: \n\tpython xuanwu.py thrift_file_path output_folder_path"
	sys.exit()

namespace = ""
thrift_file = sys.argv[1]
out_path = sys.argv[2]
filename = ".".join(path.basename(thrift_file).split(".")[:-1])

try:
	src_path = out_path.replace("\\", "/")
	src_path = src_path[src_path.index("/src/")+5:].strip("/")
except ValueError:
	print "output_folder_path should contains '/src/', for xuanwu to use absolute go path import"
	sys.exit()


if not out_path.endswith(path.sep):
		out_path = out_path + path.sep

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
			raise Exception(obj.name.value + " " + field.name.value +
			 " has invalid field annotation: " + att.name.value)
		setattr(field, att.name.value, att.value.value)

	if hasattr(field, "widget") and field.widget == "relateSelect":
		tpl = open('tmpl/field_getRelateFields.tmpl', 'r').read()
		t = Template(tpl, searchList=[{"relateFields": field.relateFields}])
		field.relateData = str(t).strip()

	if hasattr(field, "bindData"):
		col, label = field.bindData.split(".")
		tpl = open('tmpl/field_getBindData.tmpl', 'r').read()

		t = Template(tpl, searchList=[{"field": field, "col": col, "label": label}])
		field.bindData = str(t).strip()

	if hasattr(field, "meta"):
		if str(field.type) != "string":
			raise Exception("meta data must should store in string type")

		tpl = open('tmpl/field_getMeta.tmpl', 'r').read()

		t = Template(tpl, searchList=[{"field": field, "table": field.meta}])
		field.metaFunc = str(t).strip()

	if not hasattr(field, "placeholder"):
		field.placeholder = ""

	if field.required and not hasattr(field, "requiredMsg"):
		field.requiredMsg = "请输入" + field.label

def get_search(obj):
	search = None

	for field in obj.fields:
		if field.name.value == "ID" and hasattr(field, "search"):
			search = [f.strip() for f in field.search.split(",")]

	if search == None:
		return search

	for fieldName in search:
		try:
			field = obj.fieldMap[fieldName]
			if str(field.type) != "string":
				raise Exception(obj.name.value + " has non-string searchField: " + fieldName)
		except KeyError:
			raise Exception(obj.name.value + " has invalid searchField: " + fieldName)

	return search

def get_widget_type(field):
	for att in field.annotations:
		if att.name.value.lower() == "widget":
			if not att.value.value in widget_types:
				raise Exception(field.name.value + " has invalid widget type: " + att.value.value)
			return att.value.value
		if att.name.value.lower() == "dm":
			field.placeholder = att.value.value
			return "dm"
		if att.name.value.lower() == "meta":
			return "meta"
	return "text"

def transform_type(field_type):
	if isinstance(field_type, (ast.Byte, ast.I16, ast.I32, ast.I64)):
		return 'Integer'
	elif isinstance(field_type, ast.Double):
		return 'Float'
	elif isinstance(field_type, ast.Bool):
		return 'Boolean'
	# ast.Binary != String unfortunately
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
	raise ValueError('Unsupported conversion type: %s (base:%s)' % (field_type, type(field_type)))


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

def transform_struct(obj):
	idField = obj.fields[0]

	obj.relateObj = {}
	obj.fieldMap = {}
	for field in obj.fields:
		field.go_type = type_translate(field.type)
		field.type = str(field.type)
		field.widget_type = get_widget_type(field)
		obj.fieldMap[field.name.value] = field

		if field.widget_type == "relateSelect":
			obj.relateObj[field.name.value] = field
			field.relateFields = []
		else:
			add_properties(field, obj)

	for field in obj.fields:
		if hasattr(field, "relateData"):
			col, label = field.relateData.split(".")
			if col not in obj.relateObj:
				raise Exception(thrift_file + " missing realte field " + col)
			obj.relateObj[col].relateFields.append((field.name.value, label))
			delattr(field, "relateData")

	for field in obj.relateObj.values():
		add_properties(field, obj)

	obj.imports = ["bytes", "fmt"]
	obj.listedFieldStrings = []

	obj.search = get_search(obj)
	if obj.search != None:
		obj.imports.append("github.com/mattbaird/elastigo/core")

	obj.label = obj.name.value
	
	obj.listedFields = []
	if hasattr(idField, "listedFields"):
		listedFields = [f.strip() for f in idField.listedFields.split(",")]
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
			raise Exception(thrift_file + " " + obj.name.value + " missing listedFields: " + str(missingFields))

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
			raise Exception(thrift_file + " " + obj.name.value + " too many default order index")
		if len(orderFields) > 0:
			raise Exception(thrift_file + " " + obj.name.value + " missing orderFields: " + ",".join(orderFields))


	obj.filterFields = []
	obj.termKeys = []
	obj.dateKeys = []
	if hasattr(idField, "filterFields"):
		filterFields = [f.strip() for f in idField.filterFields.split(",")]
		for fieldname in filterFields:
			for field in obj.fields:
				if field.name.value == fieldname:
					obj.filterFields.append(field)
					if field.widget_type in ["date", "time", "datetime"]:
						obj.dateKeys.append(fieldname + "Start")
						obj.dateKeys.append(fieldname + "End")
					elif field.type == "string":
						obj.termKeys.append(fieldname)
					else:
						raise Exception(thrift_file + " " + obj.name.value + " invalid filterFields: " + str(missingFields))			

		if len(filterFields) > len(obj.filterFields):
			foundFields = [field.name.value for field in obj.filterFields if field.name.value in filterFields]
			missingFields = [field for field in filterFields if field not in foundFields]
			raise Exception(thrift_file + " " + obj.name.value + " missing filterFields: " + str(missingFields))

	if hasattr(idField, "label"):
		obj.label = idField.label

	if len([f for f in obj.filterFields if f.type == "string"]) > 0:
		obj.imports.append("github.com/mattbaird/elastigo/indices")

	for field in obj.fields:
		if hasattr(field, "rule"):
			obj.imports.append("regexp")

		if hasattr(field, "stringList") or hasattr(field, "enums"):
			import_module = filename
			if src_path:
				import_module = src_path + "/" + import_module
			obj.imports.append(import_module)

		if field.type in ["i32", "i64", "bool"] and field.widget_type not in ["date", "time", "datetime"]:
			obj.imports.append("strconv")

		if field.type == "list<string>":
			obj.imports.append("strings")

		if field.widget_type in ["date", "time", "datetime"]:
			obj.imports.append("time")

		field.foreign = ""
		if field.name.value.endswith("ID"):
			field.foreign = field.name.value[:-2]

		if hasattr(field, "meta"):
			obj.imports.append("encoding/json")

	obj.toList = [i.name.value for i in obj.fields if hasattr(i, "toList")]
	if "ID" not in obj.toList:
		obj.toList.append("ID")

	obj.imports = sorted(set(obj.imports))
	tpl = open('go.tmpl', 'r').read()
	t = Template(tpl, searchList=[{"namespace": namespace, "filename": filename, "obj": obj}])
	code = str(t)
	f = open(out_path + 'gen_' + obj.name.value.lower() + ".go", "w")
	f.write(code)
	f.close()

def transform_const(obj):
	obj.go_type = type_translate(obj.type)
	obj.is_map = isinstance(obj.type, ast.Map)
	obj.is_list = isinstance(obj.type, ast.List)
	return obj


def transform(module):
	const = []
	type_def = []
	enum = []

	for node in [i for i in module.values() if isinstance(i, ast.Typedef)]:
		node.go_type = type_translate(node.type)
		type_def.append(node)
		typedef[str(node.name)] = node.type

	for node in module.values():
		if not isinstance(node, ast.Node):
			continue
		if isinstance(node, ast.Enum):
			node.labels = {}
			for i in node.values:
				label_anno = [j.value for j in i.annotations if j.name.value == "label"]
				if len(label_anno) == 0:
					node.labels[i.tag] = '"%s"' % i.name
				else:
					node.labels[i.tag] = label_anno[0]
			enum.append(node)
		elif isinstance(node, ast.Const):
			const.append(transform_const(node))
		elif isinstance(node, ast.Struct):
			transform_struct(node)

	if len(const) > 0:
		tpl = open('go_const.tmpl', 'r').read()

		t = Template(tpl, searchList=[{"namespace": filename, "objs": const}])
		if not path.exists(out_path + filename):
			mkdir(out_path + filename)
		f = open(out_path + "%s/gen_%s_const.go" % (filename, filename), "w")
		f.write(str(t))
		f.close()

	if len(enum) > 0:
		tpl = open('go_enum.tmpl', 'r').read()
		t = Template(tpl, searchList=[{
			"namespace": filename,
			"objs": enum,
			"name": filename,
		}])
		if not path.exists(out_path + filename):
			mkdir(out_path + filename)
		f = open(out_path + "%s/gen_%s_enum.go" % (filename, filename), "w")
		f.write(str(t))
		f.close()

	if len(type_def) > 0:
		t = Template(typedef_tpl, searchList=[{"namespace": filename, "objs": type_def}])
		if filename == "init":
			write_path = "gen_init.go"
		else:
			if not path.exists(out_path + filename):
				mkdir(out_path + filename)
			write_path = "%s/gen_%s_typedef.go" % (filename, filename)
		f = open(out_path + write_path, "w")
		f.write(str(t))
		f.close()

def main(thrift_idl):
	global namespace
	loader = Loader(thrift_idl, lambda x: x)

	if loader.namespace == "":
		print 'namespace go not found, please add `namespace go XXXX` to ' + thrift_file + " and retry"
		sys.exit(1)

	namespace = str(loader.namespace)

	tpl = open('go_package.tmpl', 'r').read()
	t = Template(tpl, searchList=[{"namespace": namespace}])
	code = unicode(t)
	f = open(out_path + 'gen_init.go', "w")
	f.write(code)
	f.close()

	for module in loader.modules.values():
		transform(module)

main(thrift_file)
