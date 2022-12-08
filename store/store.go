package store

import "fmt"

type Key string

type StoreMain struct {
	key           map[string]StructValueObject
	putChannel    chan StructValueObjectPut
	deleteChannel chan Key
}

type StructValueObject struct {
	Value string `json:"value"`
	Owner string `json:"owner"`
}

type StructValueObjectPut struct {
	Value string `json:"value"`
	Owner string `json:"owner"`
	key   Key
}

type FormattedStructValueObject struct {
	Key   string `json:"key"`
	Owner string `json:"owner"`
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
				s.key[string(put.key)] = StructValueObject{Value: put.Value, Owner: put.Owner}

				s.putChannel <- StructValueObjectPut{Value: s.key[string(put.key)].Value}

			} else if keyToShow != "" && keyOwner == put.Owner {
				if thisProduct, ok := s.key[keyToShow]; ok {
					thisProduct.Value = put.Value
					s.key[keyToShow] = thisProduct
				}

				s.putChannel <- StructValueObjectPut{Value: s.key[string(put.key)].Value}
			} else {
				s.putChannel <- StructValueObjectPut{Value: ""}
			}

		case deleteVal := <-s.deleteChannel:
			delete(s.key, string(deleteVal))
			s.deleteChannel <- deleteVal
			fmt.Printf("Deleted %v\n", string(deleteVal))
		}
	}
}

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
	/* delete(s.key, key) */
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

	var item FormattedStructValueObject
	var items []FormattedStructValueObject

	for keyVal, value := range s.key {
		item = FormattedStructValueObject{Key: keyVal, Owner: value.Owner}
		items = append(items, item)
	}

	return items
}

// GET /list/<key>

func (s *StoreMain) StoreListKeyGet(key string) (string, string) {

	var keyToShow string
	var ownerToShow string

	for keyVal, value := range s.key {
		if string(keyVal) == key {
			keyToShow = string(keyVal)
			ownerToShow = value.Owner
			fmt.Printf("%s value is %v\n", keyVal, value)
		}
	}

	return keyToShow, ownerToShow
}
