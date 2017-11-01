package currency

import (
	"github.com/andersfylling/IMT2681-2/database/documents/document"
)

// Document to be stored in database
type Document struct {
	Base  string
	Date  string
	Rates map[string]int
}

// New Creates a new instance of the document.
// Which can then be saved, removed, find matches, etc.
func New() document.Interface {
	return &Document{}
}

// Inserts the document as a new one into the collection and returns the id
func (c *Document) Insert() (id string, err error) {
	id = ""
	err = nil

	return id, err
}

// Save updates a document that already exists, and then return the old and the new.
func (c *Document) Save() (old, new document.Interface) {

	return old, new
}

// Remove can remove documents in bulk, and deleted documents are returned.
// Any document that fits the rule will get deleted.
// If the array is empty, then no documents where deleted.
// int equals their old ID
func (c *Document) Remove() map[int]document.Interface {
	results := make(map[int]document.Interface)

	return results
}

// Find returns an empty array when no match was found
func (c *Document) Find() map[int]document.Interface {
	results := make(map[int]document.Interface)

	return results
}

// make sure struct implements interface
var _ document.Interface = (*Document)(nil)
