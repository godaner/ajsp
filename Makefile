build:buildqrpj buildpj buildsp builddfz builddjd builddsc builddsl builddzz clear
buildpj:
	-rm ./bin/pj.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/pj.exe ./pj/cmd/boot.go
	-upx -9 ./bin/pj.exe
buildsp:
	-rm ./bin/sp.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/sp.exe ./sp/cmd/boot.go
	-upx -9 ./bin/sp.exe
builddfz:
	-rm ./bin/dfz.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/dfz.exe ./daifazheng/cmd/boot.go
	-upx -9 ./bin/dfz.exe
builddjd:
	-rm ./bin/djd.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/djd.exe ./daijueding/cmd/boot.go
	-upx -9 ./bin/djd.exe
builddsc:
	-rm ./bin/dsc.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/dsc.exe ./daishencha/cmd/boot.go
	-upx -9 ./bin/dsc.exe
builddsl:
	-rm ./bin/dsl.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/dsl.exe ./daisouli/cmd/boot.go
	-upx -9 ./bin/dsl.exe
builddzz:
	-rm ./bin/dzz.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/dzz.exe ./daizhizheng/cmd/boot.go
	-upx -9 ./bin/dzz.exe
buildqrpj:
	-rm ./bin/qrpj.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/qrpj.exe ./qrpj/cmd/boot.go
	-upx -9 ./bin/qrpj.exe
clear:
	-rm ./bin/*.upx