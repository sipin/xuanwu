#!/usr/bin/env python
# -*- coding: utf-8 -*-

import sys
import os
import tinycss

from ptsd import ast
from ptsd.loader import Loader, Parser
from Cheetah.Template import Template
from os import path

if len(sys.argv) != 3:
	print "usage: \n\tpython axure.py axure_folder key"
	sys.exit()

axure_folder = sys.argv[1]
key = sys.argv[2]

def readall(fname):
	f = open(fname)
	content = f.read()
	f.close()
	return content

class Size:
	def __init__(self):
		self.top = ""
		self.left = ""

	def __str__(self):
		return "top: " + str(self.top) + " left: " + str(self.left)

def get_size(declarations):
	s = Size()
	for d in declarations:
		if d.name == "left":
			s.left = d.value[0].value
		if d.name == "top":
			s.top = d.value[0].value

	return s

def get_ruledict(rs):
	data = {}
	for r in rs:
		if len(r.selector) == 2 and r.selector[1].value[0] == "u":
			data[str(r.selector[1].value)] = get_size(r.declarations)

	return data

def gen(axure_folder, key):
	if not axure_folder.endswith(path.sep):
		axure_folder = axure_folder + path.sep

	css_path = axure_folder + key + "_files/axurerp_pagespecificstyles.css"
	html_path = axure_folder + key + ".html"

	parser = tinycss.make_parser('page3')
	css = parser.parse_stylesheet_file(css_path)

	rd = get_ruledict(css.rules)

	print rd

def main():
	gen(axure_folder, key)

main()
