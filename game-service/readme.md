# Book of Tomes

Public URL: https://gs.dev.topgaming.team/hello

### Deploy service on EKS

1. Configure EKS access as described here: https://topgaming.atlassian.net/wiki/spaces/TOPGAMING/pages/14712863/Dev+K8s
2. Push to `development` triggers the build. Track the progress here: https://eu-central-1.console.aws.amazon.com/codesuite/codepipeline/pipelines/dev-game-service-pipeline/view?region=eu-central-1  
        Navigate to "Details" link in "Build" section. Click on "Tail logs" button on the newly opened page to view the logs.
3. Once the build is complete a new docker image will be pushed to ECR repository (https://eu-central-1.console.aws.amazon.com/ecr/repositories/private/713614461671/game-service?region=eu-central-1)
4. Clone `eks` repository:  
        `git clone ssh://git-codecommit.eu-central-1.amazonaws.com/v1/repos/eks`
5. Navigate to terraform dev VPC directory:  
        `cd eks/dev/terraform`
6. Run `terraform init` to pull required modules locally
7. Find game service config located in `services/game-service.tf` and change image tag to the latest one in ECR (e.g. `fff4807`):  
        `image = "713614461671.dkr.ecr.eu-central-1.amazonaws.com/game-service:a6cf890"` -> `image = "713614461671.dkr.ecr.eu-central-1.amazonaws.com/game-service:fff4807"`
8. Run `terraform apply` to apply changes.
9. Check running pods with `kubectl get pods`  
        You should see something like this:  
        `NAME                           READY   STATUS    RESTARTS   AGE`  
        `d-store-bb9b76497-zgsh2        1/1     Running   0          101m`  
        `game-service-77bb64996-j4qxm   1/1     Running   0          44m`  
        `mock-casino-6f49f5c456-h4ghn   1/1     Running   0          3h17m`  
10. The latest version of the service is now available at `https://gs.dev.topgaming.team`

### Connecting to D-Store

D-Store API base URL: https://ds.dev.topgaming.team/v1  
Swagger specs: https://eu-central-1.console.aws.amazon.com/codesuite/codecommit/repositories/swagger-specs/browse/refs/heads/master/--/d-store/api.yaml?region=eu-central-1  
X-API-Key (for POST /round): `ae4a7fcbaa488b5fa004419d16a94ae2`
