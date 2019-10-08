/*
Copyright 2019 Sergio Rua

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	v1 "k8s.io/api/networking/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/Masterminds/sprig"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var verbose bool
var local bool
var kubeconfig string

const umlTemplate = `
{{ $first := true }}
@startuml {{.Name}}

start

if (direction?) then (ingress)
{{ range $_, $v := .Spec.Ingress -}}
{{ if $first }}
fork
{{ else }}
fork again
{{ end }}
{{ $first = false }}
{{ range $_, $f := $v.From -}}
{{if $f.NamespaceSelector -}}
{{ range $index, $label := $f.NamespaceSelector.MatchLabels -}}
    {{ (print ":" $label) | indent 4 }};
    floating note left: {{$index}}
{{ end }}
{{- end }}
{{if $f.PodSelector -}}
{{ range $index, $label := $f.PodSelector.MatchLabels -}}
	{{ (print ":" $label) | indent 4 }};
    floating note left: {{$index}}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
  end fork
else (egress)
  :Allow All;
endif

@enduml
`

func init() {
	flag.BoolVar(&verbose, "v", false, "Verbose")
	flag.BoolVar(&local, "l", false, "Run outside kube cluster (dev purposes)")

	if home := homeDir(); home != "" {
		flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	}
}

// RenderUml returns a UML diagram using http://plantuml.com/
func RenderUml(networkPolicy v1.NetworkPolicy) {
	tmpl, err := template.New("Policies").Funcs(sprig.FuncMap()).Parse(umlTemplate)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, networkPolicy)
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	var config *rest.Config
	var err error

	if local == false {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err.Error())
	}

	policies, err := clientset.NetworkingV1().NetworkPolicies("").List(metav1.ListOptions{})
	for _, pol := range policies.Items {
		RenderUml(pol)
		for _, ing := range pol.Spec.Ingress {
			for _, fr := range ing.From {
				if fr.PodSelector == nil {
					continue
				}
				for k, v := range fr.PodSelector.MatchLabels {
					fmt.Printf("%v = %v\n", k, v)
				}
			}
		}
	}

}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
