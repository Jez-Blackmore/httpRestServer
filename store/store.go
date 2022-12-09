package store

import (
	"fmt"
	"time"
)

type Key string

type StoreMain struct {
	key           map[string]StructValueObject
	putChannel    chan StructValueObjectPut
	deleteChannel chan Key
	/* 	readsChannel  chan FormattedStructValueObject
	   	writesChannel chan FormattedStructValueObject */
}

type StructValueObject struct {
	Key     string        `json:"key"`
	Value   string        `json:"value"`
	Owner   string        `json:"owner"`
	Writes  int           `json:"writes"`
	Reads   int           `json:"reads"`
	Age     time.Duration `json:"age"`
	Updated time.Time     `json:"updated"`
}

type StructValueObjectPut struct {
	Value  string `json:"value"`
	Owner  string `json:"owner"`
	key    Key
	Writes int `json:"writes"`
}

type FormattedStructValueObject struct {
	Key     string        `json:"key"`
	Owner   string        `json:"owner"`
	Writes  int           `json:"writes"`
	Reads   int           `json:"reads"`
	Age     time.Duration `json:"age"`
	Updated time.Time     `json:"updated"`
}

var (
	MainStoreMain  StoreMain
	StoreFormatted map[string]FormattedStructValueObject
)

func NewStoreMain() StoreMain {

	Store := StoreMain{

		key:           map[string]StructValueObject{},
		putChannel:    make(chan StructValueObjectPut),
		deleteChannel: make(chan Key),
		/* readsChannel:  make(chan FormattedStructValueObject),
		writesChannel: make(chan FormattedStructValueObject), */
	}

	return Store
}

func (s *StoreMain) Monitor() {
	for {
		select {
		case put := <-s.putChannel:

			var keyToShow string
			var keyOwner string

			for keyVal, val := range s.key {
				if keyVal == string(put.key) {
					keyToShow = string(keyVal)
					keyOwner = val.Owner
					fmt.Println("test :", val)
				}
			}
			if keyToShow == "" {

				s.key[string(put.key)] = StructValueObject{Key: string(put.key), Value: put.Value, Owner: put.Owner, Writes: 1, Reads: 0, Updated: time.Now(), Age: 0}

				s.putChannel <- StructValueObjectPut{Value: s.key[string(put.key)].Value}

			} else if keyToShow != "" && keyOwner == put.Owner || put.Owner == "admin" {
				if thisProduct, ok := s.key[keyToShow]; ok {
					thisProduct.Value = put.Value
					/* fmt.Println("old: ", s.key[keyToShow].Writes, "new: ", s.key[keyToShow].Writes+1) */

					/* thisProduct.Age = time.Since(s.key[keyToShow].Updated)
					thisProduct.Updated = time.Now() */
					thisProduct.Writes = s.key[keyToShow].Writes + 1
					s.key[keyToShow] = thisProduct
				}

				/* 		s.readsChannel <- FormattedStructValueObject{Key: string(put.key)} */

				s.putChannel <- StructValueObjectPut{Value: s.key[string(put.key)].Value}
			} else {
				s.putChannel <- StructValueObjectPut{Value: ""}
			}

		case deleteVal := <-s.deleteChannel:
			delete(s.key, string(deleteVal))
			s.deleteChannel <- deleteVal
			fmt.Printf("Deleted %v\n", string(deleteVal))
			/* case readVal := <-s.readsChannel:
			if thisProduct, ok := s.key[readVal.Key]; ok {
				thisProduct.Age = time.Since(s.key[readVal.Key].Updated)
				thisProduct.Updated = time.Now()
				thisProduct.Reads = s.key[readVal.Key].Reads + 1
				s.key[readVal.Key] = thisProduct
			} */

			/* 		s.readsChannel <- FormattedStructValueObject{Key: s.key[readVal.Key].Key, Owner: s.key[readVal.Key].Owner, Writes: s.key[readVal.Key].Writes, Reads: s.key[readVal.Key].Reads, Age: s.key[readVal.Key].Age} */
			/* } */

		}
	}
}

// Helper functions
// Helper functions
func (s *StoreMain) GetKeyValueOwner(key string) string {

	var owner string

	for keyVal, _ := range s.key {
		if string(keyVal) == key {
			owner = s.key[keyVal].Owner
		}
	}

	return owner
}

// GET /store/<key>
func (s *StoreMain) UpdateStoreGet(key string) string {
	var valueToShow string

	for keyVal, value := range s.key {
		if string(keyVal) == key {
			valueToShow = value.Value

			if thisProduct, ok := s.key[value.Key]; ok {
				thisProduct.Age = time.Since(s.key[value.Key].Updated)
				thisProduct.Updated = time.Now()
				thisProduct.Reads = s.key[value.Key].Reads + 1
				s.key[value.Key] = thisProduct
			}
		}
	}

	return valueToShow
}

// DELETE /store/<key>

func (s *StoreMain) UpdateStoreDelete(key Key) bool {

	s.deleteChannel <- key

	value := <-s.deleteChannel

	if value != "" {
		return true
	}

	return false
}

// PUT /store/<key>

func (s *StoreMain) UpdateStorePut(key Key, valueToUpdate StructValueObject, username string) string {

	s.putChannel <- StructValueObjectPut{key: key, Value: valueToUpdate.Value, Owner: username}

	value := <-s.putChannel

	return value.Value
}

// GET /list
func (s *StoreMain) StoreListGet() []FormattedStructValueObject {

	/* var item FormattedStructValueObject */
	var items []FormattedStructValueObject

	var item FormattedStructValueObject

	for _, value := range s.key {

		/* s.readsChannel <- FormattedStructValueObject{Key: keyVal, Owner: value.Owner, Writes: value.Writes, Reads: value.Reads, Age: value.Age, Updated: value.Updated} */

		/* item := <-s.readsChannel */

		if thisProduct, ok := s.key[value.Key]; ok {
			thisProduct.Age = time.Since(s.key[value.Key].Updated)
			thisProduct.Updated = time.Now()
			thisProduct.Reads = s.key[value.Key].Reads + 1
			s.key[value.Key] = thisProduct
		}
		item = FormattedStructValueObject{Key: s.key[value.Key].Key, Owner: s.key[value.Key].Owner, Writes: s.key[value.Key].Writes, Reads: s.key[value.Key].Reads, Age: s.key[value.Key].Age}

		items = append(items, item)
	}

	return items
}

// GET /list/<key>

func (s *StoreMain) StoreListKeyGet(key string) FormattedStructValueObject {

	/* var keyToShow string
	var ownerToShow string
	var numWritesToShow int
	var numReadsToShow int
	var ageToShow time.Duration */

	var item FormattedStructValueObject

	for keyVal, value := range s.key {
		if string(keyVal) == key {

			if thisProduct, ok := s.key[value.Key]; ok {
				thisProduct.Age = time.Since(s.key[value.Key].Updated)
				thisProduct.Updated = time.Now()
				thisProduct.Reads = s.key[value.Key].Reads + 1
				s.key[value.Key] = thisProduct
			}
			item = FormattedStructValueObject{Key: s.key[value.Key].Key, Owner: s.key[value.Key].Owner, Writes: s.key[value.Key].Writes, Reads: s.key[value.Key].Reads, Age: s.key[value.Key].Age}

			/* keyToShow = string(keyVal)
			ownerToShow = value.Owner
			numWritesToShow = value.Writes
			numReadsToShow = value.Reads
			ageToShow = time.Since(value.Updated) */
			/* fmt.Printf("%s value is %v\n", keyVal, value) */
		}
	}

	return item
}
