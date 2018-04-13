FROM golang

# Обновляем пакеты
RUN apt-get update

# Скачиваем исходный код в Docker-контейнер
ENV WORK /go/src
ADD . /go/src/github.com/Sovianum/arquest-server

# Переходим в директорию проекта
WORKDIR $WORK/github.com/Sovianum/arquest-server

# Устанавливаем dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Устанавливаем все зависимости
RUN dep ensure

# Устанавливаем nginx
RUN apt-get install -y nginx

# Собираем бинарь
RUN go build main.go

# Объявлем порт сервера
EXPOSE 80

CMD service nginx start && ./main -c /etc/ard.conf.json

#CMD /bin/bash