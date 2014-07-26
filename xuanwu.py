# -*- coding: utf-8 -*-
import sys
import base

from ptsd import ast
from ptsd.loader import Loader, Parser
from Cheetah.Template import Template
from os import path, mkdir

if len(sys.argv) != 3:
	print "usage: \n\tpython xuanwu.py thrift_file_path output_folder_path"
	sys.exit()

namespace = ""
thrift_file = sys.argv[1]
out_path = sys.argv[2]
filename = ".".join(path.basename(thrift_file).split(".")[:-1])

if not out_path.endswith(path.sep):
	out_path = out_path + path.sep

try:
	src_path = out_path.replace("\\", "/")
	src_path = src_path[src_path.index("/src/")+5:].strip("/")
except ValueError:
	print "output_folder_path should contains '/src/', for xuanwu to use absolute go path import"
	sys.exit()

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

def struct_import(obj):
	idField = obj.fields[0]
	obj.imports = set(["bytes", "fmt"])
	obj.search = get_search(obj)
	if obj.search != None:
		obj.imports.add("github.com/mattbaird/elastigo/core")


	for f in obj.fields:
		if f.foreign_package != "":
			obj.imports.add(f.foreign_package)

	obj.label = obj.name.value
	if hasattr(idField, "label"):
		obj.label = idField.label

	if len([f for f in obj.filterFields if f.type == "string"]) > 0:
		obj.imports.add("github.com/mattbaird/elastigo/indices")

	for field in obj.fields:
		if hasattr(field, "rule"):
			obj.imports.add("regexp")

		if hasattr(field, "stringList") or hasattr(field, "enums"):
			import_module = filename
			if src_path:
				import_module = src_path + "/" + import_module
			obj.imports.add(import_module)

		if field.type in ["i32", "i64", "bool"] and field.widget_type not in ["date", "time", "datetime"]:
			obj.imports.add("strconv")

		if field.type == "list<string>":
			obj.imports.add("strings")

		if field.widget_type in ["date", "time", "datetime"]:
			obj.imports.add("time")

		if hasattr(field, "meta"):
			obj.imports.add("encoding/json")
	obj.imports = sorted(set(obj.imports))


def transform_struct(obj):
	struct_import(obj)
	tpl = open('tmpl/go.tmpl', 'r').read()
	t = Template(tpl, searchList=[{"namespace": namespace, "filename": filename, "obj": obj}])
	code = str(t)
	with open(out_path + 'gen_' + obj.name.value.lower() + ".go", "w") as f:
		f.write(code)

def transform(module):
	for struct in module.structs:
		transform_struct(struct)

	if len(module.consts) > 0:
		tpl = open('tmpl/go_const.tmpl', 'r').read()
		t = Template(tpl, searchList=[{"namespace": filename, "objs": module.consts}])
		if not path.exists(out_path + filename):
			mkdir(out_path + filename)
		with open(out_path + "%s/gen_%s_const.go" % (filename, filename), "w") as fp:
			fp.write(str(t))

	if len(module.enums) > 0:
		tpl = open('tmpl/go_enum.tmpl', 'r').read()
		t = Template(tpl, searchList=[{
			"namespace": filename,
			"objs": module.enums,
			"name": filename,
		}])
		if not path.exists(out_path + filename):
			mkdir(out_path + filename)
		with open(out_path + "%s/gen_%s_enum.go" % (filename, filename), "w") as fp:
			fp.write(str(t))


def main(thrift_idl):
	loader = base.load_thrift(thrift_idl)
	global namespace
	namespace = loader.namespace
	tpl = open('tmpl/go_package.tmpl', 'r').read()
	t = Template(tpl, searchList=[{"namespace": namespace}])
	code = unicode(t)
	with open(out_path + 'gen_init.go', "w") as fp:
		fp.write(code)

	for module in loader.modules.values():
		transform(module)

main(thrift_file)
