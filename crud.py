# -*- coding: utf-8 -*-
import sys
import os
import base

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

outDir = out_path.split(path.sep)[-2]

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

def dealPermission(obj, idField, labelName):
    location = "%s.%s" % (os.path.basename(thrift_file), str(obj.name))
    obj.crud = dict()
    create = fieldElems(idField, "Create")
    obj.create = permission(location, "Create", create)
    obj.read = permission(location, "Read", fieldElems(idField, "Read"))
    obj.update = permission(location, "Update", fieldElems(idField, "Update"))
    obj.delete = permission(location, "Delete", fieldElems(idField, "Delete"))

    obj.hasUser = len([i for i in obj.fields if str(i.name) == "UsersID"]) > 0
    if len(obj.read) == 0 and obj.hasUser:
        obj.read = [[dict(name="Owner",type="Department")]]

    if len(obj.delete) == 0 and obj.hasUser:
        obj.delete = [[dict(name="Owner",type="User")]]

    if len(obj.update) == 0 and obj.hasUser:
        obj.update = [[dict(name="Owner",type="User")]]

    obj.hasDelete = len(obj.delete) > 0
    obj.hasCreate = len(obj.create) > 0
    obj.hasUpdate = len(obj.update) > 0
    obj.hasRead = len(obj.read) > 0
    if len(obj.create) + len(obj.read) + len(obj.update) + len(obj.delete) > 0:
        obj.imports.append("admin/permission")
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
        obj.imports.append(("mp", "zfw/models/permission"))

def transform_module(module):
    for obj in module.structs:        
        urlBase = ""
        obj.imports = []
        idField = obj.fields[0]        
        urlBase = fieldElem(idField, "baseurl")
        tplPackage = fieldElem(idField, "tplpackage")
        if tplPackage == "":
            tplPackage = "tpl/auto"

        if obj.label  == "" or urlBase == "":
            continue
    
        dealPermission(obj, idField, obj.label)
        if len(obj.relateObj) > 0:
            obj.imports.append("encoding/json")

        outDir = urlBase.split(path.sep)[-2]
        crud = open('tmpl/crud.tmpl', 'r').read().decode("utf8")
        t = Template(crud, searchList=[{"namespace": outDir,
                                        "className": obj.name.value,
                                        "classLabel": obj.label,
                                        "urlBase": urlBase,
                                        "tplPackage": tplPackage,
                                        "obj": obj,
                                        }])
        res = str(t)
        writeDir = getControlDir(urlBase)
        outfile = writeDir + "gen_" + obj.name.value.lower() + ".go"
        with open(outfile, "w") as fp:
            fp.write(res)         
            
    
def main(thrift_idl):
    loader = base.load_thrift(thrift_idl)
    global namespace
    namespace = loader.namespace
    for obj in loader.modules.values():
        transform_module(obj)
    updateController(out_path)

main(thrift_file)
