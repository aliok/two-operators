Commands to create two operators and install them in a cluster:
```shell
mkdir operator1
cd operator1
operator-sdk init --domain aliok --repo github.com/aliok/two-operators/operator1
operator-sdk create api --group test --version v1alpha1 --kind InstallationA --resource --controller
make docker-build docker-push IMG="docker.io/aliok/operator1:latest"
make bundle IMG="docker.io/aliok/operator1:latest"
make bundle-build bundle-push BUNDLE_IMG="docker.io/aliok/operator1-bundle:latest"
operator-sdk olm install
operator-sdk run bundle docker.io/aliok/operator1-bundle:latest
kubectl apply -f config/samples/test_v1alpha1_installationa.yaml

mkdir operator2
cd operator2
operator-sdk init --domain aliok --repo github.com/aliok/two-operators/operator2
operator-sdk create api --group test --version v1alpha1 --kind InstallationB --resource --controller
make docker-build docker-push IMG="docker.io/aliok/operator2:latest"
make bundle IMG="docker.io/aliok/operator2:latest"
make bundle-build bundle-push BUNDLE_IMG="docker.io/aliok/operator2-bundle:latest"
operator-sdk olm install
operator-sdk run bundle docker.io/aliok/operator2-bundle:latest
kubectl apply -f config/samples/test_v1alpha1_installationb.yaml
```

Commands to test OLM dependency resolution:
```shell
operator-sdk cleanup operator1
operator-sdk cleanup operator2

cd operator1
make docker-build docker-push IMG="docker.io/aliok/operator1:v0.0.1"
make bundle IMG="docker.io/aliok/operator1:v0.0.1"
make bundle-build bundle-push BUNDLE_IMG=docker.io/aliok/operator1-bundle:v0.0.1
make catalog-build catalog-push CATALOG_IMG=docker.io/aliok/operator1-catalog:v0.0.1
make bundle-build bundle-push catalog-build catalog-push

kubectl apply -f - <<EOF
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  # Replace the version with whatever version you want to install
  name: operator1-catalog
  namespace: default
spec:
  displayName: operator1-catalog
  image: docker.io/aliok/operator1-catalog:v0.0.1
  publisher: aliok
  sourceType: grpc
EOF

cd operator2
make docker-build docker-push IMG="docker.io/aliok/operator2:v0.0.1"
make bundle IMG="docker.io/aliok/operator2:v0.0.1"
make bundle-build bundle-push BUNDLE_IMG=docker.io/aliok/operator2-bundle:v0.0.1
operator-sdk run bundle docker.io/aliok/operator2-bundle:v0.0.1

```
