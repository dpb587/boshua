# GraphQL API

The `/api/v2/graphql` endpoint provides a GraphQL API with query and mutation support.


## Query: `release`

Get information about a release...

    {
      release(name: String, version: String, url: String, checksum: String, labels: [String]) {
        name
        version
        labels

        tarball { #artifact }

        compilation(os: String, version: String) {
          tarball { #artifact }
          analysis { #analysis }
        }

        analysis {
          results(analyzers: [String]) {
            analyzer
            artifact { #artifact }
          }
        }
      }
    }


## Query: `stemcell`

Get information about a stemcell...

    {
      stemcell(iaas: String, hypervisor: String, os: String, version: String, diskFormat: String, flavor: String, labels: [String]) {
        iaas
        hypervisor
        os
        flavor
        diskFormat
        version

        tarball { #artifact }
        analysis { #analysis }
      }
    }


## Mutation: `scheduleCompilation`

Schedule compilation for a release...

    mutation {
      scheduleReleaseCompilation(name: String, version: String, url: String, checksum: String, osName: String, osVersion: String) {
        status
      }
    }


## Mutation: `scheduleStemcellAnalysis`

Schedule analysis of a stemcell...

    mutation {
      scheduleStemcellAnalysis( #stemcellParams, analyzer: String) {
        status
      }
    }


## Mutation: `scheduleReleaseAnalysis`

Schedule analysis of a stemcell...

    mutation {
      scheduleReleaseAnalysis( #releaseParams, analyzer: String) {
        status
      }
    }


## Mutation: `scheduleReleaseCompilationAnalysis`

Schedule analysis of a stemcell...

    mutation {
      scheduleReleaseAnalysis( #releaseParams, #releaseCompilationParams, analyzer: String) {
        status
      }
    }

