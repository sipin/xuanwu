#!/usr/bin/env python
# -*- coding: utf-8 -*-

import sys
import os
import tinycss
from bs4 import BeautifulSoup

from ptsd.loader import Loader
from Cheetah.Template import Template
from os import path

class Size:
	def __init__(self):
		self.top = 0
		self.left = 0

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
			size = get_size(r.declarations)
			if size.top > 0:
				data[str(r.selector[1].value)] = size

	return data

def gen(axure_folder, key, output):
	if not axure_folder.endswith(path.sep):
		axure_folder = axure_folder + path.sep

	css_path = axure_folder + key + "_files/axurerp_pagespecificstyles.css"
	html_path = axure_folder + key + ".html"

	parser = tinycss.make_parser('page3')
	css = parser.parse_stylesheet_file(css_path)

	rd = get_ruledict(css.rules)

	soup = BeautifulSoup(open(html_path))

	text_dict = {}

	for i in range(1, len(rd)):
		node_id = 'u%d' % i
		nodes = soup.find_all(id = node_id)
		if len(nodes) == 1:
			key = nodes[0].get_text().encode("utf-8").strip()
			key = key.replace(" ", "").replace("*", "").replace("c2a0".decode("hex"), "")
			if len(key) < 2 or  "\n" in key[:-1]:
				continue
			text_dict[key] = node_id

	rows = {}
	data = {}
	for text in text_dict.keys():
		key = text_dict[text]
		if rd.has_key(key):
			size = rd[key]
			rows[text] = size.top + float(size.left) / (len(str(size.left)) * 10)
			data[text] = size
	top = 0
	result = ""
	for w in sorted(rows, key=rows.get):
		size = data[w]
		if size.top - top > 10:
			if top > 0:
				result += "\n"
			top = size.top
		else:
			result += "\t"
		result += w

	f = file(output, "w+")
	f.write(result)
	f.close()

	return result

def main():
	if len(sys.argv) != 4:
		print "usage: \n\tpython axure.py axure_folder key output"
		sys.exit()

	axure_folder = sys.argv[1]
	key = sys.argv[2]
	output = sys.argv[3]
	gen(axure_folder, key, output)

def assure_path_exists(path):
	dir = os.path.dirname(path)
	if not os.path.exists(dir):
		os.makedirs(dir)

def gen_axure(thrift_file, obj):
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

if __name__ == "__main__":
	main()
