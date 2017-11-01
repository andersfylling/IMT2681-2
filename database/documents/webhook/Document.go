package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andersfylling/IMT2681-2/database/dbsession"
	"github.com/andersfylling/IMT2681-2/database/documents/document"
	"gopkg.in/mgo.v2/bson"
)

type Document struct {
	ID      bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	URL     string        `json:"webhookURL"`
	Base    string        `json:"baseCurrency"`
	Target  string        `json:"targetCurrency"`
	Current float64       `json:"currentRate"`
	Min     float64       `json:"minTriggerValue"`
	Max     float64       `json:"maxTriggerValue"`
}

// Invoke For invoking discord webhooks
type Webhook struct {
	Content   string `json:"content"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

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
	return wh, err
}

// Inserts the document as a new one into the collection and returns the id
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
func (c *Document) Remove() []document.Interface {
	results := []document.Interface{}

	return results
}

// Find returns an empty array when no match was found
func (c *Document) Find() []interface{} {
	var results []*Document

	ses, con, err := dbsession.GetCollection(Collection)
	if err != nil {
		return documentsToInterfaces(results)
	}
	defer ses.Close()

	err = con.Find(bson.M{
		"base":    c.Base,
		"target":  c.Target,
		"current": c.Current,
		"min":     c.Min,
		"max":     c.Max,
	}).All(&results)

	fmt.Println(len(results))

	if err != nil {
		fmt.Println(err.Error())
	}

	return documentsToInterfaces(results)
}

func (c *Document) FindAndInvoke() []*Document {
	webhooks := c.Find()
	for _, hook := range webhooks {
		(hook.(*Document)).InvokeWebhook("test", "lol", "")
	}

	return interfaceToDocument(webhooks)
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
	fmt.Println(jsonStr)
	res, err := http.Post(c.URL, "application/json; charset=utf-8", jsonStr)

	if res.StatusCode == 200 {
		return nil
	}

	return err
}

// make sure struct implements interface
var _ document.Interface = (*Document)(nil)
