package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const MaxPokemon = 719
const BaseUrl = "http://3ds.pokemon-gl.com/share/images/pokemon"

var formNames = map[int][]string{
	201: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
		"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"exclamation", "question"},
	351: []string{"", "sunny", "rainy", "snowy"},
	386: []string{"normal", "attack", "defense", "speed"},
	412: []string{"plant", "sandy", "trash"},
	413: []string{"plant", "sandy", "trash"},
	421: []string{"overcast", "sunshine"},
	422: []string{"west", "east"},
	423: []string{"west", "east"},
	479: []string{"", "heat", "wash", "frost", "fan", "mow"},
	487: []string{"altered", "origin"},
	492: []string{"land", "sky"},
	493: []string{"normal", "fighting", "flying", "poison", "ground",
		"rock", "bug", "ghost", "steel", "fire", "water", "grass",
		"electric", "psychic", "ice", "dragon", "dark", "fairy"},
	550: []string{"red-striped", "blue-striped"},
	555: []string{"standard", "zen"},
	585: []string{"spring", "summer", "autumn", "winter"},
	586: []string{"spring", "summer", "autumn", "winter"},
	592: []string{"male", "female"},
	593: []string{"male", "female"},
	641: []string{"incarnate", "therian"},
	642: []string{"incarnate", "therian"},
	645: []string{"incarnate", "therian"},
	646: []string{"", "white", "black"},
	647: []string{"ordinary", "resolute"},
	648: []string{"aria", "pirouette"},
	649: []string{"", "douse", "shock", "burn", "chill"},
	666: []string{"icy-snow", "polar", "tundra", "continental", "garden",
		"elegant", "meadow", "modern", "marine", "archipelago",
		"high-plains", "sandstorm", "river", "monsoon", "savanna",
		"sun", "ocean", "jungle"},
	668: []string{"red", "yellow", "orange", "blue", "white"},
	669: []string{"red", "yellow", "orange", "blue", "white"},
	670: []string{"red", "yellow", "orange", "blue", "white"},
	676: []string{"natural", "heart", "star", "diamond", "deputante", "matron", "dandy", "la-reine", "kabuki", "pharaoh"},
	678: []string{"male", "female"},
	681: []string{"shield", "blade"},

	3:   []string{"", "mega"},
	6:   []string{"", "mega-x", "mega-y"},
	9:   []string{"", "mega"},
	65:  []string{"", "mega"},
	94:  []string{"", "mega"},
	115: []string{"", "mega"},
	127: []string{"", "mega"},
	130: []string{"", "mega"},
	142: []string{"", "mega"},
	150: []string{"", "mega-x", "mega-y"},
	181: []string{"", "mega"},
	212: []string{"", "mega"},
	214: []string{"", "mega"},
	229: []string{"", "mega"},
	248: []string{"", "mega"},
	257: []string{"", "mega"},
	282: []string{"", "mega"},
	303: []string{"", "mega"},
	306: []string{"", "mega"},
	308: []string{"", "mega"},
	310: []string{"", "mega"},
	354: []string{"", "mega"},
	359: []string{"", "mega"},
	445: []string{"", "mega"},
	448: []string{"", "mega"},
	460: []string{"", "mega"},
}

func getUrl(n, id int) string {
	r := 0x159a55e5 * uint(n+id*0x10000) & 0xffffff
	return fmt.Sprintf("%s/300/%06x.png", BaseUrl, r)
}

func getFilename(n, id int) string {
	names, ok := formNames[n]
	name := ""
	if ok {
		if id < len(names) {
			name = names[id]
		} else {
			name = fmt.Sprint(id)
		}
	} else {
		if id > 0 {
			name = fmt.Sprint(id)
		}
	}
	if name == "" {
		return fmt.Sprintf("%d.png", n)
	} else {
		return fmt.Sprintf("%d-%s.png", n, name)
	}
}

var NotFound = errors.New("pokemon not found")

func getFile(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return NotFound
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	for n := 1; n <= MaxPokemon+5; n++ {
		if _, ok := formNames[n]; !ok {
			continue
		}
		id := 0
		for ; ; id++ {
			url := getUrl(n, id)
			filename := getFilename(n, id)
			err := getFile(url, filename)
			if err == NotFound {
				break
			} else if err != nil {
				log.Println(err)
			}
		}
		if id == 0 {
			break
		}
	}
}
