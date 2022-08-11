package model

type Config struct {
	ID        int    `json:"id"`
	Keyid     string `json:"keyid"`
	Keysecret string `json:"keysecret"`
	Ifileurl  string `json:"ifileurl"`
}

type Users struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Ifileuid int    `json:"ifileuid"`
}
type Project struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	IfileRoot string `json:"ifile_root" db:"ifile_root"`
	UID       int    `json:"uid"`
	UserID    int    `json:"userid"`
	Auth      string `json:"auth"`
}
type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	IfileRoot string `json:"ifile_root" db:"ifile_root"`
	UID       int    `json:"uid"`
}
