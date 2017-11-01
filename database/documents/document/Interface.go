package document

import "gopkg.in/mgo.v2/bson"

type Interface interface {
	Insert() (bson.ObjectId, error)
	Save() (old, new Interface)
	Remove() map[int]Interface
	Find() map[int]Interface
}
