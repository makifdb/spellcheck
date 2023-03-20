# spellcheck

## Simple Spell Checker Package for Go
This is a simple spellchecker package for the Go programming language. It allows you to check the spelling of words in English. It also provides suggestions for misspelled words.

## Installation
To install the package, simply run the following command:

```go
go get github.com/makifdb/spellcheck
```

## Usage
To use the package, simply import it into your project and create a new spellchecker instance. You can then use the spellchecker to check the spelling of words and get suggestions for misspelled words. You can also insert new words into the dictionary.

```go
package main

import (
	"fmt"

	"github.com/makifdb/spellcheck"
)

func main() {
	// Init spellchecker
	sc, err := spellcheck.New()
	if err != nil {
		fmt.Println(err)
	}

	// Check spelling of a word
	ok := sc.SearchDirect("hllo")
	if !ok {
		fmt.Println("Word is misspelled or not in dictionary")
	}

	// Check spelling of a word and get suggestions
	ok, suggestions := sc.Search("hllo")
	if !ok {
		fmt.Println("Word is misspelled or not in dictionary")
		if len(suggestions) > 0 {
			fmt.Println("Did you mean: ", suggestions)
		}
	}

	// Insert a word into the dictionary
	sc.Insert("hllo")
}
```

## License
This package is licensed under the MIT license. See the LICENSE file for more details.