build:buildsp builddfz builddjd builddsc builddsl builddzz clear
buildsp:
	-rm sp.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o sp.exe ./sp/cmd/boot.go
	-upx -9 sp.exe
builddfz:
	-rm dfz.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dfz.exe ./daifazheng/cmd/boot.go
	-upx -9 dfz.exe
builddjd:
	-rm djd.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o djd.exe ./daijueding/cmd/boot.go
	-upx -9 djd.exe
builddsc:
	-rm dsc.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dsc.exe ./daishencha/cmd/boot.go
	-upx -9 dsc.exe
builddsl:
	-rm dsl.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dsl.exe ./daisouli/cmd/boot.go
	-upx -9 dsl.exe
builddzz:
	-rm dzz.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dzz.exe ./daizhizheng/cmd/boot.go
	-upx -9 dzz.exe
clear:
	-rm *.upx