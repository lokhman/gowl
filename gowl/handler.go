package gowl

type Handler func(r *Request) ResponseInterface

func EmptyHandler(_ *Request) ResponseInterface {
	return nil
}
