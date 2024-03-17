build:
	(cd src) && (go build -o ../bin main.go)
	copy ".env" "bin"