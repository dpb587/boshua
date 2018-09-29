# `pivnet://`

Configuration:

 * `apiToken` - API token (from [Profile](https://network.pivotal.io/users/dashboard/edit-profile) page)
 * TODO: proxy settings


## URL

Syntax:

    pivnet://{host}/api/v2/products/{productName}/releases/{releaseId}/product_files/{productFileId}/download?extract=something.tgz!else.tgz&range=XX-YY

Parameters:

 * `host` - hostname of Pivotal Network (optional; default `network.pivotal.io`)
 * `productName` - product name (e.g. `pivotal-mysql`)
 * `releaseId` - release ID of the product (e.g. `193224`)
 * `productFileId` - product file ID of the release (e.g. `221505`)
 * `extract` - path within the asset to extract (e.g. `releases/dedicated-mysql-0.72.1-ubuntu-xenial-97.17-20180919-215132-069364325.tgz`)
 * TODO: `range` - a specific byte range where the initial `extract` file can be found to optimize downloads (e.g. `10293914-19289192`)
