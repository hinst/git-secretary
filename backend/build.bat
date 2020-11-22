cd server
go build -o ../dist/ || exit /b -1
cd ..

cd story-girls-standard
go build -o ../dist/plugins/story-girls-standard.exe || exit /b -1
cd ..
