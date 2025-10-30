# List of pinned Go modules

The following modules are not upgradable without refactoring application code and need to be version-frozen:

| Module | Current | Latest |
|---|---|---|
| github.com/cilium/ebpf | v0.12.3 | v0.19.0 |
| github.com/gizak/termui | v2.3.1-0.20180817033724-8d4faad06196+incompatible | v3.1.0+incompatible |
| github.com/opencontainers/runc | v1.1.14 | v1.3.0 |

Notes:
- Updating these modules will require code changes; pin their versions in go.mod or vendor as appropriate.
- Keep this file updated when the application is refactored to support newer versions.
