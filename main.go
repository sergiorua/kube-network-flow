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
	"context"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	v1 "k8s.io/api/networking/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/Masterminds/sprig"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

var verbose bool
var kubeconfig string
var templFile string
var namespace string
var policyName string

const umlTemplate = `
@startuml {{.Name}}

start

if (direction?) then (ingress)
	fork
{{ range $_, $v := .Spec.Ingress -}}
{{ range $_, $f := $v.From -}}
{{ if $f.IPBlock -}}
	{{ (print ":" $f.IPBlock.CIDR) }};
	floating note left: IP Block
	fork again
{{- end -}}
{{if $f.NamespaceSelector -}}
{{ range $index, $label := $f.NamespaceSelector.MatchLabels -}}
    {{ (print ":" $label) | indent 4 }};
    floating note left: {{$index}}
	fork again
{{- end -}}
{{- end }}
{{if $f.PodSelector -}}
{{ range $index, $label := $f.PodSelector.MatchLabels -}}
	{{ (print ":" $label) | indent 4 }};
    floating note left: {{$index}}
	fork again
{{ end }}
{{- end }}
{{- end }}
{{ end }}
  end fork
else (egress)
{{ range $_, $v := .Spec.Egress -}}
{{ if $v }}
	fork
{{ range $_, $f := $v.To -}}
{{ if $f.IPBlock -}}
	{{ (print ":" $f.IPBlock.CIDR) }};
	floating note left: IP Block
	fork again
{{- end -}}
{{if $f.NamespaceSelector -}}
{{ range $index, $label := $f.NamespaceSelector.MatchLabels -}}
    {{ (print ":" $label) | indent 4 }};
    floating note left: {{$index}}
	fork again
{{ end }}
{{ end }}
{{if $f.PodSelector -}}
{{ range $index, $label := $f.PodSelector.MatchLabels -}}
	{{ (print ":" $label) | indent 4 }};
    floating note left: {{$index}}
	fork again
{{- end }}
{{- end }}
{{- end }}
{{ end }}
  end fork
{{ end }}
endif

@enduml
`

func init() {
	flag.BoolVar(&verbose, "v", false, "Verbose")

	if os.Getenv("KUBECONFIG") != "" {
		flag.StringVar(&kubeconfig, "kubeconfig", os.Getenv("KUBECONFIG"), "(optional) absolute path to the kubeconfig file")
	} else if home := homeDir(); home != "" {
		flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.StringVar(&templFile, "template", "", "absolute path to the template file")
	flag.StringVar(&namespace, "namespace", "", "Limit to just this namespace (default all)")
	flag.StringVar(&policyName, "policy", "", "Limit to just this policy (default all)")
}

// RenderUml returns a UML diagram using http://plantuml.com/
func RenderUml(templ string, networkPolicy v1.NetworkPolicy) {
	tmpl, err := template.New("Policies").Funcs(sprig.FuncMap()).Parse(templ)
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
	var templ string

	if templFile != "" {
		dat, err := ioutil.ReadFile(templFile)
		if err != nil {
			panic(err)
		}
		templ = string(dat)
	} else {
		templ = umlTemplate
	}

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err.Error())
	}

	filter := metav1.ListOptions{}
	if policyName != "" {
		filter = metav1.ListOptions{
			FieldSelector: fields.OneTermEqualSelector("metadata.name", policyName).String(),
		}
	}

	policies, err := clientset.NetworkingV1().NetworkPolicies(namespace).List(context.TODO(), filter)
	for _, pol := range policies.Items {
		RenderUml(templ, pol)
	}

}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
