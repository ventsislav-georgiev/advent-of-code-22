export SESSION_KEY=${skey}

.DEFAULT_GOAL := go
.PHONY: go
go:
	@go run golang/cmd/${year}/${day}/main.go --year=${year} --day=${day} --task=${task}

.PHONY: rust
rust:
	@cargo run -q -p aoc-${year}-${day} -- --year=${year} --day=${day} --task=${task}

bench:
	@go test -bench=. -benchtime=20s -benchmem -run=^$$ ./golang/cmd/day${day}
