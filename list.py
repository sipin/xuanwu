#!/usr/bin/env python
# -*- coding: utf-8 -*-

import sys
import base

thrift_file = sys.argv[1]

def main(thrift_idl):
	loader = base.load_thrift(thrift_idl)
	namespace = loader.namespace
	print '	_ "zfw/models/' + namespace + '"'

main(thrift_file)
