package Utils

import (
	"strings"
	"strconv"
)

func tokenizer(c rune) bool {
	return strings.ContainsRune(" #''\":*~?!%(),<=>[]^{}\r\n\t\f.", c)
}

func ZdbBuildQuery(field, searchRequest string) string {
	words := strings.FieldsFunc(searchRequest, tokenizer)
	wordsLen := len(words);
	fixedTerm:= strings.Join(words, " ")
	for i := 0; i < wordsLen; i++ {
		if len(words[i]) == 0 {
			words = append(words[:i], words[:i + 1]...)
			i--;
		} else {
			if zdbShouldEscape(words[i]) {
				words[i] = strconv.Quote(words[i])
			}
		}
	}

	var nameQuery string;
	if wordsLen == 1 {
		nameQuery = " '" + field + ":(" + words[0] + ")'";
	} else {
		nameQuery = " '" + field + ":(" + fixedTerm + " or " + strings.Join(words, " or ") + ")'"
	}

	return nameQuery;
}

var zdbKeyWords []string = []string{
	"with",
	"and",
	"or",
	"not",
}

func zdbShouldEscape(s string) bool {
	return IndexOf(zdbKeyWords, strings.ToLower(s)) > -1;
}