package domain

// UserUpdate carries a partial update to a User. A nil field means "leave
// unchanged"; a non-nil field is applied. For TelegramID, a non-nil pointer to
// the empty string clears (unsets) the value.
type UserUpdate struct {
	ID          string
	DisplayName *string
	TelegramID  *string
}
