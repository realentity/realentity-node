# HTTPS Deployment Guide for RealEntity Node

This guide covers different approaches to deploy your RealEntity node with HTTPS support.

## Quick Start - Generate Development Certificates

For local development and testing:

### Windows
```bash
cd scripts
.\generate-tls-cert.bat
```

### Linux/macOS
```bash
cd scripts
chmod +x generate-tls-cert.sh
./generate-tls-cert.sh
```

Then update your `config.json`:
```json
{
  "server": {
    "bind_address": "0.0.0.0",
    "port": 4001,
    "http_port": 8080,
    "https_port": 8443,
    "tls_cert_file": "certs/server.crt",
    "tls_key_file": "certs/server.key"
  }
}
```

## Production Deployment Options

### Option 1: Built-in TLS Support (Recommended for VPS)

1. **Obtain SSL certificates:**
   - **Let's Encrypt (Free)**: Use `certbot` to get certificates
   - **Commercial CA**: Purchase certificates from providers like Cloudflare, DigiCert
   - **Self-signed**: For testing only (not recommended for production)

2. **Let's Encrypt Example:**
```bash
# Install certbot
sudo apt update
sudo apt install certbot

# Get certificate for your domain
sudo certbot certonly --standalone -d yourdomain.com

# Certificates will be in /etc/letsencrypt/live/yourdomain.com/
```

3. **Update config.json:**
```json
{
  "server": {
    "bind_address": "0.0.0.0",
    "port": 4001,
    "http_port": 8080,
    "https_port": 443,
    "tls_cert_file": "/etc/letsencrypt/live/yourdomain.com/fullchain.pem",
    "tls_key_file": "/etc/letsencrypt/live/yourdomain.com/privkey.pem",
    "public_ip": "YOUR_VPS_IP"
  }
}
```

4. **Deploy:**
```bash
./deploy/universal.sh vps-bootstrap --public-ip YOUR_VPS_IP
```

### Option 2: Reverse Proxy (nginx/Cloudflare)

This is the most common production approach.

#### Using nginx as reverse proxy:

1. **Install nginx:**
```bash
sudo apt update
sudo apt install nginx
```

2. **Configure nginx (`/etc/nginx/sites-available/realentity`):**
```nginx
server {
    listen 80;
    server_name yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    
    # Modern SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-SHA384;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;

    # Proxy to RealEntity node
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

3. **Enable the site:**
```bash
sudo ln -s /etc/nginx/sites-available/realentity /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

4. **Configure RealEntity to only listen on HTTP (behind proxy):**
```json
{
  "server": {
    "bind_address": "127.0.0.1",
    "port": 4001,
    "http_port": 8080,
    "https_port": 0,
    "public_ip": "YOUR_VPS_IP"
  }
}
```

### Option 3: Cloudflare Tunnel (Zero Trust)

1. **Install cloudflared:**
```bash
wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
sudo dpkg -i cloudflared-linux-amd64.deb
```

2. **Authenticate with Cloudflare:**
```bash
cloudflared tunnel login
```

3. **Create tunnel:**
```bash
cloudflared tunnel create realentity-node
```

4. **Configure tunnel (`~/.cloudflared/config.yml`):**
```yaml
tunnel: realentity-node
credentials-file: /home/ubuntu/.cloudflared/your-tunnel-id.json

ingress:
  - hostname: api.yourdomain.com
    service: http://localhost:8080
  - service: http_status:404
```

5. **Route DNS:**
```bash
cloudflared tunnel route dns realentity-node api.yourdomain.com
```

6. **Run tunnel:**
```bash
cloudflared tunnel run realentity-node
```

### Option 4: Docker with TLS

1. **Update Dockerfile to include certificates:**
```dockerfile
# Add to your Dockerfile
COPY certs/ /app/certs/
```

2. **Docker run with TLS:**
```bash
docker run -d \
  --name realentity-node \
  -p 8080:8080 \
  -p 8443:8443 \
  -v $(pwd)/certs:/app/certs \
  -v $(pwd)/config.json:/app/config.json \
  realentity-node:latest
```

### Option 5: Kubernetes with Cert-Manager

1. **Install cert-manager:**
```bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml
```

2. **Create ClusterIssuer:**
```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: your-email@domain.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
```

3. **Create Ingress with TLS:**
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: realentity-ingress
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  tls:
  - hosts:
    - api.yourdomain.com
    secretName: realentity-tls
  rules:
  - host: api.yourdomain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: realentity-service
            port:
              number: 8080
```

## Testing HTTPS

After setting up HTTPS, test your endpoints:

```bash
# Health check
curl -k https://yourdomain.com/health

# Node info
curl -k https://yourdomain.com/api/node

# With proper certificate (remove -k)
curl https://yourdomain.com/health
```

## Security Best Practices

1. **Use strong TLS configuration:**
   - TLS 1.2+ only
   - Strong cipher suites
   - Perfect Forward Secrecy

2. **Regular certificate renewal:**
   - Set up automatic renewal for Let's Encrypt
   - Monitor certificate expiration

3. **Firewall configuration:**
   - Only open necessary ports (443, 4001)
   - Block direct access to HTTP port if using reverse proxy

4. **Headers for security:**
   - HSTS (HTTP Strict Transport Security)
   - CSP (Content Security Policy)
   - X-Frame-Options

## Troubleshooting

### Common Issues:

1. **Certificate not found:**
   - Check file paths in config.json
   - Verify file permissions

2. **Port binding issues:**
   - Ensure ports aren't already in use
   - Check firewall rules

3. **Certificate validation errors:**
   - Verify domain matches certificate
   - Check certificate chain

4. **Let's Encrypt rate limits:**
   - Use staging environment for testing
   - Wait for rate limit reset

### Debugging Commands:

```bash
# Check certificate details
openssl x509 -in certs/server.crt -text -noout

# Test TLS connection
openssl s_client -connect yourdomain.com:443

# Check port availability
netstat -tlnp | grep :443
```

## Performance Considerations

1. **TLS overhead:** HTTPS adds ~1-2% CPU overhead
2. **Connection reuse:** Enable HTTP/2 and keep-alive
3. **Certificate caching:** Use session resumption
4. **Load balancing:** Distribute TLS termination across multiple servers

Choose the deployment method that best fits your infrastructure and security requirements. For most production deployments, Option 2 (reverse proxy) or Option 3 (Cloudflare) are recommended.
