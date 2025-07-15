# Kubernetes deployment for social-network

## Steps for local launch in Minikube

1. **Create namespace:**
   ```sh
   kubectl apply -f k8s/namespace.yaml
   ```
2. **Deploy the application and service:**
   ```sh
   kubectl apply -f k8s/deployment.yaml -n social-network
   kubectl apply -f k8s/service.yaml -n social-network
   ```
3. **Check status:**
   ```sh
   kubectl get pods -n social-network
   kubectl get svc -n social-network
   ```
4. **Open the service in your browser:**
   ```sh
   minikube service social-network -n social-network
   ```

## Notes
- The image used is public: `ghcr.io/olegmedvedev/social-network:latest`
- For production, it is recommended to use Ingress and secrets for environment variables.
- The healthcheck uses the `/healthz` endpoint (make sure it is implemented in your application). 