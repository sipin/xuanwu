# -*- coding: utf-8 -*-
import sys

from ptsd import ast
from ptsd.loader import Loader
from Cheetah.Template import Template

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
		"i32": "int32"
	}
	for field in obj.fields:
		field.type = str(field.type)
		field.go_type = types[str(field.type)]

	tpl = open('go.tmpl', 'r').read()
	t = Template(tpl, searchList=[{"namespace": "models" , "obj": obj}])
	code = unicode(t)
	f = open('../src/kp/models/gen_' + obj.name.value.lower() + ".go", "w")
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
	t = Template(tpl, searchList=[{"namespace": "models"}])
	code = unicode(t)
	f = open('../src/kp/models/gen.go', "w")
	f.write(code)
	f.close()

	for module in loader.modules.values():
		transform(module)


main(sys.argv[1])
