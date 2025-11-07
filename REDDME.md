// docker buiold .
//docker tag go-api-web tonxay/go-api-web create tag for deploy images
docker build -t go-fiber-docker-api .
docker tag go-fiber-docker-api:latest your-dockerhub-username/go-fiber-docker-api:latest
docker login
docker push your-dockerhub-username/go-fiber-docker-api:latest
GOOS=linux GOARCH=amd64 go build -o server ./cmd/heroku
docker build -t tonxay/gowebhook:latest .
docker build --no-cache -t tonxay/gowebhook:latest .

heroku buildpacks:clear

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

set Domain

sudo mkdir -p /var/www/api.chat-dd.uk
echo '{"status":"ok"}' | sudo tee /var/www/api.chat-dd.uk/index.json
sudo nano /etc/nginx/sites-available/api.chat-dd.uk
sudo apt install certbot python3-certbot-nginx -y
sudo certbot --nginx -d api.chat-dd.uk

server {
listen 80;
server_name api.chat-dd.uk;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

}


docker exec golang_app ls -la /go-api/images/bills/
docker exec golang_app rm /go-api/images/bills/8252705319817.jpg
docker exec golang_app rm -f /go-api/images/bills/*




docker system df -v

Your output from the df command shows that your root filesystem (/dev/vda1 mounted on /) is 95% full. This is a critical issue that needs immediate attention to prevent system instability.The high usage is directly tied to your Docker installation, as shown by the many overlay entries also reporting 95% usage.Here is a breakdown of the problem and the immediate steps you should take.🚨 Critical Disk Space IssueFilesystemTotal SizeUsed SpaceUse%Mount Point/dev/vda1~23.1GB~21.9GB95%/ (Root)overlay~23.1GB~21.9GB95%/var/lib/docker/overlay2/...The overlay entries are Docker's filesystem layers. Since they report the same usage as the root filesystem, it means the vast majority of your disk space is being consumed by Docker images, containers, volumes, and caches.🗑️ How to Free Up Space (Docker Cleanup)The fastest and most effective way to reclaim disk space is by using the Docker system prune command. This removes resources that are no longer actively running or referenced.1. Prune Unused Docker ResourcesThis command will remove all stopped containers, all dangling images (images without a tag), all unused networks, and all dangling build cache.Bashdocker system prune
It will ask for confirmation: Are you sure you want to continue? [y/N] Type y and press Enter.2. Prune Everything (Aggressive Cleanup)If docker system prune doesn't free up enough space, you can use the -a and -v flags for a more aggressive clean up, which will also remove:All unused images (not just dangling ones)All unused volumes (if you add the -v flag)⚠️ Warning: If you have data stored in a Docker Volume that is not currently attached to a container, this command will delete that data. Only run this if you are sure you don't need the volumes.