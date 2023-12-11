docker build -t registry.mobilisis.com/user/docs/mysql .
docker push registry.mobilisis.com/user/docs/mysql
docker save -o mysql.tar registry.mobilisis.com/user/docs/mysql

timeout /t 30
