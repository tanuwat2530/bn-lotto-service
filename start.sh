rm -rf bn.log bn-lotto-app
go get
go build -o bn-lotto-app
./bn-lotto-app > bn.log 2>&1 &
