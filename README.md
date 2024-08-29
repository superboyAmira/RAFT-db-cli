# RAFT-db-cli
Fault-tolerant distributed local database system using the RAFT algorithm. Golang

## Tech Stack

* Go 1.22.3
* GRPC + protobuf

# Resources (links)

* go implemenation of RAFT protocol https://github.com/hashicorp/raft
* https://medium.com/@govinda.attal/raft-consensus-leader-election-with-golang-89bfdbd471cb

# About project

Как только кластер получает лидера, он может принимать новые записи журнала.
Клиент может запросить лидера добавить новую запись журнала,
которая представляет собой непрозрачный двоичный объект в Raft.
Затем лидер записывает запись в долговременное хранилище и пытается реплицировать ее
в кворум последователей. Как только запись журнала считается зафиксированной ,
ее можно применить к конечному автомату.
Конечный автомат зависит от приложения и реализуется с помощью интерфейса.

# Quick start

# Demonstration of the functionality