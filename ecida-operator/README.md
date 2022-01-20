# ECiDA Kubernetes operator

This operator is instantiated with `operator-sdk`. operator-sdk gives you a
comprehensive Makefile with several commands. The most important commands are

## !! Important !!

Before running any of the following commands, make sure that `kubectl` points
to a cluster that is intended to be used for ECiDA development, as the operator
will deploy itself onto this cluster.

## Requirements

- The operator requires access to a Kubernetes cluster. It should be sufficient
  to have  `kubectl` set up to point to a cluster through the `~/.kube/config`
  file. A cluster that is set up with minikube will work just fine.

- Go needs to be installed. According to `go.mod` at least version 1.15 needs
  to be installed. I am currently using go 1.16.9. I have not tested if there
  is a specific version that needs to be used for this to work.

- `gnumake` is required to run the `make` commands from the `Makefile` that are
  outlined below

## Commands

- ```
  make run
  ```

  This compiles the operator, and starts it. The operator runs locally and
  attempts to connect to whichever cluster is currently configured in
  `~/.kube/config`. This means that any Resources that exist on the cluster
  that this operator interacts with, will be read by the operator.
  Corresponding pods will be deployed on the cluster as well if the resources
  indicate so. For this command to be successfully ran, make sure to run `make
  install` first.

- ```
  make install
  ```

  This installs all the necessary CustomResourceDefinitions onto the cluster.
  **This command should be ran before running the operator.** To undo this
  command, you can run `make uninstall`.

- ```
  make deploy
  ```
  
  This deploys the operator as a pod on the cluster. Make sure to run `make
  install` first to ensure that the CRDs are available. To undo this command,
  run `make undeploy`.

- ```
  make generate
  ```

  The CustomResourceDefinitions are built from structs in `./api/v1alpha1/`,
  for example the `pipelinemodule_types.go` file. When you make changes in
  these files, make sure to run this command to update the
  CustomResourceDefinitions. The compiled output of this command for the
  `pipelinemodule_types.go` file can be seen in
  `./config/crd/bases/ecida.researchable.nl_pipelinemodules.yaml`.

- ```
  make docker-build
  make docker-push
  ```
  
  These commands are used to create and push docker images for the operator.
  The full name of the image is defined inside of the `Makefile` on line 32.
  Currently this is set to `maxverbeek/ecida-operator:v0.1`. For these commands
  to work, obvisouly `docker` needs to be installed.
