#!/usr/bin/env python3

PDW_TABLE = str.maketrans('0123456789', 'abcdefghij')
URL = "http://cdn2.pokemon-gl.com/src/swf/theme/assets/global/parts/pokemon/scaled/{0}.swf"
NUTS_PATH = "http://cdn2.pokemon-gl.com/src/swf/theme/assets/global/parts/nuts/90/{0}.swf"
ITEM_PATH = "http://cdn2.pokemon-gl.com/src/swf/theme/assets/global/parts/item/90/{0}.swf"

# Item sizes: 28, 90
# Nuts sizes: 28, 35, 57, 90, 

FORMS = {
    201: list('abcdefghijklmnopqrstuvwxyz') + ['exclamation', 'question'],
    351: ['', 'sunny', 'rainy', 'snowy'],
    386: ['normal', 'attack', 'defense', 'speed'],
    412: ['plant', 'sandy', 'trash'],
    413: ['plant', 'sandy', 'trash'],
    421: ['overcast', 'sunshine'],
    422: ['west', 'east'],
    423: ['west', 'east'],
    479: ['', 'heat', 'wash', 'frost', 'fan', 'mow'],
    487: ['altered', 'origin'],
    492: ['land', 'sky'],
    493: ['normal', 'fighting', 'flying', 'poison', 'ground', 'rock', 'bug',
          'ghost', 'steel', 'fire', 'water', 'grass', 'electric', 'psychic',
          'ice', 'dragon', 'dark'],
    550: ['red-striped', 'blue-striped'],
    555: ['standard', 'zen'],
    585: ['spring', 'summer', 'autumn', 'winter'],
    586: ['spring', 'summer', 'autumn', 'winter'],
    592: ['', 'female'],
    593: ['', 'female'],
    641: ['incarnate', 'therian'],
    642: ['incarnate', 'therian'],
    645: ['incarnate', 'therian'],
    646: ['', 'black', 'white'],
    647: ['', 'resolute'],
    648: ['aria', 'pirouette'],
    649: ['', 'douse', 'shock', 'burn', 'chill'],
}


def encrypt_kinomi(n):
    x = (n << 8) ^ 0xc3c3c3
    return str(x).translate(PDW_TABLE)

def geturl(number, form=0):
    return URL.format(encrypt_kinomi(number << 8 | form))


# Items:
# if 149 <= id <= 213:
#     nut!
#     id - 148

def getitem(item_id):
    if 149 <= item_id <= 213:
        return NUTS_PATH.format(encrypt_kinomi(item_id - 148))
    else:
        return ITEM_PATH.format(encrypt_kinomi(item_id))


#for number in range(1, 649 + 1):
#    print(geturl(number))

#for form in range(4):
#    print(geturl(386, form))

import sys
if len(sys.argv) == 1:
    for number in range(1, 649 + 1):
        forms = FORMS.get(number, [''])
        for form_index, form_name in enumerate(forms):
            if form_name:
                name = "{}-{}.swf".format(number, form_name)
            else:
                name = "{}.swf".format(number)
            print(geturl(number, form_index), name)
elif len(sys.argv) == 2:
    if sys.argv[1] == 'missing':
        for number in range(1, 649+1):
            form_index = len(FORMS.get(number, ['']))
            print(geturl(number, form_index))
    else:
        print(geturl(int(sys.argv[1]), 0))
elif len(sys.argv) == 3:
    print(geturl(int(sys.argv[1]), int(sys.argv[2])))
else:
    print("Usage: pdw_urls.py [number [form]]", file=sys.stderr)

