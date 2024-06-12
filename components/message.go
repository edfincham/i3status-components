package components

import (
	"encoding/json"
	"log"
)

type Message struct {
	Align               string `json:"align"`
	Border              string `json:"border"`
	BorderLeft          int    `json:"border_left"`
	BorderRight         int    `json:"border_right"`
	BorderTop           int    `json:"border_top"`
	BorderBottom        int    `json:"border_bottom"`
	Colour              string `json:"color"`
	FullText            string `json:"full_text"`
	Instance            string `json:"instance"`
	MinWidth            int    `json:"min_width"`
	Name                string `json:"name"`
	Urgent              bool   `json:"urgent"`
	Separator           bool   `json:"separator"`
	SeparatorBlockWidth int    `json:"separator_block_width"`
	ShortText           string `json:"short_text"`
}

func (m *Message) ToJson() string {
	s, err := json.Marshal(m)
	if err != nil {
		log.Fatal("failed to encode message")
	}

	return string(s)
}

func NewMessage() *Message {
	return &Message{
		Separator:           false,
		Align:               "left",
		SeparatorBlockWidth: 10,
	}
}
