package webhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/andersfylling/IMT2681-2/database/dbsession"
	"github.com/andersfylling/IMT2681-2/database/documents/document"
	"github.com/andersfylling/IMT2681-2/utils"
	"gopkg.in/mgo.v2/bson"
)

// Document is a valid MongoDB document that can be used for database operands
type Document struct {
	ID      bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	URL     string        `json:"webhookURL"`
	Base    string        `json:"baseCurrency"`
	Target  string        `json:"targetCurrency"`
	Current float64       `json:"currentRate,omitempty"`
	Min     float64       `json:"minTriggerValue"`
	Max     float64       `json:"maxTriggerValue"`
}

// Webhook For invoking discord webhooks
type Webhook struct {
	Content   string `json:"content"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

// Collection Which collection this document resides in
const Collection = "Webhook"

// New Creates a new instance of the document.
// Which can then be saved, removed, find matches, etc.
func New() *Document {
	return &Document{
		//ID:     nil,
		URL:    "",
		Base:   "",
		Target: "",
		Min:    0.0,
		Max:    0.0,
	}
}

// NewFromRequest Uses a http.Request object to populate the Document from body content
func NewFromRequest(r *http.Request) (*Document, error) {
	decoder := json.NewDecoder(r.Body)
	wh := New()
	err := decoder.Decode(wh)

	// validate the webhook url
	if err == nil && wh.URL != "" && !utils.ValidURL(wh.URL) {
		err = errors.New("The given URL is not valid: " + wh.URL)
	}

	return wh, err
}

// Insert the document as a new one into the collection and returns the id
func (c *Document) Insert() (id bson.ObjectId, err error) {
	id = ""

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

	err = con.Remove(bson.M{"_id": c.ID})
	if err != nil {
		// reset array
		if len(results) > 0 {
			results = results[:0]
		}
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
			"base":    c.Base,
			"target":  c.Target,
			"current": c.Current,
			"min":     c.Min,
			"max":     c.Max,
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

// FindAndInvoke ..
func (c *Document) FindAndInvoke() []*Document {
	var results []*Document

	ses, con, err := dbsession.GetCollection(Collection)
	if err != nil {
		return results
	}
	defer ses.Close()

	err = con.Find(bson.M{
		"base":   c.Base,
		"target": c.Target,
		"$or": []bson.M{
			bson.M{"min": bson.M{"$lt": c.Current}},
			bson.M{"max": bson.M{"$gt": c.Current}},
		},
	}).All(&results)

	if err != nil {
		fmt.Println(err.Error())
	}
	for _, hook := range results {
		hook.InvokeWebhook("Triggered by value "+strconv.FormatFloat(c.Current, 'f', -1, 64), "Invoked by request", "")
	}

	return results
}

// InvokeWebhook invokes the webhook based on data from the Document
func (c *Document) InvokeWebhook(content, username, avatarURL string) error {
	body := &Webhook{}

	if content != "" {
		body.Content = content
	}
	if username != "" {
		body.Username = username
	}
	if avatarURL != "" {
		body.AvatarURL = avatarURL
	}

	jsonStr := new(bytes.Buffer)
	json.NewEncoder(jsonStr).Encode(body)
	res, err := http.Post(c.URL, "application/json; charset=utf-8", jsonStr)

	if res.StatusCode == 200 {
		return nil
	}

	return err
}

// InvokeAll invokes all stored webhooks
func InvokeAll() ([]*Document, error) {
	var results []*Document

	ses, con, err := dbsession.GetCollection(Collection)
	if err != nil {
		return results, err
	}
	defer ses.Close()

	err = con.Find(nil).All(&results)
	for _, doc := range results {
		doc.InvokeWebhook("Invoked for evaluation testing", "", "")
	}
	return results, err
}

// make sure struct implements interface
var _ document.Interface = (*Document)(nil)
