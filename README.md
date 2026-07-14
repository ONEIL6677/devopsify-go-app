# DEVOPSIFY GO APP FOR COMPLETE BEGINNERS
> if you are an experience perserson go and follow this readme /k8s/README.md

## Step 1: run the app locally




## step 2: Dockerise the app
> We will use Docker to containerize the Go web application. Docker is a container platform that allows you to > build, ship, and run containers.
> create a file named `Dockerfile` and past the following code inside

```yaml
# Containerize the go application that we have created
# This is the Dockerfile that we will use to build the image
# and run the container

# Start with a base image
FROM golang:1.22.5 AS base

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod .

# Download all the dependencies
RUN go mod download

# Copy the source code to the working directory
COPY . .

# Build the application
RUN go build -o main .

#######################################################
# Reduce the image size using multi-stage builds
# We will use a distroless image to run the application
FROM gcr.io/distroless/base

# Copy the binary from the previous stage
COPY --from=base /app/main .

# Copy the static files from the previous stage
COPY --from=base /app/static ./static

# Expose the port on which the application will run
EXPOSE 8080

# Command to run the application
CMD ["./main"]
```

>Commands to build the Docker container:
> NB replace oneil6677 with your docker hub username

```bash
docker build -t oneil6677/go-app:v1 .
```

>Command to run the Docker container:

```bash
docker run -p 8080:8080 oneil6677/go-app:v1
```

>Command to push the Docker container to Docker Hub:

```bash
docker push oneil6677/go-app"v1
```

# Step 3: Write kubernetes manifest
> in the cureent working directory create a folder called `k8s` then another folder inside called `manifest`

> inside the manifest foldder create the following files
### 1: A deployment file
> A Deployment in Kubernetes is an object that manages a set of identical pods for you it handles creating  them, keeping the right number running, and updating them safely over time.
```yaml
# This is a sample deployment manifest file for a simple web application.
apiVersion: apps/v1        # Kubernetes API version for Deployment objects
kind: Deployment            # Tells k8s we're creating a Deployment (manages pods)
metadata:
  name: go-app               # Name of this Deployment
  labels:
    app: go-app               # Label to identify/group this deployment
spec:
  replicas: 1                 # Number of pod instances to run
  selector:
    matchLabels:
      app: go-app               # Selects pods with this label to manage
  template:                   # Blueprint for the pods this Deployment creates
    metadata:
      labels:
        app: go-app               # Label applied to each pod (must match selector above)
    spec:
      containers:
      - name: go-app                # Name of the container inside the pod
        image: oneil6677/go-app:v1  # Docker image to run (dont forget to replace oneil6677 with your docker # # username)
        ports:
        - containerPort: 8080         # Port the app listens on inside the container
```

### 2: A service file
> A Service in Kubernetes is what exposes a set of pods to network traffic either to other pods inside the cluster, or to the outside world.
> Here's the problem it solves: pods are temporary. They get created, destroyed, and replaced constantly (especially by Deployments), and each new pod gets a new IP address. If something needed to talk directly to a pod's IP, that connection would break every time a pod restarted.

```yaml
# Service for the application
apiVersion: v1
kind: Service              # Exposes a set of pods under a stable network identity
metadata:
  name: go-app
  labels:
    app: go-app
spec:
  ports:
  - port: 80                 # Port exposed by the Service
    targetPort: 8080          # Port the container actually listens on
    protocol: TCP
  selector:
    app: go-app                # Routes traffic to pods with this label
  type: ClusterIP             # Internal-only; change to NodePort/LoadBalancer for external access
  ```

### An ingress file


> An Ingress is an object that manages external HTTP/HTTPS access to services inside your cluster think of it as a smart traffic router sitting in front of your Services.

> Here's why it exists: you could expose every Service externally using LoadBalancer, but that means spinning up a separate cloud load balancer (and cost) for every single service. That gets expensive and messy fast if you have many services.

```yaml
# Ingress resource for the application
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: go-app
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /   # rewrites request path before forwarding to service
spec:
  ingressClassName: nginx        # tells k8s to use the NGINX Ingress Controller
  rules:
  - host: go-app.local             # domain that routes to this service
    http:
      paths: 
      - path: /                      # match all paths under /
        pathType: Prefix
        backend:
          service:
            name: go-app                # must match your Service's name
            port:
              number: 80                  # must match your Service's exposed port
```


