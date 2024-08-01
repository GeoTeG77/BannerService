package models 

type Content struct {
	Title *string `json:"title"`
	Text  *string `json:"text"`
	URL   *string `json:"url" validate:"required,url"`
}

type CreateRequest struct {
	Tag_ids    []int64  `json:"tag_ids" validate:"required"`
	Feature_id *int64   `json:"feature_id" validate:"required"`
	Content    *Content `json:"content"`
	Is_active  *bool    `json:"is_active" validate:"required"`
}

type GetBannerRequest struct {
	Tag_id           *int64
	Feature_id       *int64
	Use_last_version *bool
}

type GetBannersRequest struct {
	Tag_id     *int64
	Feature_id *int64
	Limit      *int64
	Offset     *int64
}

type GetBannersResponce struct {
	Banner_id  *int64   `json:"banner_id" validate:"required"`
	Tag_ids    []int64  `json:"tag_ids" validate:"required"`
	Feature_id *int64   `json:"feature_id,omitempty"`
	Content    *Content `json:"content,omitempty"`
	Is_active  *bool    `json:"is_active" validate:"required"`
	Created_at *string  `json:"created_at"`
	Updated_at *string  `json:"upadated_at"`
}

type Response struct {
	//resp.Response
	Status    *string `json:"status"`
	Error     *string `json:"error,omitempty"`
	Banner_id *string `json:"alias,omitempty"`
}