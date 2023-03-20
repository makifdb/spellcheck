package spellcheck

import (
	"bufio"
	"io"
	"net/http"
	"sync"
)

const (
	// DEFAULT_ALPHABET is the default alphabet used for generating variations
	DEFAULT_ALPHABET = "abcdefghijklmnopqrstuvwxyz"
	// DEFAULT_REMOTE_WORDS_URL is the default url containing the words
	DEFAULT_REMOTE_WORDS_URL = "https://raw.githubusercontent.com/makifdb/spellcheck/main/words.txt"
	// DEFAULT_DEPTH is the default depth used for generating variations
	DEFAULT_DEPTH = 1
)

type LetterNode struct {
	// children is a map of the children nodes
	children map[rune]*LetterNode
	// isWord is true if the node is the end of a word
	isWord bool
}

type Trie struct {
	//	root is the root node of the trie
	root *LetterNode
	// depth is the depth used for generating variations
	depth int
	// mutex is used for locking the trie
	sync.RWMutex
}

// Insert adds a word to the trie
func (t *Trie) Insert(word string) {
	t.Lock()
	defer t.Unlock()

	// create the root node if it doesn't exist
	if t.root == nil {
		t.root = &LetterNode{make(map[rune]*LetterNode), false}
	}

	// check if the word is already in the trie
	node := t.root
	for _, char := range word {
		if _, ok := node.children[char]; !ok {
			node.children[char] = &LetterNode{make(map[rune]*LetterNode), false}
		}
		node = node.children[char]
	}

	// check if the word is already in the trie
	if node.isWord {
		return
	}

	// mark the end of the word
	node.isWord = true
}

// InsertReader adds words from a reader to the trie
func (t *Trie) InsertReader(r io.Reader) *bufio.Scanner {
	// create a scanner to read the file line by line
	scanner := bufio.NewScanner(r)

	// read the file line by line
	for scanner.Scan() {
		word := scanner.Text()
		t.Insert(word)
	}

	// check for errors
	if err := scanner.Err(); err != nil {
		return nil
	}

	return scanner
}

// Search checks if a word is in the trie and returns with a list of suggestions
func (t *Trie) Search(word string) (bool, []string) {
	t.RLock()
	defer t.RUnlock()

	node := t.root

	// check the tire for the full word
	for _, char := range word {
		if _, ok := node.children[char]; !ok {
			// word not found, generate suggestions

			// generate variations
			suggestions := generateVariations(word, t.depth)
			var validSuggestions []string

			// check if the suggestions are in the trie
			for _, suggestion := range suggestions {
				if found := t.SearchDirect(suggestion); found {

					// add the suggestion to the list of valid suggestions
					validSuggestions = append(validSuggestions, suggestion)
				}
			}

			// return false and the list of suggestions
			return false, validSuggestions
		}

		// move to the next node
		node = node.children[char]
	}

	// word found
	return node.isWord, nil
}

// SearchDirect checks if a word is in the trie
func (t *Trie) SearchDirect(word string) bool {
	t.RLock()
	defer t.RUnlock()

	// check the tire for the full word
	node := t.root
	for _, char := range word {
		if _, ok := node.children[char]; !ok {
			return false
		}
		node = node.children[char]
	}

	return node.isWord
}

// generateVariations generates variations of a word with a given depth
func generateVariations(word string, depth int) []string {

	// check if the depth is 0 or the word is empty
	if depth == 0 || len(word) == 0 {
		return []string{}
	}

	var result []string

	// Deletes
	// word -> ord, wod, wrd, wor
	for i := 0; i < len(word); i++ {
		result = append(result, word[:i]+word[i+1:])
	}

	// Transposes
	// word -> owrd, wrod, wodr, word
	for i := 0; i < len(word)-1; i++ {
		result = append(result, word[:i]+string(word[i+1])+string(word[i])+word[i+2:])
	}

	// Replaces
	// word -> aord, bord, cord, ... zord
	for i := 0; i < len(word); i++ {
		for _, c := range DEFAULT_ALPHABET {
			result = append(result, word[:i]+string(c)+word[i+1:])
		}
	}

	// Inserts
	// word -> aword, bword, cword, ... zword
	for i := 0; i < len(word)+1; i++ {
		for _, c := range DEFAULT_ALPHABET {
			result = append(result, word[:i]+string(c)+word[i:])
		}
	}

	// Recursive call
	// generate variations for each variation
	for _, variation := range result {
		result = append(result, generateVariations(variation, depth-1)...)
	}

	return clearSuggestions(result)
}

// clearSuggestions removes duplicates from a list of suggestions
func clearSuggestions(suggestions []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, suggestion := range suggestions {
		if encountered[suggestion] {
			continue
		}

		encountered[suggestion] = true
		result = append(result, suggestion)
	}

	return result
}

// New creates a new trie and populates it with the words
func New() (*Trie, error) {
	// create a new trie
	t := &Trie{&LetterNode{make(map[rune]*LetterNode), false}, DEFAULT_DEPTH, sync.RWMutex{}}

	resp, err := http.Get(DEFAULT_REMOTE_WORDS_URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// create a scanner to read the file line by line and insert to trie
	t.InsertReader(resp.Body)

	return t, nil
}
