#!/usr/bin/env python3

"""svgclean.py - Clean up svg files produced by pyswf."""

import sys
import re
from lxml import etree
from warnings import warn
from collections import Counter

namespaces = {
    'svg': "http://www.w3.org/2000/svg",
    'xlink': "http://www.w3.org/1999/xlink",
}


def svgclean(tree):
    # Embed groups in <defs> which are only used once.
    use_elements = tree.xpath("//svg:use[@xlink:href]", namespaces=namespaces)
    g_elements = tree.xpath("/svg:svg/svg:defs/svg:g[@id]", namespaces=namespaces)

    #print(use_elements, g_elements, sep='\n', file=sys.stderr)

    revlinks = {}

    for e in use_elements:
        href = e.get('{http://www.w3.org/1999/xlink}href')
        if href is not None and href.startswith("#"):
            revlinks.setdefault(href, []).append(e)

    for g in g_elements:
        try:
            href = "#" + g.attrib['id']
            uses = revlinks[href]
        except LookupError:
            pass
        finally:
            if len(uses) == 1:
                use = uses[0]
                new_g = etree.Element('{http://www.w3.org/2000/svg}g')
                for k in use.attrib:
                    if k not in {'x', 'y', 'width', 'height', '{http://www.w3.org/1999/xlink}href'}:
                        new_g.attrib[k] = use.attrib[k]

                point = (float(use.get('x', '0')), float(use.get('y', '0')))
                if point != (0.0, 0.0):
                    translate = 'translate(%s, %s)' % point
                    if new_g.get('transform'):
                        new_g.set('transform', new_g.get('transform') + ' ' + translate)
                    else:
                        new_g.set('transform', translate)

                del g.attrib['id']
                new_g.append(g)

                use.getparent().replace(use, new_g)

    # Simplify transforms
    for e in tree.xpath("//*[@transform]", namespaces=namespaces):
        def transform_transformer(transform):
            function, args = transform
            if function == 'matrix':
                if args == [1.0, 0.0, 0.0, 1.0, 0.0, 0.0]:
                    # Identity transform
                    return
                elif args[:4] == [1.0, 0.0, 0.0, 1.0]:
                    return 'translate', args[4:6]
                elif args[1:3] == args[4:6] == [0.0, 0.0]:
                    return 'scale', [args[0], args[3]]
                else:
                    return function, args
            else:
                return function, args

        transforms = parse_transform(e.get('transform'))
        transforms = list(filter(None, map(transform_transformer, transforms)))
        e.set('transform', unparse_transform(transforms))

        if not e.get('transform'):
            del e.attrib['transform']

    # Remove empty defs
    for defs in tree.xpath("//svg:defs", namespaces=namespaces):
        if len(defs) == 0:
            defs.getparent().remove(defs)

    # Collapse nested groups
    for g in tree.xpath("//svg:g", namespaces=namespaces):
        if len(g) == 1:
            child = list(g)[0]
            #print(child, file=sys.stderr)
            if child.tag == '{http://www.w3.org/2000/svg}g':
                for k in g.attrib:
                    if k == 'transform' and k in child.attrib:
                        child.attrib[k] = g.attrib[k] + " " + child.attrib[k]
                    elif k not in child.attrib:
                        child.attrib[k] = g.attrib[k]
                g.getparent().replace(g, child)

    for g in tree.xpath("//svg:g", namespaces=namespaces):
        if len(g) == 1 and not g.attrib:
            child = list(g)[0]
            if child.tag == '{http://www.w3.org/2000/svg}path':
                g.getparent().replace(g, child)

    # Strip useless stroke attributes
    for path in tree.xpath("//svg:path[@stroke='none']", namespaces=namespaces):
        del path.attrib['stroke']

    # Cleanup namespaces
    etree.cleanup_namespaces(tree)

    return tree


def parse_transform(s):
    transforms = []
    for function, argstring in re.findall(r'(\w+)\(([-0-9., ]*)\)', s):
        args = list(map(float, re.split('\s+|\s*,\s*', argstring)))
        transforms.append((function, args))

    return transforms

def unparse_transform(transforms):
    return ' '.join(
        "%s(%s)" % (function, ','.join(map(str, args)))
        for function, args in transforms
    )

def main(args):
    try:
        filename = args[0]
    except IndexError:
        filename = '-'

    if filename == '-':
        filename = sys.stdin

    tree = etree.parse(filename)
    tree = svgclean(tree)
    s = etree.tostring(tree, encoding="utf-8", xml_declaration=True)
    sys.stdout.buffer.write(s)

if __name__ == '__main__':
    main(sys.argv[1:])
