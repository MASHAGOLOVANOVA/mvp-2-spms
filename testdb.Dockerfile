FROM mysql:latest

ENV MYSQL_ROOT_PASSWORD=root
ENV MYSQL_DATABASE=student_project_management_test
ENV MYSQL_USER=owner
ENV MYSQL_PASSWORD=root

RUN rm -rf /docker-entrypoint-initdb.d/*
# Копируем SQL-файлы в папку инициализации
COPY database/migrations/testing/up/ /docker-entrypoint-initdb.d/