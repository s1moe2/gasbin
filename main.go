package main

func main() {
	s := NewAPIServer("0.0.0.0", 6789)
	s.SetRoutes(GetAPI())
	s.Start()
}
