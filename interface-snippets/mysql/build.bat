docker build -t mysql-img .
docker save -o mysql.tar mysql-img

timeout /t 30