## Step 4: Creat a kubernetes cluster (EKS)

### prerequisites
#### instal kubectl
>kubectl – A command line tool for working with Kubernetes clusters. For more information, see [Installing or updating kubectl]("https://docs.aws.amazon.com/eks/latest/userguide/install-kubectl.html").

#### Install eksctl
>eksctl – A command line tool for working with EKS clusters that automates many individual tasks. For more information, see [Installing or updating]("https://docs.aws.amazon.com/eks/latest/userguide/eksctl.html").

>Download the kubectl binary for your cluster’s Kubernetes version from Amazon S3.
>Kubernetes 1.35
```bash
curl -O https://s3.us-west-2.amazonaws.com/amazon-eks/1.35.3/2026-04-08/bin/linux/amd64/kubectl
```
```bash
chmod +x ./kubectl
```
> copy binary to a folder in your path
```bash
mkdir -p $HOME/bin && cp ./kubectl $HOME/bin/kubectl && export PATH=$HOME/bin:$PATH
```
#### install aws cli
>AWS CLI – A command line tool for working with AWS services, including Amazon EKS. For more information, see [Installing, updating, and uninstalling the AWS CLI]("https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html") in the AWS Command Line Interface User Guide. After installing the AWS CLI, we recommend that you also configure it. For more information, see [Quick configuration]("https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html#cli-configure-quickstart-config") with aws configure in the AWS Command Line Interface User Guide.>

### Install EKS

> Install a EKS cluster with EKSCTL

```bash
eksctl create cluster --name demo-cluster --region us-east-1 
```

> Delete the cluster

```bash
eksctl delete cluster --name demo-cluster --region us-east-1
```































# DevOpsify the go web application

The main goal of this project is to implement DevOps practices in the Go web application. The project is a simple website written in Golang. It uses the `net/http` package to serve HTTP requests.

DevOps practices include the following:

- Creating Dockerfile (Multi-stage build)
- Containerization
- Continuous Integration (CI)
- Continuous Deployment (CD)

## Summary Diagram
![image](https://github.com/user-attachments/assets/45f4ef12-c5b5-4247-9d43-356b5dfb671b)


## Creating Dockerfile (Multi-stage build)

The Dockerfile is used to build a Docker image. The Docker image contains the Go web application and its dependencies. The Docker image is then used to create a Docker container.

We will use a Multi-stage build to create the Docker image. The Multi-stage build is a feature of Docker that allows you to use multiple build stages in a single Dockerfile. This will reduce the size of the final Docker image and also secure the image by removing unnecessary files and packages.

## Containerization

Containerization is the process of packaging an application and its dependencies into a container. The container is then run on a container platform such as Docker. Containerization allows you to run the application in a consistent environment, regardless of the underlying infrastructure.

We will use Docker to containerize the Go web application. Docker is a container platform that allows you to build, ship, and run containers.



## Continuous Integration (CI)

Continuous Integration (CI) is the practice of automating the integration of code changes into a shared repository. CI helps to catch bugs early in the development process and ensures that the code is always in a deployable state.

We will use GitHub Actions to implement CI for the Go web application. GitHub Actions is a feature of GitHub that allows you to automate workflows, such as building, testing, and deploying code.

The GitHub Actions workflow will run the following steps:

- Checkout the code from the repository
- Build the Docker image
- Run the Docker container
- Run tests

## Continuous Deployment (CD)

Continuous Deployment (CD) is the practice of automatically deploying code changes to a production environment. CD helps to reduce the time between code changes and deployment, allowing you to deliver new features and fixes to users faster.

We will use Argo CD to implement CD for the Go web application. Argo CD is a declarative, GitOps continuous delivery tool for Kubernetes. It allows you to deploy applications to Kubernetes clusters using Git as the source of truth.

The Argo CD application will deploy the Go web application to a Kubernetes cluster. The application will be automatically synced with the Git repository, ensuring that the application is always up to date.





