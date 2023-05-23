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
