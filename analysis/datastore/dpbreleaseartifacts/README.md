Assuming a directory structure like...

  ./{compiled_release_prefix}
    ./{os_name}
      ./{os_version}
        ./analysis
          ./{analyzer}
            ./v{version}.meta4
        ./v{version}.meta4
  ./{release_prefix}
    ./analysis
      ./{analyzer}
        ./v{version}.meta4
    ./v{version}.meta4
