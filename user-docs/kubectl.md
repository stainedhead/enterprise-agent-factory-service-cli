# kubectl access to a team's namespace

```bash
afs teams get-kubeconfig --team-id team_01J...
kubectl get pods
```

`get-kubeconfig` asks the API for the team's kubeconfig and merges it into your own. The API's kubeconfig has a
context named `afs-<team-slug>`, a cluster and a user whose exec credential plugin runs `aws eks get-token`. It carries
**no token**; kubectl asks AWS for one each time, with the role your identity provider federates to.

## Where it is written

| You give | The file |
|---|---|
| nothing | the first path of `$KUBECONFIG`, else `~/.kube/config` |
| `--kubeconfig ./kc.yaml` | that file; `~/.kube/config` is not touched |
| `--kubeconfig -` | nothing is written; the API's YAML is printed unchanged to stdout |

Merging replaces only the cluster, context and user entries that have the same names as the incoming ones; every other
entry, setting and comment is kept, and `current-context` becomes the team's context. Running the command twice leaves one
context. The file is written atomically (a temporary file, then a rename), so a failure never leaves half a file.
A new file is created with mode 0600 (and `~/.kube` with 0700). An existing file keeps its mode, and a warning
says so when it is wider than 0600.

## What `afs` refuses to merge

Because kubectl runs the plugin command on every call, `afs` checks the service's kubeconfig before it touches your
file. It exits 255 and leaves the file unchanged when a user's exec plugin is not `aws` with arguments that start
`eks get-token`, when a user has no exec plugin (for example a static token), when the plugin sets an environment
variable that does not start with `AWS_`, or when a cluster, context, user or `current-context` name does not start
with `afs`. Entries you own with other names are therefore never replaced. `--kubeconfig -` prints the response as is
and merges nothing.

## Another AWS profile

```bash
afs teams get-kubeconfig --team-id team_01J... --aws-profile ops
```

adds `AWS_PROFILE=ops` to the plugin's `env` block, so kubectl asks AWS for the token with that profile.

## Requirements

kubectl needs the AWS CLI v2 on `PATH` to run `aws eks get-token`. `afs` warns when it cannot find `aws`; the merge still
happens. If kubectl access is switched off for your role, the API answers `AccessDenied` and `afs` exits 254.
