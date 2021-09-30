# mirror-docker-tags-action
This action insert a matrix variable that contains all the tags that needs to be updated.

## Inputs

## `from`
**Required** The repository and tags to mirror, example `"ubuntu[20.10,latest],alpine[latest]"`.

## `to`
**Required** The repository owner to create the new tags (GITHUB_REPOSITORY_OWNER by default). 

## `extraRegistry`
**Optional** By default only created tags to docker hub, you can add an extra registry, example `"ghcr.io"`

## `updateAll`
**Optional** Boolean that indicates if all the tags needs to be updated, if false only updates tags that are new or has recently being updated.

## Outputs

## `matrix`

A matrix that can be used for other jobs
```
{ include: [
    base_img : ubuntu,
    tags:  repo/ubuntu:latest,repo/ubuntu:20.10
    platforms:  Linux/amd64,Linux/arm/v7
]}
```

## Example usage

```
jobs:
  setup-matrix:
    runs-on: ubuntu-latest
    steps:
      - name: Get missing tags
        id: tags
        uses: s6on/mirror-docker-tags-action@v1.0.0
        with:
          from: alpine[latest,3],debian[latest,11,11-slim],ubuntu[latest,20.04,18.04]
          to: ${{ github.repository_owner }}
          extra-registry: ghcr.io
          updateAll: ${{ github.event_name != 'schedule' }}
    outputs:
      matrix: ${{ steps.tags.outputs.matrix }}
  build-push:
    runs-on: ubuntu-latest
    needs: setup-matrix
    if: ${{ needs.setup-matrix.outputs.matrix }}
    strategy:
      fail-fast: false
      matrix: ${{ fromJson(needs.setup-matrix.outputs.matrix) }}
    name: ${{ matrix.base_img }}
```