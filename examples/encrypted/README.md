# Encrypted Habitat example

By default supervisors will communicate with no encryption. This example demostrates how to secure the communication.

## Workflow

The user needs to generate a [ring
key](https://www.habitat.sh/docs/run-packages-security/) using `hab ring key generate foobar`, and then base64
encode it (on Linux) with `hab ring export foobar | base64 -w 0` (please refer to
[this
document](https://kubernetes.io/docs/concepts/configuration/secret/#creating-a-secret-manually)
for platform-specific instructions on base64 encoding).

The encoded key can then be used as the value of the `ring-key` key in a Kubernetes
secret.

The secret's name must be the same as the filename of the key, minus the
extension.

For example, for a key named `foobar`, the key file might be something like
`foobar-20170824094632.sym.key`, and the secret name must be
`foobar-20170824094632`.

The secret's name must additionally be referenced in the `Habitat` object's `ringSecretName` key.

After the Habitat operator is up and running, execute the following command from the root of this repository:

```
kubectl create -f examples/encrypted/habitat.yml
```

## Deletion

The Habitat operator does not delete the Secret on Habitat deletion. This is
because the user might want to re-use the secret across multiple
`Habitat`s and `Habitat` lifecycles.

To delete the Secret simply run:

```
  kubectl delete secret example-encrypted-ring-20170829113029
```
