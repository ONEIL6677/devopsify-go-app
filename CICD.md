## Helm

### What is Helm

Helm is a package manager for Kubernetes. It bundles YAML manifests (Deployments, Services, Ingress, etc.) into reusable packages called **charts**, so you can install, upgrade, and manage complex applications with a single command instead of applying many YAML files by hand.

### Why It's Important

- Avoids repeating and manually managing dozens of YAML files
- Lets you parameterize configs (`values.yaml`) instead of hardcoding
- Makes upgrades and rollbacks of entire applications simple
- Provides access to a huge ecosystem of pre-built charts (databases, monitoring tools, ingress controllers, etc.)

### How It's Used

- Install a chart to deploy an app to your cluster
- Override default settings via a `values.yaml` file or `--set` flags
- Upgrade a release when you change values or chart version
- Roll back to a previous release if something breaks

### Installation on Ubuntu

>Download and run Helm's official install script:
```bash
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

>Verify installation:
```bash
helm version
```

>Add a chart repository (example: Bitnami):
```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
```

>Update repo index:
```bash
helm repo update
```

>Search for a chart:
```bash
helm search repo bitnami
```

>Install a chart:
```bash
helm install my-release bitnami/nginx
```

>List installed releases:
```bash
helm list
```

>Uninstall a release:
```bash
helm uninstall my-release
```

## Helm chart creation

> Follow the commands bellow

> Create a help folder
```bash
mkdir helm
```
```bash
cd helm
```
> Create helm chart for this project
```bash
help create go-app-chart
```
```bash
cd go-app-chart
```
>within the helm folder you should see `go-app-chart` folder if you move in the folder you will see the follwing configuration files and folders
**Chart.yaml charts, templates, values.yaml**

> Switch to templates folder
```bash
cd templates
```
>the delelte everything in the templates folder
```bash
rm -rf *
```
> then copy all your k8s manifests and paste in the templates folder
```bash
cp ../../../k8s/manifests/* .
```
>verify if they were copied
```bash
ls
```

>delete all resources you create manualy earlier and use a single helm command to create all also monitor the image tage whic will be update dynamically in helm's values.yaml file because it was made to update dynamically inside the deployment.yaml file
```bash
kubectl delete ing go-app
```
```bash
kubectl delete deploy go-app
```
```bash
kubectl delete svc go-app
```

>Now use the command bellow and create resources the monitor image tage in the values.yaml file it will be dynamically update

> first go back to help folder
```bash
cd ../../..
```
>Now create app
```bash
helm install go-app ./go-app-chart
```
>Monitor the values.yaml file image tage update on your docker hub account should be same in this file

> Use the commands bellow to check if everything was created
```bash
kubectl get ing go-app
```
```bash
kubectl get deploy go-app
```
```bash
kubectl get svc go-app
```

> uninstall everything
```bash
helm uninstall go-app
```
> veriify if everything was removed
```bash
kubectl get all
```
>it should return nothing about go-app

# CICD
## CI stages (Github Actions)
- stage 1: Build and unit test
- stage 2: Static code analyses
- stage 3: Docker image build and push
- stage 4: update help with docker image create

## CD Stages (ArgoCD)















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
