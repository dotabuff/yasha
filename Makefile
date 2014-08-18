SteamKit := $(wildcard ~/github/SteamRE/SteamKit/Resources/Protobufs)

default : proto

proto: dota/google/protobuf/descriptor.pb.go
	rm -f dota/*.proto
	cp -r $(SteamKit)/dota/*.proto ./dota/
	rm -f dota/*.steamworkssdk.proto
	dos2unix -q dota/*.proto
	sed -i 's/^\(\s*\)\(optional\|repeated\|required\)\s*\./\1\2 /' dota/*.proto
	sed -i '1ipackage dota;\n' dota/*.proto
	protoc -I$(SteamKit) -Idota --gogo_out=dota dota/*.proto
	sed -i 's|google/protobuf|github.com/dotabuff/yasha/dota/google/protobuf|' dota/*.pb.go

dota/google/protobuf/descriptor.pb.go : $(SteamKit)/google/protobuf/descriptor.proto
	mkdir -p dota
	protoc -I$(SteamKit)/ --gogo_out=dota $<
