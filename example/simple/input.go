//go:generate go run github.com/polpo-studio/goverter/cmd/goverter github.com/polpo-studio/goverter/example/simple
package simple

// goverter:converter
type Converter interface {
	Convert(source []Input) []Output
}

type Input struct {
	Name string
	Age  int
}

type Output struct {
	Name string
	Age  int
}
