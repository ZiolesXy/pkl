package respons

type User struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Role Role `json:"role"`
	Barangs []Barang `json:"barangs"`
}