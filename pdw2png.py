import sys
from cStringIO import StringIO

from swf.movie import SWF
from swf.export import SVGExporter

def getimagedata(swf):
    for tag in swf.tags:
        if tag.name == 'DefineBinaryData':
            return tag.data


def doit(inf, outf):
    swf = SWF(inf)

    png = getimagedata(swf)

    outf.write(png)

def main():
    doit(sys.stdin, sys.stdout)

if __name__ == '__main__':
    main()
