package respons

type User struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Role Role `json:"role,omitempty"`
	// Barangs []Barang `json:"barangs"`
}