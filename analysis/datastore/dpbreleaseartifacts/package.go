// Package dpbreleaseartifacts provides the "dpbreleaseartifacts" provider which
// supports querying and storing analysis results in a git repository.
//
// It supports releases, release compilations, and stemcell analysis results.
//
// The following directory structure is used within the repository:
//
//     ./{release_compilation_path}
//       ./{os_name}
//         ./{os_version}
//           ./analysis
//             ./{analyzer}
//               ./v{version}.meta4
//     ./{release_path}
//       ./analysis
//         ./{analyzer}
//           ./v{version}.meta4
package dpbreleaseartifacts
