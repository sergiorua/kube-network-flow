# kube-network-flow
Traffic flow chart from NetworkPolicy

## Usage

```sh
Usage of ./kube-network-flow:
  -kubeconfig string
        (optional) absolute path to the kubeconfig file (default "/Users/srua/.kube/config")
  -template string
        absolute path to the template file
  -v    Verbose
```

The command execution returns a list of all the network policies in UML format. There is a default [PlantUML](http://plantuml.com/) template similar to the file `net_template.plantuml.tmpl` this repository but you can create your own using Golang templates.

## Display

There are many ways to display the results, see this website: http://plantuml.com/running

## Sample

![Sample Output UML](https://github.com/sergiorua/kube-network-flow/raw/master/sample-output.png)