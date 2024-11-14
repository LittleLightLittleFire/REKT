package types

type CreateOutput struct {
	Data struct {
		ID   *string `json:"id"`
		Text *string `json:"text"`
	} `json:"data"`
}

func (r *CreateOutput) HasPartialError() bool {
	return false
}

type DeleteOutput struct {
	Data struct {
		Deleted *bool `json:"deleted"`
	} `json:"data"`
}

func (r *DeleteOutput) HasPartialError() bool {
	return false
}
