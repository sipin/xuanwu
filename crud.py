# -*- coding: utf-8 -*-
import sys
import os

from ptsd import ast
from ptsd.loader import Loader, Parser
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

def updateController(out_path):
    subdirs = [o for o in os.listdir(out_path) if os.path.isdir(os.path.join(out_path, o))]
    content = "package controller\n\nimport (\n"
    for x in subdirs:
        l = "\t_ \"controller/{0}\"\n".format(x)
        content = content + l

    content = content + ")\n"
    f = open(out_path + "gen_init.go", "w")
    f.write(content)
    f.close()

outDir = out_path.split(path.sep)[-2]

def getControlDir(urlBase):
    outdir = out_path + \
      urlBase.split(path.sep)[-2] + path.sep
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

def permission(name, action, obj):
    validType = ["User", "Role", "Department", "Organization"]
    ret = []
    for p in obj:
        item = []
        pis = p.split(",")
        for pi in pis:
            pos = "(%s.%s=%s)" % (name, action, pi)
            nameType = pi.split(":")
            if len(nameType) != 2:
                raise Exception("%s.%s is invalid: '%s'" % (name, action, pi))
            perm = dict(name=nameType[0], type=nameType[1])
            if perm["type"] not in validType:
                raise Exception("type not in %s %s" % (validType, pos))
            if perm["type"] != "Role" and perm["name"] != "Owner":
                raise Exception("%s name must be Owner %s" % (perm["type"], pos))
            if perm["type"] == "Role" and perm["name"] == "Owner":
                raise Exception("Owner can not use Role type %s" % pos)
            item.append(perm)
        ret.append(item)
    return ret

def dealPermission(obj, idField, thrift_idl, labelName):
    location = "%s.%s" % (os.path.basename(thrift_idl), str(obj.name))
    obj.crud = dict()
    obj.create = permission(location, "Create", fieldElems(idField, "Create"))
    obj.read = permission(location, "Read", fieldElems(idField, "Read"))
    obj.update = permission(location, "Update", fieldElems(idField, "Update"))
    obj.delete = permission(location, "Delete", fieldElems(idField, "Delete"))
    obj.hasUser = len([i for i in obj.fields if str(i.name) == "UsersID"]) > 0
    if len(obj.create) + len(obj.read) + len(obj.update) + len(obj.delete) > 0:
        if not obj.hasUser:
            raise Exception("%s use crud and miss UsersID" % (location))
    obj.permission = []

    for items in obj.create + obj.read + obj.update + obj.delete:
        for perm in items:
            if perm["name"] == "Owner": continue
            p = dict(name=perm["name"],type=perm["type"],category=labelName)
            if p["type"] == "Role":
                p["type"] = ["Role", "User"]
            else:
                p["type"] = [p["type"]]
            p["type"] = ", ".join(["mp.Type." + i for i in p["type"]])
            obj.permission.append(p)
    if len(obj.permission) > 0:
        obj.imports.append("admin/permission")
        obj.imports.append(("mp", "zfw/models/permission"))


def main(thrift_idl):
    source = open(thrift_idl, 'r').read()
    loader = Loader(thrift_idl, lambda x: x)
    namespace = str(loader.namespace)

    tpl = open('tpl/create.tmpl', 'r').read().decode("utf8")
    crud = open('tmpl/crud.tmpl', 'r').read().decode("utf8")
    thrift = Parser().parse(source)

    ## generate controller files
    for obj in thrift.body:
        labelName = ""
        urlBase = ""
        obj.filterFields = []
        obj.imports = []

        if not hasattr(obj, "fields"):
            continue

        idField = obj.fields[0]
        labelName = fieldElem(idField, "label")
        urlBase = fieldElem(idField, "baseurl")

        if labelName == "" or urlBase == "":
            continue

        filterFields = fieldElem(idField, "filterfields")
        if filterFields != "":
            filterFields = [f.strip() for f in filterFields.split(",")]
        for fieldname in filterFields:
            for field in obj.fields:
                if field.name.value == fieldname:
                    obj.filterFields.append(field)

        dealPermission(obj, idField, thrift_idl, labelName)

        # make sure all field in filter fields exists
        if len(filterFields) > len(obj.filterFields):
            foundFields = [field.name.value for field in obj.filterFields if field.name.value in filterFields]
            missingFields = [field for field in filterFields if field not in foundFields]
            raise Exception("missing filterFields: " + str(missingFields))

        outDir = urlBase.split(path.sep)[-2]
        t = Template(crud, searchList=[{"namespace": outDir,
                                        "className": obj.name.value,
                                        "classLabel": labelName,
                                        "urlBase": urlBase,
                                        "obj": obj,
                                        }])
        res = str(t)
        writeDir = getControlDir(urlBase)
        outfile = writeDir + "gen_" + obj.name.value.lower() + ".go"
        f = open(outfile, "w")
        f.write(res)
        f.close()

        ## print "generate: {}".format(outfile)
        ## update init.go in controller
        updateController(out_path)

main(thrift_file)
