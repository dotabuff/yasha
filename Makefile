SteamKit := $(wildcard ~/github/SteamRE/SteamKit/Resources/Protobufs)

default : test

test:
	go test -cover -v

proto:
	rm -rf dota
	make dota

dota: dota/google/protobuf/descriptor.pb.go
	rm -f dota/*.proto
	cp game-tracking/Protobufs/dota/*.proto dota/
	sed -i 's/^\(\s*\)\(optional\|repeated\|required\)\s*\./\1\2 /' dota/*.proto
	sed -i 's!^\s*rpc\s*\(\S*\)\s*(\.\([^)]*\))\s*returns\s*(\.\([^)]*\))\s*{!rpc \1 (\2) returns (\3) {!' dota/*.proto
	sed -i '1ipackage dota;\n' dota/*.proto
	protoc -I$(SteamKit)/ -Idota --go_out=dota dota/*.proto
	sed -i 's|google/protobuf/descriptor.pb|github.com/dotabuff/yasha/dota/google/protobuf|' dota/*.pb.go

dota/google/protobuf/descriptor.pb.go : google/protobuf/descriptor.proto
	mkdir -p dota/google/protobuf
	protoc -I. --go_out=dota $<

sync-replays:
	s3cmd --region=us-west-2 sync ./replays/*.dem s3://yasha.dotabuff/
