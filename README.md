# quest-server

### Как запустить локально
Выяснить, на каком ip докер видит локальную машину. Для этого выполнить команду
```
sudo ip addr show docker0
```
Вторая строка будет иметь вид
```
inet 172.17.0.1/16 brd 172.17.255.255 scope global docker0
```
Первый IP (в моем случае 172.17.0.1) - это <IP хоста>
Отредактировать файл /etc/postgresql/9.6/main/postgresql.conf, заменив в нем
```
listen_addresses = 'localhost'
```
на
```
listen_addresses = '*'
```
Отредактировать файл /etc/postgresql/9.6/main/pg_hba.conf, вставив первой строкой
```
host    all             all             <IP хоста>/0            md5
```
Поднять postgresql на порте, указанном в конфиге (по умолчанию поднимается на порте 5432).
Создать базу по схеме в папке resources.

Скачать репозиторий и перейти в его папку

Собрать image командой
```
docker build -t ard .
```
Запустить докер командой:
```
docker run -p <MACHINE PORT>:80 -v <DEMON CONFIG PATH>.json:/etc/ard.conf.json:ro -v <NGINX CONFIG PATH>:/etc/nginx/nginx.conf:ro -v <QUESTS DIR>:/ard/data/quests --add-host="outer_host:<IP хоста>" -t ard
```
Конфиг демона находится в resources/ard.conf.json. Конфиг nginx лежит в resources/nginx.conf. QUESTS DIR - папка с ресурсами квестов.

Предоплагается, что ресурсом квеста будет архив с файлами, необходимым для квеста. Имя архива - id квеста в базе.

