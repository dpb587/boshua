# TODO

A place to keep track of intentions and goals...

## General

## Analysis

 * namespace built-in analyzers into `boshua.*`; allows custom, external ones to be registered, eventually
 * refactor/convert artifactfiles.v1 to refactored, recursive artifact analyzer (like tiles)


## Scheduler

 * default ops file/image for container-based compilation instead of external director


## Technical

 * grep for `TODO` for some specific code concerns


## Random

 * `... datastore filter` ~> `versions`
 * `analysis generate --analyzer=x` ?> `analysis generate x`
 * configurable analyzers
 * api download mirror/proxy
 * nicer wrapper libraries - `r, _ := boshua.Release("openvpn/5.1.0"); p, _ := r.Packages(); return p[0].Name`
 * testing
 * webui
 * `deployment download-artifacts-cache-something` for filling ~/.bosh/cache with compiled versions?
