docker build -t registry.mobilisis.com/elena.krzina/docs/mysql .
docker push registry.mobilisis.com/elena.krzina/docs/mysql
docker save -o mysql.tar registry.mobilisis.com/elena.krzina/docs/mysql

timeout /t 30