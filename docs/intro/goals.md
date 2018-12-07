# Goals

 * provide consistent ways to generate and store metadata about releases, compilations, and stemcells
 * support private and public configurations with access to different sets of resources
 * deprecate one-off, duplicated scripts which have historically implemented these processes
 * allow compiled releases to be dynamically discovered and avoid hard-coding references
 * support just-in-time release compilations, shared across environments
 * focus on CLI and API extensibility to enable this as a building block
 * support both local and remote execution of any commands

