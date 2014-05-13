# -*- coding: utf-8 -*-
import sys

from ptsd import ast
from ptsd.loader import Parser
from Cheetah.Template import Template
from os import path


if len(sys.argv) != 2:
	print "usage: \n\tpython crud.py thrift_file_path"
	sys.exit()

thrift_file = sys.argv[1]

def add_properties(field):
	field.label = field.name.value

	for att in field.annotations:
		if att.name.value.lower() == "label":
			field.label = att.value.value


def main(thrift_idl):
	source = open(thrift_idl, 'r').read()

	tpl = open('tpl/create.tmpl', 'r').read().decode("utf8")

	thrift = Parser().parse(source)
	for obj in thrift.body:
		for field in obj.fields:
			add_properties(field)


		t = Template(tpl, searchList=[{"obj": obj}])
		code = unicode(t)
		print code.encode("utf8")


main(thrift_file)
