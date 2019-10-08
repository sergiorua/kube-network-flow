# kube-network-flow
Traffic flow chart from NetworkPolicy

## Example

@startuml component
actor client
node app
database db

db -> app
app -> client
@enduml