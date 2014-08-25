#!/usr/bin/env python
# -*- coding: utf-8 -*-

import sys
import os
import base
import fcntl

from ptsd.loader import Loader
from Cheetah.Template import Template
from os import path

if len(sys.argv) != 3:
	print "usage: \n\tpython crud.py thrift_file_path output_folder_path"
	sys.exit()

thrift_file = sys.argv[1]
out_path = sys.argv[2]
filename = ".".join(path.basename(thrift_file).split(".")[:-1])

if not out_path.endswith(path.sep):
	out_path = out_path + path.sep

outDir = out_path.split(path.sep)[-2]

modelsNamespace = ""

def updateController(out_path):
	subdirs = [o for o in os.listdir(out_path) if os.path.isdir(os.path.join(out_path, o))]
	content = "package controller\n\nimport (\n"
	for x in subdirs:
		l = "\t_ \"controller/{0}\"\n".format(x)
		content = content + l

	content = content + ")\n"
	f = open(out_path + "gen_init.go", "w")
	fcntl.flock(f, fcntl.LOCK_EX)
	f.write(content)
	fcntl.flock(f, fcntl.LOCK_UN)
	f.close()

def getControlDir(urlBase):
	outdir = out_path + urlBase.split(path.sep)[-2] + path.sep
	if not os.path.exists(outdir):
		os.makedirs(outdir)
	return outdir

def fieldElem(field, key):
	for att in field.annotations:
		if att.name.value.lower() == key:
			return att.value.value
	return ""

def fieldElems(field, key):
	key = key.lower()
	ret = []
	for att in field.annotations:
		if att.name.value.lower() == key:
			ret.append(att.value.value)
	return ret

def transform_tpl(obj, name):
	idField = obj.fields[0]
	tpl = fieldElem(idField, name)
	namespace = tpl

	if tpl == "axure":
		return tpl

	if tpl != "":
		name1 = ""
		for item in obj.imports:
			if isinstance(item, tuple):
				if item[1] == namespace:
					return item[0]
		obj.imports.add((name, namespace))
		if tpl == "axure":
			return tpl
		return name
	return None

def transform_tpls(obj):
	obj.edittpl = transform_tpl(obj, "edittpl")
	obj.viewtpl = transform_tpl(obj, "viewtpl")
	obj.indextpl = transform_tpl(obj, "indextpl")
	obj.createtpl = transform_tpl(obj, "createtpl")

def assure_path_exists(path):
    dir = os.path.dirname(path)
    if not os.path.exists(dir):
            os.makedirs(dir)

def gen_axure(obj):
	return
	with open("../axure/" + obj.baseURL[1:].replace("/", "_") + ".txt", "r") as f:
		axure = f.read()

	def get_field_by_label(obj, label):
		for f in obj.fields:
			if f.label == label:
				return f
		raise Exception(thrift_file + " " + obj.name.value + " has no label: " + label)

	rows = []
	max_fields = 1
	for row in axure.split("\n"):
		fields = []
		labels = row.split("\t")
		if len(labels) > max_fields:
			max_fields = len(labels)
		for label in labels:
			fields.append(get_field_by_label(obj, label))
		rows.append(fields)

	crud = open('tmpl/axure_create.tmpl', 'r').read().decode("utf8")
	obj.max_fields = max_fields
	res = Template(crud, searchList=[{"namespace": outDir,
									"className": obj.name.value,
									"obj": obj,
									"rows": rows,
									}])

	outfile = "../tpl" + obj.baseURL + "/create.gohtml"
	assure_path_exists(outfile)
	with open(outfile, "w+") as fp:
		fp.write(str(res))

	print outfile

def transform_module(module):
	for obj in module.structs:
		urlBase = ""
		obj.imports = set()
		idField = obj.fields[0]

		urlBase = fieldElem(idField, "baseurl")

		transform_tpls(obj)

		tplPackage = fieldElem(idField, "tplpackage")
		if tplPackage == "":
			tplPackage = "tpl/auto"

		if urlBase == "":
			continue

		if obj.label != obj.name.value:
			obj.perm = obj.label
			if hasattr(idField, "perm"):
				obj.perm = idField.perm

			obj.imports.add("admin/permission")
			obj.hasUser = len([i for i in obj.fields if str(i.name) == "UsersID"]) > 0

			if len(obj.relateObj) > 0:
				obj.imports.add("encoding/json")

			outDir = urlBase.split(path.sep)[-2]
			crud = open('tmpl/crud.tmpl', 'r').read().decode("utf8")

			res = Template(crud, searchList=[{"namespace": outDir,
											"className": obj.name.value,
											"urlBase": urlBase,
											"tplPackage": tplPackage,
											"modelsNamespace": modelsNamespace,
											"obj": obj,
											}])
			writeDir = getControlDir(urlBase)
			outfile = writeDir + "gen_" + obj.name.value.lower() + ".go"
			with open(outfile, "w") as fp:
				fp.write(str(res))

			if obj.createtpl == "axure":
				gen_axure(obj)
		else:
			for field in obj.fields:
				if field.widget_type in ("relateSelect", "relateAjaxSelect") and (not hasattr(field, "bindFunc")):
					outDir = urlBase.split(path.sep)[-2]
					tmpl = file("tmpl/relateOnly.tmpl").read().decode("u8")
					ret = Template(tmpl, searchList=[{"namespace": outDir,
											"className": obj.label,
											"urlBase": urlBase,
											"tplPackage": tplPackage,
											"modelsNamespace": field.bindModels,
											"obj": obj,
											}])
					writeDir = getControlDir(urlBase)
					outfile = writeDir + "gen_relate_" + obj.name.value.lower() + ".go"
					with open(outfile, "w") as fp:
						fp.write(str(ret).strip())
					break


def main(thrift_idl):
	loader = base.load_thrift(thrift_idl)
	global namespace, modelsNamespace
	namespace = loader.namespace
	modelsNamespace = namespace
	for obj in loader.modules.values():
		transform_module(obj)
	updateController(out_path)

main(thrift_file)
