package mydict

import "errors"

type Dictionary map[string]string

var (
	errNoKey = errors.New("No such key")
	errPreoccupiedKey = errors.New("The key already exists")
)

func (d Dictionary) Search(word string) (string, error) {
	value, exists := d[word]
	switch exists {
	case true:
		return value, nil
	}
	return "", errNoKey
}

func (d Dictionary) Add(word, def string) error {
	_, err := d.Search(word)
	switch err {
	case errNoKey:
		d[word] = def
		return nil
	}
	return errPreoccupiedKey
}

func (d Dictionary) Update(word, def string) error {
	_, err := d.Search(word)
	switch err {
	case nil:
		d[word] = def
	}
	return err
}

func (d Dictionary) Delete(word string) error {
	// or just delete(d, word)
	_, err := d.Search(word)
	switch err {
	case nil:
		d[word] = ""
	}
	return err
}