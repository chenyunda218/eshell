package client

type Doc struct {
	Name      string `json:"name"`
	Timestamp int64  `json:"timestamp"`
	Person    Person `json:"person"`
}

type Person struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Book   Book   `json:"book"`
	Gender string `json:"gender"`
}

type Book struct {
	Author string `json:"author"`
	Title  string `json:"title"`
}

type Index struct {
	Health       string `json:"health"`
	Status       string `json:"status"`
	Index        string `json:"index"`
	UUID         string `json:"uuid"`
	Pri          string `json:"pri"`
	Rep          string `json:"rep"`
	DocsCount    string `json:"docs.count"`
	DocsDeleted  string `json:"docs.deleted"`
	StoreSize    string `json:"store.size"`
	PriStoreSize string `json:"pri.store.size"`
}

type CreateSessionRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
