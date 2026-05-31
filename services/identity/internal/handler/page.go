package handler

import identityv1 "github.com/Mond1c/lms/gen/go/lms/v1"

// Offset-based pagination: page_token is the row offset. nextPageToken advances
// by the page size while a full page was returned, and is 0 once exhausted.
const (
	defaultPageSize int32 = 50
	maxPageSize     int32 = 200
)

func pageParams(p *identityv1.PageRequest) (limit, offset int32) {
	size := p.GetPageSize()
	if size <= 0 {
		size = defaultPageSize
	}
	if size > maxPageSize {
		size = maxPageSize
	}
	offset = p.GetPageToken()
	if offset < 0 {
		offset = 0
	}
	return size, offset
}

func nextPageToken(offset, limit, returned int32) int32 {
	if returned < limit {
		return 0
	}
	return offset + limit
}
