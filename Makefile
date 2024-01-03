# PLATFORMS := darwin/amd64 darwin/arm64 freebsd/amd64 linux/386 linux/amd64 linux/arm linux/arm64 linux/mipsle windows/386 windows/amd64 windows/arm windows/arm64
PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/386 linux/arm linux/arm64 linux/mipsle  linux/riscv64 windows/amd64 windows/arm64
temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

all:  web.tar.gz $(PLATFORMS) fixwindows wdos_file_checksum.sha1

binary: $(PLATFORMS)

hash: wdos_file_checksum.sha1

web: web.tar.gz

clean: 
	rm -f wdos_*_*
	rm -f web.tar.gz

$(PLATFORMS):
	@echo "Building $(os)/$(arch)"
	GOROOT_FINAL=Git/ GOOS=$(os) GOARCH=$(arch) GOARM=6 go build -o './dist/wdos_$(os)_$(arch)'  -ldflags "-s -w" -trimpath

fixwindows:
	-mv ./dist/wdos_windows_amd64 ./dist/wdos_windows_amd64.exe
	-mv ./dist/wdos_windows_arm64 ./dist/wdos_windows_arm64.exe

web.tar.gz:
	
	@echo "Removing old build resources, if exists"
	-rm -rf ./dist/
	-mkdir ./dist/

	@echo "Moving subfolders to build folder"
	-cp -r ./app ./dist/app/
	# cp -r ./subservice ./dist/subservice/
	-cp -r ./system ./dist/system/

	@ echo "Remove sensitive information"
	-rm ./dist/system/dev.uuid
	-rm ./dist/system/ao.db
	-rm ./dist/system/smtp_conf.json
	-rm ./dist/system/storage.json
	# -mv ./dist/system/storage.json ./dist/system/storage.json.example
	-rm -rf ./dist/system/aecron/
	-rm -rf ./dist/system/logs/
	-rm -rf ./dist/system/storage/
	-rm ./dist/system/cron.json
	-rm ./dist/system/bridge.json
	-rm ./dist/system/auth/authlog.db
	
	@ echo "Remove ip2country, todo: Add binary base searching for this function"
	-rm -rf "./dist/system/ip2country"

	@ echo "Remove modules that should not go into the build folder"
	-rm -rf "./dist/app/Cyinput"
	-rm -rf "./dist/app/Label Maker"
	-rm -rf "./dist/app/Dummy"
	

	@echo "Creating tarball for all required files"
	-rm ./dist/web.tar.gz
	-cd ./dist/ && tar -czf ./web.tar.gz system/ web/

	@echo "Clearing up unzipped folder structures"
	-rm -rf "./dist/web"
	-rm -rf "./dist/system"

wdos_file_checksum.sha1:
	@echo "Generating the checksum, if sha1sum installed"
	-sha1sum ./dist/web.tar.gz > ./dist/wdos_file_checksum.sha1