package models

type MessageModel struct {
	Sender    string `db:"sender_id"`
	Recipient string `db:"recipient_id"`
	Content   string `db:"message"`
}
