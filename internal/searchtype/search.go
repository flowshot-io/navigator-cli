package searchtype

import (
	"encoding/json"
	"errors"
)

type SearchType int32

const (
	Unknown SearchType = iota
	ID
	TEXT
	Image
)

var (
	searchTypeNames = map[SearchType]string{
		Unknown: "unknown",
		ID:      "id",
		TEXT:    "text",
		Image:   "image",
	}

	ErrInvalidType = errors.New("invalid file status")
)

func (s SearchType) String() string {
	if name, ok := searchTypeNames[s]; ok {
		return name
	}

	return searchTypeNames[0]
}

func (s SearchType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *SearchType) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	*s = Unknown
	for k, v := range searchTypeNames {
		if v == str {
			*s = k
			return nil
		}
	}

	return ErrInvalidType
}

func (s SearchType) IsValid() bool {
	_, ok := searchTypeNames[s]
	return ok
}

func FromString(str string) (SearchType, error) {
	for k, v := range searchTypeNames {
		if v == str {
			return k, nil
		}
	}

	return Unknown, ErrInvalidType
}
