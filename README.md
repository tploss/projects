# Projects
Program for managing your projects (repositories) on your local system.

## `projects` program
This repository is used to develop the `projects` CLI.
It is developed in Go 1.18.

## Project file
Example:
```yaml
projects:
  - name: biosim
    repos:
      - name: sim
        url: git@gitlab.com:tploss/biosim.git
  - name: utils
    repos:
      - name: godirserver
        url: git@gitlab.com:tploss/godirserver.git
      - name: ssc
        url: git@gitlab.com:tploss/script-snippet-collection.git
```
All `name` attributes are used for paths so only letters, digits, - and _ are allowed.
