package magazine

import (
	"emarket/internal/records"
	"sort"
)

type Model struct {
	dataPath string
	byid     map[string]*records.Magazine
	sorted   []*records.Magazine
}

func NewModel(dataPath string) (*Model, error) {
	magazs, err := loadEnabledMagazines(dataPath)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(magazs, func(i, j int) bool {
		return magazs[i].Title < magazs[j].Title
	})

	return &Model{
		dataPath: dataPath,
		byid:     buildMapByID(magazs),
		sorted:   magazs,
	}, nil
}

func (m *Model) GetRecord(id string) *records.Magazine {
	return m.byid[id]
}

func (m *Model) GetSortedRecords() []*records.Magazine {
	return m.sorted
}

func loadEnabledMagazines(dataPath string) ([]*records.Magazine, error) {
	all := make([]*records.Magazine, 0)
	if err := loadDataFromJsonFile(dataPath, &all); err != nil {
		return nil, err
	}

	filtered := make([]*records.Magazine, 0)
	for _, m := range all {
		if m.Enable {
			filtered = append(filtered, m)
		}
	}

	return filtered, nil
}

func buildMapByID(magazs []*records.Magazine) map[string]*records.Magazine {
	byid := make(map[string]*records.Magazine)
	for _, m := range magazs {
		byid[m.ID] = m
	}
	return byid
}
