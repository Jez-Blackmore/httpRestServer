package store

import "fmt"

/* type Key string */

type StructValueObject struct {
	Value string `json:"value"`
	Owner string `json:"owner"`
}

type FormattedStructValueObject struct {
	Key   string `json:"key"`
	Owner string `json:"owner"`
}

var (
	Store          map[string]StructValueObject
	StoreFormatted map[string]FormattedStructValueObject
)

func CreateStore() {

	Store = map[string]StructValueObject{"test": {Value: "1", Owner: "Jez"}}

	storeChannel := make(chan string)

	fmt.Println("Close store channel")
	close(storeChannel)
}

// Helper functions
func GetKeyValueOwner(key string) string {

	var owner string

	for keyVal, _ := range Store {
		if string(keyVal) == key {
			owner = Store[keyVal].Owner
		}
	}

	return owner
}

// GET /store/<key>
func UpdateStoreGet(key string) string {
	var valueToShow string

	for keyVal, value := range Store {
		if string(keyVal) == key {
			valueToShow = value.Value
		}
	}

	return valueToShow
}

// PUT /store/<key>
func UpdateStorePut(key string, valueToUpdate StructValueObject) string {
	var keyToShow string

	for keyVal, _ := range Store {
		if string(keyVal) == key {
			keyToShow = string(keyVal)
		}
	}

	if keyToShow == "" {

		Store[key] = StructValueObject{Value: valueToUpdate.Value, Owner: valueToUpdate.Owner}

	} else {

		if thisProduct, ok := Store[keyToShow]; ok {
			thisProduct.Value = valueToUpdate.Value
			Store[keyToShow] = thisProduct
		}

	}

	return Store[key].Value
}

// GET /list
func StoreListGet() []FormattedStructValueObject {

	var item FormattedStructValueObject
	var items []FormattedStructValueObject

	for keyVal, value := range Store {
		item = FormattedStructValueObject{Key: keyVal, Owner: value.Owner}
		items = append(items, item)
	}

	return items
}

// GET /list/<key>

func StoreListKeyGet(key string) (string, string) {

	var keyToShow string
	var valueToShow string

	for keyVal, value := range Store {
		if string(keyVal) == key {
			keyToShow = string(keyVal)
			valueToShow = value.Value
			fmt.Printf("%s value is %v\n", keyVal, value)
		}
	}

	return keyToShow, valueToShow
}
