# Release Notes

Generates markdown-style release notes from the pull requests merged between
two given branches or tags

### Usage

```sh
go run main.go \
  --username gerad \
  --password my-github-password \
  --repo release-notes \
  --base v1.0 \
  --head v2.0
```
