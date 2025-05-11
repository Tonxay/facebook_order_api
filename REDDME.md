// docker buiold .
//docker tag go-api-web tonxay/go-api-web create tag for deploy images
docker build -t go-fiber-docker-api .
docker tag go-fiber-docker-api:latest your-dockerhub-username/go-fiber-docker-api:latest
docker login
docker push your-dockerhub-username/go-fiber-docker-api:latest
GOOS=linux GOARCH=amd64 go build -o server ./cmd/heroku
docker build -t tonxay/gowebhook:latest .
docker build --no-cache -t tonxay/gowebhook:latest .

Full Guide: Nginx + WebSocket + HTTPS on DigitalOcean

Full Guide: Nginx + WebSocket + HTTPS on DigitalOcean
✅ 1. Install Nginx on your Droplet

sudo apt update
sudo apt install nginx
✅ 2. Install Certbot for HTTPS (Let's Encrypt)

sudo apt install certbot python3-certbot-nginx
✅ 3. Set Up Nginx for WebSocket Reverse Proxy
Let's assume your Go WebSocket server runs on port 8080.

Edit your Nginx config:

sudo nano /etc/nginx/sites-available/default
Or create a new one in /etc/nginx/sites-available/myapp and symlink it.

Replace the server block with:

server {
listen 80;
server_name yourdomain.com;

    location /ws/ {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
    }

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

}
Replace yourdomain.com with your actual domain.

Then test and reload:

sudo nginx -t
sudo systemctl reload nginx
✅ 4. Enable HTTPS with Certbot
Run:
sudo certbot --nginx -d yourdomain.com
Certbot will automatically configure SSL in your Nginx config. It adds a new server block with listen 443 ssl;.

Now your client can use:js

const ws = new WebSocket("wss://yourdomain.com/ws/9392341827487393");
✅ WebSocket upgrade works securely.

✅ 5. Make Sure Your Go Server Binds to localhost:8080 or 0.0.0.0:8080
Since Nginx is reverse proxying to port 8080, your Go app must listen there:

http.ListenAndServe(":8080", nil)
✅ 6. Open Firewall Ports (Optional)
Make sure UFW allows HTTP/HTTPS:

sudo ufw allow 'Nginx Full'
sudo ufw enable
Would you like me to generate a working Nginx config file or Docker-compatible setup for this?
