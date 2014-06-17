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

outDir = out_path.split(path.sep)[-1]

def getViewDir(urlBase):
    outdir = out_path + "tpl/" + \
      (path.sep).join(urlBase.split(path.sep)[:-1]) + path.sep
    if not os.path.exists(outdir):
        os.makedirs(outdir)
    return outdir

def getControlDir(urlBase):
    outdir = out_path + \
      (path.sep).join(urlBase.split(path.sep)[:-1]) + path.sep
    if not os.path.exists(outdir):
        os.makedirs(outdir)
    return outdir

def add_properties(field):
    field.label = field.name.value
    for att in field.annotations:
        if att.name.value.lower() == "label":
            field.label = att.value.value


def fieldElem(field, key):
    for att in field.annotations:
        if att.name.value.lower() == key:
            return att.value.value
    return ""


def main(thrift_idl):
    source = open(thrift_idl, 'r').read()
    loader = Loader(thrift_idl, lambda x: x)
    namespace = str(loader.namespace)

    tpl = open('tpl/create.tmpl', 'r').read().decode("utf8")
    crud = open('tmpl/crud.tmpl', 'r').read().decode("utf8")
    thrift = Parser().parse(source)

    ## generate view files
    # for obj in thrift.body:
    #     for field in obj.fields:
    #         add_properties(field)

    #     t = Template(tpl, searchList=[{"obj": obj}])
        ##print code.encode("utf8")

    ## generate controller files
    for obj in thrift.body:
        labelName = ""
        urlBase = ""
        if not hasattr(obj, "fields"):
            continue
        for field in obj.fields:
            if field.name.value == "ID":
                labelName = fieldElem(field, "label")
                urlBase = fieldElem(field, "baseurl")
                break

        if len(labelName) == 0:
            pass
        elif len(urlBase) == 0:
            pass
        else:
            for field in obj.fields:
                add_properties(field)
            t = Template(tpl, searchList=[{"obj": obj}])
            res = str(t)
            outfile = getViewDir(urlBase) + "gen_" + obj.name.value.lower() + ".go"

            t = Template(crud, searchList=[{"namespace": outDir,
                                            "className": obj.name.value,
                                            "classLabel": labelName,
                                            "urlBase": urlBase,
                                            }])
            res = str(t)
            outfile = out_path + "gen_" + obj.name.value.lower() + ".go"
            f = open(outfile, "w")
            f.write(res)
            f.close()

main(thrift_file)
