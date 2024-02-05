docker build --no-cache -t myapp:1.0.0 .
docker save -o directdioled.tar myapp:1.0.0
timeout /t 30
