package main

import (
	"encoding/json"
)

type currencies map[string]float64

type state struct {
	Currencies currencies
}

func newState(current []byte) (*state, error) {
	var err error
	s := &state{}
	if len(current) > 0 {
		err = json.Unmarshal(current, s)
	}
	if s.Currencies == nil {
		s.Currencies = make(currencies)
	}
	return s, err
}

func (s state) set(currency string, value float64) {
	s.Currencies[currency] = value
}

func (s state) get(currency string) float64 {
	current, ok := s.Currencies[currency]
	if ok {
		return current
	}
	return 0
}

func (s state) toJSON() ([]byte, error) {
	return json.Marshal(s)
}
