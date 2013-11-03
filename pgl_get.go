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
	201: {"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
		"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"exclamation", "question"},
	351: {"", "sunny", "rainy", "snowy"},
	386: {"normal", "attack", "defense", "speed"},
	412: {"plant", "sandy", "trash"},
	413: {"plant", "sandy", "trash"},
	421: {"overcast", "sunshine"},
	422: {"west", "east"},
	423: {"west", "east"},
	479: {"", "heat", "wash", "frost", "fan", "mow"},
	487: {"altered", "origin"},
	492: {"land", "sky"},
	493: {"normal", "fighting", "flying", "poison", "ground",
		"rock", "bug", "ghost", "steel", "fire", "water", "grass",
		"electric", "psychic", "ice", "dragon", "dark", "fairy"},
	550: {"red-striped", "blue-striped"},
	555: {"standard", "zen"},
	585: {"spring", "summer", "autumn", "winter"},
	586: {"spring", "summer", "autumn", "winter"},
	592: {"male", "female"},
	593: {"male", "female"},
	641: {"incarnate", "therian"},
	642: {"incarnate", "therian"},
	645: {"incarnate", "therian"},
	646: {"", "white", "black"},
	647: {"ordinary", "resolute"},
	648: {"aria", "pirouette"},
	649: {"", "douse", "shock", "burn", "chill"},
	666: {"icy-snow", "polar", "tundra", "continental", "garden",
		"elegant", "meadow", "modern", "marine", "archipelago",
		"high-plains", "sandstorm", "river", "monsoon", "savanna",
		"sun", "ocean", "jungle"},
	669: {"red", "yellow", "orange", "blue", "white"},
	670: {"red", "yellow", "orange", "blue", "white"},
	671: {"red", "yellow", "orange", "blue", "white"},
	676: {"natural", "heart", "star", "diamond", "debutante",
		"matron", "dandy", "la-reine", "kabuki", "pharaoh"},
	678: {"male", "female"},
	681: {"shield", "blade"},

	3:   {"", "mega"},
	6:   {"", "mega-x", "mega-y"},
	9:   {"", "mega"},
	65:  {"", "mega"},
	94:  {"", "mega"},
	115: {"", "mega"},
	127: {"", "mega"},
	130: {"", "mega"},
	142: {"", "mega"},
	150: {"", "mega-x", "mega-y"},
	181: {"", "mega"},
	212: {"", "mega"},
	214: {"", "mega"},
	229: {"", "mega"},
	248: {"", "mega"},
	257: {"", "mega"},
	282: {"", "mega"},
	303: {"", "mega"},
	306: {"", "mega"},
	308: {"", "mega"},
	310: {"", "mega"},
	354: {"", "mega"},
	359: {"", "mega"},
	445: {"", "mega"},
	448: {"", "mega"},
	460: {"", "mega"},
}

func getUrl(n, id int) string {
	r := 0x159a55e5 * uint(n+id*0x10000) & 0xffffff
	return fmt.Sprintf("%s/300/%06x.png", BaseUrl, r)
}

func getFilename(n, id int) string {
	names, ok := formNames[n]
	if !ok {
		names = []string{""}
	}
	name := ""
	if id < len(names) {
		name = names[id]
	} else {
		name = fmt.Sprint(id)
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
	for n := 1; n <= MaxPokemon; n++ {
		/*if _, ok := formNames[n]; !ok {
			continue
		}*/
		id := 0
		for ; ; id++ {
			url := getUrl(n, id)
			filename := getFilename(n, id)
			err := getFile(url, filename)
			if err == NotFound {
				break
			}
			if err != nil {
				log.Println(err)
			}
		}
		if id == 0 {
			break
		}
	}
}
