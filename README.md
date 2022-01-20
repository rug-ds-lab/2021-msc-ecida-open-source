# 2021-msc-ecida-open-source

This repository consists of two projects. The first one being a CLI
application, originally meant to set up a module repository in the form of a
Helm chart. This project can be found in `src`.

The `repo` directory contains several examples that have been made with the
CLI. These examples are actually just Helm charts. What is special about the
`repo` directory is that its contents get deployed to a Helm repository, hosted
on the github pages of this github repository (https://rug-ds-lab.github.io/2021-msc-ecida-open-source/index.yaml).
There is a github actions pipeline in the `.github` folder that makes this work.

The `ecida-operator` directory contains a Kubernetes operator instantiated with operator-sdk.
For more information about the operator, see the README.md in that directory.
