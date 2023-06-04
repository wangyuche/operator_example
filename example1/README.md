# example1
----
Example 1 will let everyone know how to quickly create a PostgreSQL Operator through teaching


## Install Operator SDK
See documentation on [Operator SDK](https://sdk.operatorframework.io/docs/installation/).
```shell
brew install operator-sdk
```
## Creat PostgreSQL Operator
```shell
operator-sdk init --domain example --repo github.com/wangyuche/operator_example/example1
operator-sdk create api --group example --version v1 --kind Example --resource --controller
```
example yaml
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgresql-data
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: standard
  resources:
    requests:
      storage: 60Gi
---
apiVersion: v1
kind: Service
metadata:
  name: postgresql
spec:
  ports:
    - name: tcp
      port: 5432
      protocol: TCP
      targetPort: 5432
  selector:
    app: postgresql
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgresql
spec:
  serviceName: "postgresql"
  replicas: 1
  selector:
    matchLabels:
      app: postgresql
  template:
    metadata:
      labels:
        app: postgresql
    spec:
      containers:
      - image: postgres:13.11
        imagePullPolicy: Always
        name: postgresql
        ports:
        - containerPort: 5432
          protocol: TCP
        env:
        - name: POSTGRES_PASSWORD
          value: "sdffhgn"
        volumeMounts:
        - mountPath: /var/lib/postgresql
          name: postgresql-data
      terminationGracePeriodSeconds: 60
      volumes:
        - name: postgresql-data
          persistentVolumeClaim:
            claimName: postgresql-data
```
modify api/v1/example_type.go
```go
type ExampleSpec struct {
	Images   string `json:"images,omitempty"`
	Password string `json:"password,omitempty"`
}
```
modify controllers/example_controller.go
```go
func (r *ExampleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	config := &examplev1.Example{}
	pg_pvc := r.svcpg_standalone(config)
	pg_svc := r.svcpg_standalone(config)
	pg_sts := r.stspg_standalone(config)
	err := r.Create(ctx, pg_pvc)
	if err != nil {
		l.Error(err, "pg_pvc")
		return ctrl.Result{}, err
	}
	err = r.Create(ctx, pg_svc)
	if err != nil {
		l.Error(err, "pg_svc")
		return ctrl.Result{}, err
	}
	err = r.Create(ctx, pg_sts)
	if err != nil {
		l.Error(err, "pg_sts")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
```
add controllers/pd_standalone.go

build images
```shell
export images=arieswangdocker/operator_example1:latest
docker rmi -f ${images}
docker build . -t ${images}
docker push ${images}
cd config/manager && kustomize edit set image controller=${images}
cd ..
cd ..
kustomize build config/default > deploy.yaml
kubectl apply -f deploy.yaml
```
modify config/samples/example_v1_example.yaml
```yaml
apiVersion: example.example/v1
kind: Example
metadata:
  labels:
    app.kubernetes.io/name: example
    app.kubernetes.io/instance: example-sample
    app.kubernetes.io/part-of: example1
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/created-by: example1
  name: example-sample
spec:
  images: "postgres:13.11"
  password: "12345"
```
```shell
kubectl apply -f -<<EOF
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: example1-clusterrole
rules:
  - apiGroups: ["*"]
    resources: ["*"]
    verbs: ["create", "delete"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: example1-clusterrolebinding
subjects:
  - kind: ServiceAccount
    name: example1-controller-manager
    namespace: example1-system
roleRef:
  kind: ClusterRole
  name: example1-clusterrole
  apiGroup: rbac.authorization.k8s.io
EOF
```
```shell
kubectl apply -f config/samples/example_v1_example.yaml
```
