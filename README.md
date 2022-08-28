# lrview

[![Deploy Status][deploy-workflow-badge]][deploy-workflow-url]
[![Docker Status][docker-workflow-badge]][docker-workflow-url]

lrview let you visualize a Lightroom catalog in the browser by exploiting the previews computed by Lightroom.  
Its main use case is to display a Lightroom catalog stored on a NAS without a Lightroom installation.

A demonstration instance is live at [lrview.fly.dev](https://lrview.fly.dev).

## Usage

Run lrview using one of the method below, then open your browser at [localhost:8080](http://localhost:8080).

### Docker

```bash
docker run \
  -e LRVIEW_CATALOG_PATH=/MyCatalog.lrcat \
  -p 8080:8080 \
  -v "$(pwd)/MyCatalog.lrcat":"/MyCatalog.lrcat" \
  -v "$(pwd)/MyCatalog Previews.lrdata":"/MyCatalog Previews.lrdata" \
  ghcr.io/maxmouchet/lrview:main
```

### Nix

```bash
export LRVIEW_CATALOG_PATH=MyCatalog.lrcat
nix run github:maxmouchet/lrview
```

### Source

```bash
git clone git@github.com:maxmouchet/lrview.git && cd lrview
export LRVIEW_CATALOG_PATH=MyCatalog.lrcat
go run main.go
```

[deploy-workflow-badge]: https://img.shields.io/github/workflow/status/maxmouchet/lrview/Deploy?logo=github&label=deploy

[deploy-workflow-url]: https://github.com/maxmouchet/lrview/actions/workflows/deploy.yml

[docker-workflow-badge]: https://img.shields.io/github/workflow/status/maxmouchet/lrview/Docker?logo=github&label=docker

[docker-workflow-url]: https://github.com/maxmouchet/lrview/actions/workflows/docker.yml

