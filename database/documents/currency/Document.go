package currency

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andersfylling/IMT2681-2/database/dbsession"
	"github.com/andersfylling/IMT2681-2/database/documents/document"
	"gopkg.in/mgo.v2/bson"
)

// Document to be stored in database
type Document struct {
	ID    bson.ObjectId      `json:"_id" bson:"_id,omitempty"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

const Collection = "CurrencyRates"

// New Creates a new instance of the document.
// Which can then be saved, removed, find matches, etc.
func New() *Document {
	return &Document{
		//ID:     nil,
		Base:  "",
		Date:  "",
		Rates: map[string]float64{},
	}
}

// NewFromRequest Uses a http.Request object to populate the Document from body content
func NewFromRequest(r *http.Request) (*Document, error) {
	decoder := json.NewDecoder(r.Body)
	d := New()
	err := decoder.Decode(d)

	return d, err
}

// Insert the document as a new one into the collection and returns the id
func (c *Document) Insert() (id bson.ObjectId, err error) {
	id = ""
	err = nil

	ses, con, err := dbsession.GetCollection(Collection)
	if err != nil {
		return id, err
	}
	defer ses.Close()

	c.ID = bson.NewObjectId()
	err = con.Insert(c)

	return c.ID, err
}

// Save updates a document that already exists, and then return the old and the new.
func (c *Document) Save() (old, new document.Interface) {

	return old, new
}

func interfaceToDocument(interfaces []interface{}) []*Document {
	docs := make([]*Document, len(interfaces))
	for i, v := range interfaces {
		docs[i] = v.(*Document)
	}

	return docs
}

func documentsToInterfaces(docs []*Document) []interface{} {
	interfaces := make([]interface{}, len(docs))
	for i, v := range docs {
		interfaces[i] = v
	}

	return interfaces
}

// Remove can remove documents in bulk, and deleted documents are returned.
// Any document that fits the rule will get deleted.
// If the array is empty, then no documents where deleted.
// int equals their old ID
func (c *Document) Remove() []interface{} {
	var results []*Document

	if c.ID.Hex() == "" {
		return documentsToInterfaces(results)
	}

	ses, con, err := dbsession.GetCollection(Collection)
	if err != nil {
		return documentsToInterfaces(results)
	}
	defer ses.Close()

	err = con.Find(bson.M{
		"_id": c.ID,
	}).All(&results)

	if err != nil {
		fmt.Println(err.Error())
	}

	return documentsToInterfaces(results)
}

// Find returns an empty array when no match was found
func (c *Document) Find() []interface{} {
	var results []*Document

	ses, con, err := dbsession.GetCollection(Collection)
	if err != nil {
		return documentsToInterfaces(results)
	}
	defer ses.Close()

	// If no _id is set, use the content and assume we are invoking a document
	if c.ID.Hex() == "" {
		err = con.Find(bson.M{
			"date": c.Date,
		}).All(&results)
	} else {
		err = con.Find(bson.M{
			"_id": c.ID,
		}).All(&results)
	}

	if err != nil {
		fmt.Println(err.Error())
	}

	return documentsToInterfaces(results)
}

// make sure struct implements interface
var _ document.Interface = (*Document)(nil)
