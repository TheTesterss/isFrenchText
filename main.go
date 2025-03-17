package main;

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	"net/http"
	"strings"
)

var chars []string = []string{
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	"é", "è", "ç", "ê", "à", "â", "ë", "ä",
}
var path string = "./content.txt"
var larousseApiPath string = "https://larousse.fr/dictionnaires/francais"

func main() {
	// First we read the content from the file we may verify its wrote in french.
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("An error occured while reading: %s\n%s", path, err)
		return
	}

	var r int = 0
	var wordsAmount int = 0
	for _, word := range strings.Split(string(content), " ") {
		word = cleanWord(strings.ToLower(word))
		if word == "" {
			continue
		}
		wordsAmount++

		if !isWordFrench(word) {
			r++
		}
	}
	
	if r == 0 {
		fmt.Println("The text is 100% french.")
	} else {
		fmt.Printf("The text is %.2f%% french.", (float64(wordsAmount-r)/float64(wordsAmount))*100)
	}
}

func isWordFrench(word string) bool {
	// Parsing the larousse's website to retrieve the result page for a specific word.
	// 		=> In case it does exist, a page with such informations as definition, equivalents, opposites... is given.
	//		=> In case it doesn't exist, the same page as above is rendered but the word defined is a similary wrote word not the requested one.

	// First we start reading the page html balises.
	html, err := http.Get(larousseApiPath+"/"+word)
	if err != nil {
		return false
	}
	defer html.Body.Close()

	// Now we use the parser to find informations for a specific class/id.
	document, err := goquery.NewDocumentFromReader(html.Body)
	if err != nil {
		return false
	}

	// Makes a list of every words that larousse found for us.
	var foundWords []string
	document.Find("h2.AdresseDefinition").Each(func(index int, element *goquery.Selection) {
		for _, element_1 := range strings.Split(element.Text(), ", ") {
			foundWords = append(foundWords, cleanWord(element_1))
		}
	})

	// Verifies the presence of the word in the made list.
	return contains(foundWords, word)
}

func cleanWord(word string) string {
	// Removes accentued chars but also digits, punctuations...
	var cleaned strings.Builder
	for _, char := range word {
		if contains(chars, string(char)) {
			cleaned.WriteRune(char)
		}
	}
	return cleaned.String()
}

func contains(list []string, char string) bool {
	// Verifies if a value is present in the list.
	for _, c := range list {
		if c == char {
			return true
		}
	}
	return false
}