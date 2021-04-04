countwords:
	go run cmd/countwords/main.go < cmd/countwords/words.txt

test:
	gotestsum .