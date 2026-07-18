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

### create .github/workflows/actions.yaml
>past the code bellow in the actions.yaml file
```bash
# CICD using GitHub actions

name: CI/CD  # name of workflow or pipeline

# Exclude the workflow to run on changes to the helm chart (avoids a trigger loop, since this workflow itself commits to helm/**)
on:
  push:                       # this workflow runs when someone pushes code
    branches:
      - main                    # only the main branch triggers this workflow
    paths-ignore:              # changes to these files/folders will NOT trigger the workflow
      - 'helm/**'                 # ignore helm chart changes (this workflow updates these itself, so ignoring avoids a loop)
      - 'k8s/**'                  # ignore raw kubernetes manifest changes
      - 'README.md'                # ignore documentation-only changes
      - 'CICD.md'                  # ignore documentation-only changes

jobs:                          # a workflow is made of one or more jobs, each running on its own fresh machine

  build:                         # job 1: compiles the Go app and runs tests
    runs-on: ubuntu-latest          # the machine (runner) this job executes on

    steps:                          # a job is a sequence of steps run in order
    - name: Checkout repository       # pulls your repo's code onto the runner so later steps can use it
      uses: actions/checkout@v4

    - name: Set up Go 1.22             # installs Go on the runner so we can build/test the app
      uses: actions/setup-go@v5          # updated from v2, current version, includes built-in module caching
      with:
        go-version: 1.22

    - name: Build                       # compiles the app into a binary named "go-app"
      run: go build -o go-app

    - name: Test                        # runs all Go tests in the project
      run: go test ./...

  code-quality:                  # job 2: checks code style/quality issues, runs in parallel with "build"
    runs-on: ubuntu-latest

    steps:
    - name: Checkout repository
      uses: actions/checkout@v4

    - name: Set up Go 1.22
      uses: actions/setup-go@v5          # added, golangci-lint needs Go available to resolve imports
      with:
        go-version: 1.22

    - name: Run golangci-lint           # scans the code for style issues, bugs, and bad patterns
      uses: golangci/golangci-lint-action@v6
      with:
        version: v1.56.2                  # pins a specific linter version so results stay consistent over time

  push:                           # job 3: builds a Docker image and pushes it to DockerHub
    runs-on: ubuntu-latest

    needs: [build, code-quality]    # waits for both "build" and "code-quality" to succeed first (a failing lint blocks the image push)

    steps:
    - name: Checkout repository
      uses: actions/checkout@v4

    - name: Set up Docker Buildx        # sets up Docker's modern build engine (needed for advanced build features)
      uses: docker/setup-buildx-action@v1

    - name: Login to DockerHub          # authenticates with DockerHub so we're allowed to push images
      uses: docker/login-action@v3
      with:
        username: ${{ secrets.DOCKERHUB_USERNAME }}   # stored securely in repo secrets, never hardcoded
        password: ${{ secrets.DOCKERHUB_TOKEN }}

    - name: Build and Push action       # builds the Docker image from the Dockerfile and pushes it
      uses: docker/build-push-action@v6
      with:
        context: .                        # build context = current directory (where Dockerfile expects files)
        file: ./Dockerfile                 # path to the Dockerfile to use
        push: true                          # actually push the image to DockerHub (not just build it locally)
        tags: |                             # image tags to publish
          ${{ secrets.DOCKERHUB_USERNAME }}/go-app:${{ github.run_id }}   # unique tag per run, useful for tracking exact deployments
          ${{ secrets.DOCKERHUB_USERNAME }}/go-app:latest                   # also tag as "latest" for convenience

  update-newtag-in-helm-chart:   # job 4: updates the Helm chart with the new image tag and commits it back to the repo
    runs-on: ubuntu-latest

    needs: push                     # only runs after the Docker image has been successfully pushed

    steps:
    - name: Checkout repository
      uses: actions/checkout@v4
      with:
        token: ${{ secrets.GH_PAT }}       # renamed from TOKEN, rename the actual repo secret to match; needs a personal access token (not the default token) to push commits back

    - name: Update tag in Helm chart    # replaces the "tag:" line in values.yaml with this run's image tag
      run: |
        sed -i 's/tag: .*/tag: "${{ github.run_id }}"/' helm/go-app-chart/values.yaml

    - name: Commit and push changes     # commits the updated Helm chart back to the repo
      run: |
        git config --global user.email "oneilkimbi1@gmail.com"     # git needs an identity to make commits
        git config --global user.name "Oneil Kimbi"
        git add helm/go-app-chart/values.yaml
        git diff --staged --quiet && echo "No changes to commit" || (git commit -m "Update tag in Helm chart" && git push)
        # this line avoids the job failing if there's genuinely nothing new to commit
```
> settings>secrets and variables>actions> newrepository secrets

>name*
>DOCKERHUB_USERNAME
>secret*
>oneil6677 # replace with your dockerhub username

### Add a new repository secret for dockerhub token
>name*
>DOCKERHUB_TOKEN

>secret*
>Paste your token here

### How to get dockerhub token
>log int to dockerhub> click on `myprofile` on top right> select `account settings`> click on `personal access token`> click on `Generate new token`> give a description then on second box give read and write permisions


### Add a new repository secret for github token



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
