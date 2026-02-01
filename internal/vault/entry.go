package vault

type Entry struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Username string `json:"username"`
	Password []byte `json:"password"`
	URL      string `json:"url"`
	Notes    string `json:"notes"`
}
