package fs

import (
	"emarket/internal/emarket"
	"encoding/json"
	"io/ioutil"
	"sort"
)

func NewMagazineStorage(filename string) *MagazineStorage {
	return &MagazineStorage{filename: filename}
}

type MagazineStorage struct {
	filename string
}

func (r *MagazineStorage) Find() ([]*emarket.Magazine, error) {
	content, err := ioutil.ReadFile(r.filename)

	if err != nil {
		return nil, err
	}

	var magazines []*emarket.Magazine
	err = json.Unmarshal(content, &magazines)

	magazines = filterDisabledMagaz(magazines)

	sort.SliceStable(magazines, func(i, j int) bool {
		return magazines[i].Title < magazines[j].Title
	})

	return magazines, err
}

func filterDisabledMagaz(mm []*emarket.Magazine) []*emarket.Magazine {
	var found int

	for _, m := range mm {
		if !m.Enable {
			continue
		}

		mm[found] = m
		found++
	}

	return mm[:found]
}
