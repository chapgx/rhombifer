


# Run test 
test TEST SUITE:
	go test -v -run "^Test{{TEST}}$" {{SUITE}}


# Run all tests
testall SUITE="./tests":
	go test -v {{SUITE}}

