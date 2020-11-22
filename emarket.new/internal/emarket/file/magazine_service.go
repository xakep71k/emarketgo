package file

import (
	"emarket/internal/emarket"
	"encoding/json"
	"io/ioutil"
	"sort"
)

func NewMagazineService(filename string) *MagazineService {
	return &MagazineService{filename: filename}
}

type MagazineService struct {
	filename string
}

func (r *MagazineService) Find() ([]*emarket.Magazine, error) {
	content, err := ioutil.ReadFile(r.filename)

	if err != nil {
		return nil, err
	}

	var magazines []*emarket.Magazine
	err = json.Unmarshal(content, &magazines)

	sort.SliceStable(magazines, func(i, j int) bool {
		return magazines[i].Title < magazines[j].Title
	})

	return magazines, err
}
