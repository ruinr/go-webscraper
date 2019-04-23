package v1

import (
	"html"
	"strings"
)

//ConvertHTMLEntities converts scraped html entities and returns clean string
func ConvertHTMLEntities(htmlString string) (clean string) {
	clean = html.UnescapeString(htmlString)
	//some entities cannot be convert by using UnescapeString
	//It will requires extra string replacements using string replacer

	replacer := strings.NewReplacer(
		//set with pairs
		"\u0026", "&",
		"&amp;", "&",
	)

	clean = replacer.Replace(clean)
	return
}
