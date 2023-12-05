# Contributing

Thank you for being interested in contributing to the project, all contributions are welcome!
See the [kube-logging organization's CONTRIBUTING.md](https://github.com/kube-logging/.github/blob/main/CONTRIBUTING.md) file for general contribution guidelines.

Please read the [Organisation's Code of Conduct](https://github.com/kube-logging/.github/blob/main/CODE_OF_CONDUCT.md)!

## Ways of contributing

There are multiple ways you can help with the project.

### Questions

The following channels are available for logging operator related discussions:

- [Github discussions](https://github.com/orgs/kube-logging/discussions)
- [#logging-operator Discord channel](https://discord.gg/eAcqmAVU2u)

### Reporting an issue

Please use the appropriate template for your issue type (if applicable):

- [Bug report](https://github.com/kube-logging/logging-operator/issues/new?assignees=&labels=bug&projects=&template=---bug-report.md&title=)
- [Feature request](https://github.com/kube-logging/logging-operator/issues/new?assignees=&labels=&projects=&template=--feature-request.md&title=)
- [Don’t see your issue here? Open a blank issue.](https://github.com/kube-logging/logging-operator/issues/new)

Security related issue (security vulnerability)
The Kube Logging team and community take all security issues seriously. Thank you for improving the security of our projects. We appreciate your efforts and responsible disclosure and will make every effort to acknowledge your contributions.
Please follow the [Security Policy](https://github.com/kube-logging/logging-operator/security/policy) for the next steps.

### Opening a pull request (PR)

Steps:

1. Check if there is an [open issue](https://github.com/kube-logging/logging-operator/issues) that you can reference in your pull request
2. Fork the [repository](https://github.com/kube-logging/logging-operator/fork)
3. Create a new branch using `git switch -c <branch-name>`
4. Make your modifications and add unit tests to cover most code branches, input fields, and functionality if possible[^1]
5. Run `make fmt` to run code formatting
6. Run `make generate` to generate code, documentation, manifests, etc.
7. Run `make check` to run license checks, linter, tests
8. (Optional) run `make lint-fix` to fix trivial linting problems
9. Commit your modifications (preferably many smaller commits) using `git commit -s`, because [DCO check](https://github.com/apps/dco) is required
10. Push your commits (`git push <fork_origin> <fork_name>`) and open a pull request

See the [For developers](https://kube-logging.dev/docs/developers/) page in the documentation for details (e.g. setting up a development environment, adding documentation).

[^1]: For bigger features please open a new [Feature request issue](https://github.com/kube-logging/logging-operator/issues/new?assignees=&labels=&projects=&template=--feature-request.md&title=) where we can have a discussion about the proposed feature and the suggested test cases.

### Add yourself to the production adopters list

If you use the Logging operator in a production environment, add yourself to the list of production [adopters](https://github.com/kube-logging/logging-operator/blob/master/ADOPTERS.md).🤘
