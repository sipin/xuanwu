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

def add_properties(field):
	field.label = field.name.value
	field.extra = []

	for att in field.annotations:
		if att.name.value.lower() == "label":
			field.label = att.value.value

def get_search(obj):
	search = {}
	
	for field in obj.fields:
		for att in field.annotations:
			if att.name.value.lower() == "search":
				if str(field.type) != "string":
					raise Exception(obj.name.value + ":" + field.name.value + " must be string to be searchable")
				search_key = att.value.value.split("-")[0]
				if search.has_key(search_key):
					search[search_key].append(field.name.value)
				else:
					search[search_key] = [field.name.value]
	return search


def transform_enum(enum):
	name = enum.name.value

	return '%s = Enum(%r, %r)' % (
		 name,
		 name,
		 tuple(value.name.value for value in enum.values))


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



def transform_struct(obj):
	obj.need_strconv = False
	obj.search = get_search(obj)
	obj.need_search = len(obj.search) > 0

	for field in obj.fields:
		add_properties(field)
		field.type = str(field.type)
		field.go_type = type_translate(field.type)

		if field.go_type != "string":
			obj.need_strconv = True

		field.foreign = ""
		if field.name.value.endswith("ID"):
			field.foreign = field.name.value[:-2]

	tpl = open('go.tmpl', 'r').read()
	t = Template(tpl, searchList=[{"namespace": namespace, "obj": obj}])
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

	for node in [i for i in module.values() if isinstance(i, ast.Typedef)]:
		node.go_type = type_translate(node.type)
		type_def.append(node)
		typedef[str(node.name)] = node.type

	for node in module.values():
		if not isinstance(node, ast.Node):
			continue
		if isinstance(node, ast.Enum):
			transform_enum(node)
		elif isinstance(node, ast.Const):
			const.append(transform_const(node))
		elif isinstance(node, ast.Struct):
			transform_struct(node)

	name = ".".join(path.basename(thrift_file).split(".")[:-1])
	if len(const) > 0:
		tpl = open('go_const.tmpl', 'r').read()

		t = Template(tpl, searchList=[{"namespace": name, "objs": const}])
		if not path.exists(out_path + name):
			mkdir(out_path + name)
		f = open(out_path + name + "/gen_const.go", "w")
		f.write(str(t))
		f.close()
	if len(type_def) > 0:
		t = Template(typedef_tpl, searchList=[{"namespace": name, "objs": type_def}])
		if name == "init":
			write_path = "gen_init.go"
		else:
			if not path.exists(out_path + name):
				mkdir(out_path + name)
			write_path = name + "/gen_typedef.go"
		f = open(out_path + write_path, "w")
		f.write(str(t))
		f.close()

def main(thrift_idl):
	global namespace
	loader = Loader(thrift_idl, lambda x: x)

	if loader.namespace == "":
		print 'namespace go not found, please add `namespace go XXXX` to ' + thrift_file + " and retry"
		sys.exit(1)

	namespace = loader.namespace

	tpl = open('go_package.tmpl', 'r').read()
	t = Template(tpl, searchList=[{"namespace": namespace}])
	code = unicode(t)
	f = open(out_path + 'gen_init.go', "w")
	f.write(code)
	f.close()

	for module in loader.modules.values():
		transform(module)

main(thrift_file)
