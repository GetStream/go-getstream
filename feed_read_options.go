package getstream

import "fmt"

type FeedReadOptions struct {
	limit  *int
	offset *int

	idGte string
	idGt  string
	idLte string
	idLt  string

	ranking string
}

func NewFeedReadOptions() FeedReadOptions {
	return FeedReadOptions{}
}

func (i FeedReadOptions) AddLimit(limit int) FeedReadOptions {
	i.limit = &limit
	return i
}

func (i FeedReadOptions) AddOffset(offset int) FeedReadOptions {
	i.offset = &offset
	return i
}

func (i FeedReadOptions) AddIdGte(idGTE string) FeedReadOptions {
	i.idGte = idGTE
	return i
}

func (i FeedReadOptions) AddIdGt(idGT string) FeedReadOptions {
	i.idGt = idGT
	return i
}

func (i FeedReadOptions) AddIdLte(idLTE string) FeedReadOptions {
	i.idLte = idLTE
	return i
}

func (i FeedReadOptions) AddIdLt(idLT string) FeedReadOptions {
	i.idLt = idLT
	return i
}

func (i FeedReadOptions) Params() (params map[string]string) {
	params = make(map[string]string)

	if i.limit != nil {
		params["limit"] = fmt.Sprintf("%d", *i.limit)
	}
	if i.offset != nil {
		params["offset"] = fmt.Sprintf("%d", *i.offset)
	}
	if i.idGte != "" {
		params["id_gte"] = i.idGte
	}
	if i.idGt != "" {
		params["id_gt"] = i.idGt
	}
	if i.idLte != "" {
		params["id_lte"] = i.idLte
	}
	if i.idLt != "" {
		params["id_lt"] = i.idLt
	}
	return params
}
