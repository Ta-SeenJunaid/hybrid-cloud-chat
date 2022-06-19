# hybrid-cloud-chat

## Docker Build & Push Image to Docker Hub

- To build docker image & push it to Docker Hub:

    ```bash
    $ cd chat-app

    # docker build -t DOCKER_REGESTRY/chat-app .
    $ docker build -t taseenjunaid/chat-app .

    # docker push -t DOCKER_REGESTRY/chat-app:latest .
    $ docker push taseenjunaid/chat-app:latest
    ```


Update `deploy-local.yaml` and `deploy-remote.yaml` file with your image (such as `taseenjunaid/chat-app:latest`).

## Prepare Remote Cluster (AWS EKS) 

- Create Cluster:
    ```bash
    $ eksctl create cluster --name remote-cluster
    ```
- Update cluster config with awscli:
    ```bash
    $ aws eks update-kubeconfig --region eu-central-1  --name remote-cluster
    ```
- Install NATS:
    ```bash
    $ helm repo add nats https://nats-io.github.io/k8s/helm/charts/
    $ helm repo update
    $ helm install my-nats nats/nats
    ```
- Go to K8s folder from project directory
    ```bash
    $ cd K8s
    ```

- Expose NATS: Create load-balancer type service:
    ```bash
    $ kubectl apply -f nats.yaml
    ```

- Get NATS public address: `LoadBalancer Ingress` 
   
    ```bash
    $ kubectl describe services nats-lb
    ```
    Output:
    ```
    Name:                     nats-lb
    Namespace:                default
    Labels:                   <none>
    Annotations:              <none>
    Selector:                 app.kubernetes.io/name=nats
    Type:                     LoadBalancer
    IP Family Policy:         SingleStack
    IP Families:              IPv4
    IP:                       10.100.202.198
    IPs:                      10.100.202.198
    LoadBalancer Ingress:     a7d62f596d6eb4e5dbc5e6b87ac63192-1340131761.eu-central-1.elb.amazonaws.com
    Port:                     nats  4222/TCP
    TargetPort:               4222/TCP
    NodePort:                 nats  31083/TCP
    Endpoints:                192.168.77.234:4222
    Port:                     leafnodes  7422/TCP
    ```
Here, it is `a7d62f596d6eb4e5dbc5e6b87ac63192-1340131761.eu-central-1.elb.amazonaws.com`

- Update `deploy-remote.yaml`: Set `NAT_URL` env to the latest one which we get from NATS public address: `LoadBalancer Ingress`
- Deploy Chat-App to cluster: 
    ```bash
    $ kubectl apply -f deploy-remote.yaml
    ```
- Describe the `chat-app-remote` service to get public address: `LoadBalancer Ingress`

    ```bash
    $ kubectl describe services chat-app-remote
    ```
    Output:
    ```
    Name:                     chat-app-remote
    Namespace:                default
    Labels:                   <none>
    Annotations:              <none>
    Selector:                 app=chat-app-remote
    Type:                     LoadBalancer
    IP Family Policy:         SingleStack
    IP Families:              IPv4
    IP:                       10.100.133.76
    IPs:                      10.100.133.76
    LoadBalancer Ingress:     a30ebc0f526294f889093f141841c3b6-1399523946.eu-central-1.elb.amazonaws.com
    Port:                     <unset>  8000/TCP
    TargetPort:               8000/TCP

    ```
    Go to public address: `LoadBalancer Ingress` link with port. Here it is:  `a30ebc0f526294f889093f141841c3b6-1399523946.eu-central-1.elb.amazonaws.com:8000` link. Make sure to mention the port `8000`. 
- To get the pods, services and statefulsets:
     ```bash
    $ kubectl get all
    ```

## Prepare Local Cluster (Kind Cluster)
- Go to K8s folder from project directory if you are not there.
    ```bash
    $ cd K8s
    ```

- Create kind cluster with given config file (i.e. `kind-config.yaml`):
    ```bash
    $ kind create cluster --config=kind-config.yaml
    ```
- Update `deploy-local.yaml`: Set `NAT_URL` env to the latest one which we get from NATS public address: `LoadBalancer Ingress`
- Deploy the Chat-App to local cluster:
    ```bash
    $ kubectl apply -f deploy-local.yaml
    ```
- To get the pods, services and statefulsets:
     ```bash
    $ kubectl get all
    ``` 
- Now your app is accessible at localhost:30000 url.

You are all set to go. Happy Hybrid-Cloud chatting!!

## Delete Clusters
- To get all contexts:
     ```bash
    $ kubectl config get-contexts
    ``` 
- To go to remote context:
     ```bash
    $ kubectl config use-context <remote cluster context name>
    ```
- To delete remote cluster from remote context:
     ```bash
    $ eksctl delete cluster --name=remote-cluster
    ```
- To go to local context:
     ```bash
    $ kubectl config use-context <local cluster context name>
    ```
- To delete local cluster from local context:
     ```bash
    $ kind delete cluster
    ```


 