#!/usr/bin/python
# -*- coding: utf-8 -*-

from __future__ import print_function
import sys
import urllib,urllib2
from bs4 import BeautifulSoup

ip = sys.argv[1]

url = "http://site.ip138.com/" + ip

req = urllib2.Request(url)
bshtml = BeautifulSoup(urllib2.urlopen(req).read(), features="lxml")

lis = bshtml.find_all("li")

i = 0
l = len(lis)

while(i < l):
    li = lis[i]
    if li.find("span") is not None:
        s = unicode(li.find("span").get_text())
#        print s
        if s.find(u"绑定过的域名如下") != -1:
            i = i + 1
            break
    i = i + 1

li = lis[i]

print(li.find("a").get_text(), end="")
