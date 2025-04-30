package lookatmego

type Plugin interface {
	Render(input []byte) string
}
