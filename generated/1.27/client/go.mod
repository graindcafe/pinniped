// This go.mod file is generated by ./hack/update.sh.
module go.pinniped.dev/generated/1.27/client

go 1.13

replace go.pinniped.dev/generated/1.27/apis => ../apis

require (
	go.pinniped.dev/generated/1.27/apis v0.0.0
	k8s.io/apimachinery v0.27.16
	k8s.io/client-go v0.27.16
	k8s.io/kube-openapi v0.0.0-20230501164219-8b0f38b5fd1f
)
