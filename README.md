# file-sync

A GitHub Action to sync files from another repository

![test](https://github.com/champ-oss/file-sync/workflows/gotest/badge.svg)

## Example

```yaml
jobs:
  run:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          token: ${{ secrets.GITHUB_TOKEN }}

      - uses: champ-oss/file-sync
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
          repo: champ-oss/terraform-module-template
          files: |
            .gitignore
            LICENSE
```