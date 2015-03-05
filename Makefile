SteamKit := $(wildcard ~/github/SteamRE/SteamKit/Resources/Protobufs)

default : proto

proto:
	git -C game-tracking pull
	rm -rf dota
	make dota

game-tracking:
	git init game-tracking
	cd game-tracking && \
	git remote add -f origin https://github.com/SteamDatabase/GameTracking && \
	git config core.sparseCheckout true && \
	echo Protobufs/dota/ >> .git/info/sparse-checkout && \
	echo Protobufs/dota_s2/ >> .git/info/sparse-checkout && \
	echo Protobufs/dota_test/ >> .git/info/sparse-checkout && \
	git pull --depth=1 origin master

dota: dota/google/protobuf/descriptor.pb.go game-tracking
	rm -f dota/*.proto
	cp game-tracking/Protobufs/dota/*.proto dota/
	sed -i 's/^\(\s*\)\(optional\|repeated\|required\)\s*\./\1\2 /' dota/*.proto
	sed -i 's!^\s*rpc\s*\(\S*\)\s*(\.\([^)]*\))\s*returns\s*(\.\([^)]*\))\s*{!rpc \1 (\2) returns (\3) {!' dota/*.proto
	sed -i '1ipackage dota;\n' dota/*.proto
	protoc -I$(SteamKit)/ -Idota --go_out=dota dota/*.proto
	sed -i 's|google/protobuf/descriptor.pb|github.com/dotabuff/yasha/dota/google/protobuf|' dota/*.pb.go

dota/google/protobuf/descriptor.pb.go : descriptor.proto
	mkdir -p dota/google/protobuf
	protoc -I. --go_out=dota/google/protobuf $<
