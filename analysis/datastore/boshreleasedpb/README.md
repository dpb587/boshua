Assuming a directory structure like...

  ./compiled-release
    ./{channel}
      ./{os_name}
        ./{os_version}
          ./analysis
            ./{analyzer}
              ./v{version}.meta4
          ./v{version}.meta4
  ./release
    ./{channel}
      ./analysis
        ./{analyzer}
          ./v{version}.meta4
      ./v{version}.meta4
