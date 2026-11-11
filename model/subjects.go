package model

type Subjects struct {
	ID     int    `json:"id"`
	UserId int    `json:"user_id"`
	Alg    string `json:"alg"`
	Bio    string `json:"bio"`
	Whis   string `json:"whis"`
	Geo    string `json:"geo"`
	Lang   string `json:"lang"`
	It     string `json:"it"`
	Khis   string `json:"khis"`
	Kaz    string `json:"kaz"`
	Rights string `json:"rights"`
	Lit    string `json:"lit"`
	Rus    string `json:"rus"`
	Phys   string `json:"phys"`
	Chem   string `json:"chem"`
}
