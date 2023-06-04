# example2
----

## Creat example2
```shell
operator-sdk init --domain example --repo github.com/wangyuche/operator_example/example2
operator-sdk create api --group example --version v2 --kind Example --resource --controller
```

modify api/v1/example_type.go
```go
type ExampleStatus struct {
Time string `json:"time,omitempty"`
}
```
modify controllers/example_controller.go
```go
func (r *ExampleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
l := log.FromContext(ctx)
config := &examplev2.Example{}
err := r.Get(ctx, req.NamespacedName, config)
if err != nil {
l.Error(err, "get config err")
return ctrl.Result{}, err
}
config.Status.Time = time.Now().String().Format("2006-01-02 15:04:05")
r.Status().Update(ctx, config)
return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
}
```
build images
```shell
export images=arieswangdocker/operator_example2:latest
docker rmi -f ${images}
docker build . -t ${images}
docker push ${images}
cd config/manager && kustomize edit set image controller=${images}
cd ..
cd ..
kustomize build config/default > deploy.yaml
kubectl apply -f deploy.yaml
```
```shell
kubectl apply -f config/samples/example_v2_example.yaml
```
modify deploy.yaml
```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.11.1
  creationTimestamp: null
  name: examples.example.example
spec:
  group: example.example
  names:
    kind: Example
    listKind: ExampleList
    plural: examples
    singular: example
  scope: Namespaced
  versions:
  - name: v2
    additionalPrinterColumns:
      - name: Time
        type: string
        jsonPath: .status.time
...
```
```shell
kubectl apply -f deploy.yaml
```