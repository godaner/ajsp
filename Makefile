build:builddfz builddjd builddsc builddsl builddzz
builddfz:
	-rm dfz.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dfz.exe ./daifazheng/boot.go
	-upx -9 dfz.exe
builddjd:
	-rm djd.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o djd.exe ./daijueding/boot.go
	-upx -9 djd.exe
builddsc:
	-rm dsc.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dsc.exe ./daishencha/boot.go
	-upx -9 dsc.exe
builddsl:
	-rm dsl.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dsl.exe ./daisouli/boot.go
	-upx -9 dsl.exe
builddzz:
	-rm dzz.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dzz.exe ./daizhizheng/boot.go
	-upx -9 dzz.exe
