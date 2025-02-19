build-docker-images: ## собираем образы для будущих сборки и запуска приложения
	docker build --pull --rm -f build/build.Dockerfile -t ggghfffg/spams:build .

build-in-docker: ## запускаем контейнер с биндом к корневой папке хоста, в нем собираются приложения и сохраняются на хосте ## docker run --name docker-build-web-app-container -it -v ${PWD}:/app ggghfffg/spams:build
	docker run --name docker-build-web-app-container -it -v $(PWD):/mvp-2-spms ggghfffg/spams:build
	docker rm docker-build-web-app-container
	

run-commands-to-build-go: ## вызываются при запуске контейнера сборки приложения
	go mod download
	go build -o ./web_server/cmd/web_app/web_app ./web_server/cmd/web_app/

run-server: ## запуск сервера
	./web_server/cmd/web_app/web_app
