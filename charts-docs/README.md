# Helm chart documentation

README files for Helm charts are generated using [helm-docs](https://github.com/norwoodj/helm-docs).

Each chart should contain a `README.md.gotmpl` file that describes how
the `README.md` of the chart should be generated.

Normally, this file can be the same as the primary template in [docs/templates/README.md.gotmpl] or a symlink pointing to it:

```bash
cd charts/CHART
ln -s ../../charts-docs/templates/README.md.gotmpl
```

Copy the file to the chart directory if you want to customize the template.

**Note:** Don't forget to add `README.md.gotmpl` to `.helmignore`.

Then run `make docs` in the repository root.
