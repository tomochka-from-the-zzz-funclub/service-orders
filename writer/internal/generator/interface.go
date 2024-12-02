package generator

type Generate interface {
	GenerateOrder() (string, error)
}

//mockgen -source=internal/generator/interface.go -destination=internal/generator/mocks/mock_generate.go -package=mocks
