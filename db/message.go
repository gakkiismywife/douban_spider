package db

type Message struct {
	*BaseModel
	Title string ` json:"title"`
	Url   string ` json:"url"`
}

func HasSend(title, url string) bool {

	message := &Message{}
	where := &Message{Title: title, Url: url}
	Db.Where(where).Find(message)

	return message.Url == url
}

func CreateMessage(title, url string) {
	m := &Message{Title: title, Url: url}
	Db.Create(m)
}
