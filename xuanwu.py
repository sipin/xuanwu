# -*- coding: utf-8 -*-
import sys

from ptsd import ast
from ptsd.loader import Loader, Parser
from Cheetah.Template import Template
from os import path


if len(sys.argv) != 3:
	print "usage: \n\tpython xuanwu.py thrift_file_path output_folder_path"
	sys.exit()

thrift_file = sys.argv[1]
out_path = sys.argv[2]
if not out_path.endswith(path.sep):
		out_path = out_path + path.sep


def get_namespace(thrift_file):
	tpl = open(thrift_file, 'r').read()
	tree = Parser().parse(tpl)
	for space in tree.namespaces:
		if space.language_id == "go":
			return space.name

	print 'namespace go not found, please add `namespace go XXXX` to ' + thrift_file + " and retry"
	sys.exit(1)

namespace = get_namespace(thrift_file)

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
	types = {
		"string": "string",
		"bool": "bool",
		"i32": "int32"
	}
	for field in obj.fields:
		field.type = str(field.type)
		field.go_type = types[str(field.type)]

	tpl = open('go.tmpl', 'r').read()
	t = Template(tpl, searchList=[{"namespace": namespace, "obj": obj}])
	code = unicode(t)
	f = open(out_path + 'gen_' + obj.name.value.lower() + ".go", "w")
	f.write(code)
	f.close()


def transform(module):
	for node in module.values():
		if not isinstance(node, ast.Node):
			continue
		if isinstance(node, ast.Enum):
			transform_enum(node)
		elif isinstance(node, ast.Struct):
			transform_struct(node)
		# TODO(constants)
		# TODO(typedefs)


def main(thrift_idl):
	loader = Loader(thrift_idl, lambda x: x)

	tpl = open('go_package.tmpl', 'r').read()
	t = Template(tpl, searchList=[{"namespace": namespace}])
	code = unicode(t)
	f = open(out_path + 'gen_init.go', "w")
	f.write(code)
	f.close()

	for module in loader.modules.values():
		transform(module)

main(thrift_file)
