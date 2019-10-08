# kube-network-flow
Traffic flow chart from NetworkPolicy

## Example

```plantuml
@startuml component
actor client
node app
database db

db -> app
app -> client
@enduml
```