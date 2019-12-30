# Birthday-server
Google Container Engine  based app

[![LICENSE](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Build Status](https://travis-ci.com/toddlers/birthday-server.svg?branch=master)](https://travis-ci.org/toddlers/birthday-server)

## The application


 * It's a web application to wish happy birthday or how many days left to your birthday
 * Takes inputs from the URL query form and return the result or error message in JSON format
 * Please seed the test data in tests directory.
 * in the `Makefile` there's a rudimentary example of spawning all of the service and connecting them, plus sending a test request

## Environment

 * the Go compiler is installed and the `go` command is on the path
 * `make` and anything else needed by the Makefile is installed and is on the path
 * `docker` is on the path and the user running `make` is in the `docker` group
 * a `gcloud login` was already performed so `docker push` just works
 * `kubectl` is on the path
 * a Kubernetes cluster is set up and the the default `kubectl` context points to that cluster. 
   We can create the cluster using below commands

   ```
   gcloud config set project my-project-id
   gcloud config set compute/zone us-central1-f
   gcloud container clusters create example --zone us-central1-a \
    --num-nodes 2 --enable-autoscaling --min-nodes 1 --max-nodes 3 \
    --node-locations us-central1-a,us-central1-b,us-central1-f
   gcloud container clusters get-credentials --cluster example
   ```

## Building the Docker Image
 
 * Do a git init in the code and set respective user and email address
  ```
  git config --global user.email "[EMAIL_ADDRESS]"
  git config --global user.name "[USERNAME]"
  ```
  
 * Make initial commits to the source code repository
 * Create a repository to host your code
 ```
  gcloud source repos create sample-app
  git config credential.helper gcloud.sh
 ```
 * Add your newly created repository as remote
 ```
 export PROJECT=$(gcloud info --format='value(config.project)')
 git remote add origin https://source.developers.google.com/p/$PROJECT/r/example
 ```
 * Push your code to the new repository's master branch
 ```
 git push origin master
 ```
 
 * Create a Tag for the code and push 
 ```
 git tag v1.0.0
 git push --tags
```

### Configuring the build triggers :
 
 * In GCP console click `Build Triggers` in the container registry seciotn
 * Select `Cloud Source Repository` and click `Continue`
 * Select example repository and select cloudbuild.yaml location: `gcp/cloudbuild.yaml`.
 * Select Trigger type as `Tag` and in the Tag regex provide `v.*`


## Build and Test

  * Build images using  `make docker-images DOCKER_REPO=gcr.io/birthday-server`
  * Create and test cluster using `make kubernetes-test`
  * Clean up after testing using `make kubernetes-clean`


## Deployment
 * Do the necessary code changes , tag them and push to the repo. 
   Once pushed we can check the auto build in the build history
 * We have three ways to do rolling update.
   - Using `set image`
   
     ```
     $ kubectl set image deployment <deployment> <container>=<image> --record
     # example
     git:mastkubectl set image deployment web-controller  web=example:v1.0.0 --record
     deployment "web-controller" image updated
     
     ```
     
  - Modify the kubernetes manifest and update the docker image version.
  
    ```
    # format
    $ kubectl replace -f <yaml> --record
    # example
    $ kubectl replace -f gcp/birthday-server.yml --record
    
    ```
   
  - Using `Edit`
  
    ```
    # format
    $ kubectl edit deployment <deployment> --record
    git:master *⚡ kubectl edit deployment web-controller --record
    
    ```
 * Check Status/Pause/Resume Deployments
 
   ```
   kubectl rollout status deployment web-controller # Check Status
   kubectl rollout pause deployment web-controller # Pause Rolling Update
   kubectl rollout resume deployment web-controller # Resume Rolling Update
   ```
 * To create and list the firewall rule that allows traffic into the cluster I use a gcloud compute command 
 
 ```
 gcloud compute firewall-rules create --allow=tcp:80 --target-tags=gke-example-default-pool-12d321682 example #creating the firewall rule
 gcloud compute forwarding-rules list

 ```


## Autoscaling
 
 * When you use `kubectl autoscale`, you specify a maximum and minimum number of replicas for your        application, as well as a CPU utilization target. 
 * For example, to set the maximum number of replicas to 4 and the minimum to 3, with a CPU utilization target of 50% utilization, run the following command:
   
```
git:master *⚡ kubectl  get deployments  # getting deployments list
NAME             DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
web-controller   2         2        2            2           1h

git:master *%=⚡ kubectl autoscale deployment web-controller --max 4 --min 3 --cpu-percent 50
deployment "web-controller" autoscaled

```


  ## Rollback
  ### to previous revision
  
  ```
   $ kubectl rollout undo deployment web-controller
   git:master *⚡ kubectl rollout undo deployment web-controller
   deployment "web-controller" rolled back
 ```
  ### to specific revision
  
 ```
   $ kubectl rollout undo deployment web-controller --to-revision=<revision>
   # exmaple
   $ kubectl rollout undo deployment web-controller --to-revision=v1.0.1
 
 ```
