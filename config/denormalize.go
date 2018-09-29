package config

import "fmt"

// TODO refactor!
// TODO error check, too?
// TODO simplify provider lookups once this is verified
func Denormalize(in Config) (*Config, error) {
	out := Config{
		RawConfig:   in.RawConfig,
		Global:      in.Global,
		Server:      in.Server,
		Scheduler:   in.Scheduler,
		Downloaders: in.Downloaders,
	}

	for _, datastore := range in.Analyses.Datastores {
		out.Analyses.Datastores = append(out.Analyses.Datastores, datastore)
	}

	for _, pivnetfile := range in.PivnetFiles.Datastores {
		a := pivnetfile.AnalysisDatastore

		if a != nil {
			if a.Type != "" {
				a.Name = fmt.Sprintf("internal/pivnetfile/%s", pivnetfile.Name)

				out.Analyses.Datastores = append(
					out.Analyses.Datastores,
					AnalysisDatastoreConfig{
						AbstractComponentConfig: AbstractComponentConfig{
							Name:    a.Name,
							Type:    a.Type,
							Options: a.Options,
						},
					},
				)

				a.Type = ""
				a.Options = nil
			}
		} else {
			pivnetfile.AnalysisDatastore = &AnalysisDatastoreConfig{
				AbstractComponentConfig: AbstractComponentConfig{
					Name: "default",
				},
			}
		}

		out.PivnetFiles.Datastores = append(out.PivnetFiles.Datastores, pivnetfile)
	}

	for _, stemcell := range in.Stemcells.Datastores {
		a := stemcell.AnalysisDatastore

		if a != nil {
			if a.Type != "" {
				a.Name = fmt.Sprintf("internal/stemcell/%s", stemcell.Name)

				out.Analyses.Datastores = append(
					out.Analyses.Datastores,
					AnalysisDatastoreConfig{
						AbstractComponentConfig: AbstractComponentConfig{
							Name:    a.Name,
							Type:    a.Type,
							Options: a.Options,
						},
					},
				)

				a.Type = ""
				a.Options = nil
			}
		} else {
			stemcell.AnalysisDatastore = &AnalysisDatastoreConfig{
				AbstractComponentConfig: AbstractComponentConfig{
					Name: "default",
				},
			}
		}

		out.Stemcells.Datastores = append(out.Stemcells.Datastores, stemcell)
	}

	for _, releaseCompilation := range in.ReleaseCompilations.Datastores {
		a := releaseCompilation.AnalysisDatastore

		if a != nil {
			if a.Type != "" {
				a.Name = fmt.Sprintf("internal/release-compilation/%s", releaseCompilation.Name)

				out.Analyses.Datastores = append(
					out.Analyses.Datastores,
					AnalysisDatastoreConfig{
						AbstractComponentConfig: AbstractComponentConfig{
							Name:    a.Name,
							Type:    a.Type,
							Options: a.Options,
						},
					},
				)

				a.Type = ""
				a.Options = nil
			}
		} else {
			releaseCompilation.AnalysisDatastore = &AnalysisDatastoreConfig{
				AbstractComponentConfig: AbstractComponentConfig{
					Name: "default",
				},
			}
		}

		out.ReleaseCompilations.Datastores = append(out.ReleaseCompilations.Datastores, releaseCompilation)
	}

	for _, release := range in.Releases.Datastores {
		{
			a := release.AnalysisDatastore

			if a != nil {
				if a.Type != "" {
					a.Name = fmt.Sprintf("internal/release/%s", release.Name)

					out.Analyses.Datastores = append(
						out.Analyses.Datastores,
						AnalysisDatastoreConfig{
							AbstractComponentConfig: AbstractComponentConfig{
								Name:    a.Name,
								Type:    a.Type,
								Options: a.Options,
							},
						},
					)

					a.Type = ""
					a.Options = nil
				}
			} else {
				release.AnalysisDatastore = &AnalysisDatastoreConfig{
					AbstractComponentConfig: AbstractComponentConfig{
						Name: "default",
					},
				}
			}
		}

		{
			c := release.CompilationDatastore

			if c != nil {
				if c.Type != "" {
					c.Name = fmt.Sprintf("internal/release/%s", release.Name)

					out.ReleaseCompilations.Datastores = append(
						out.ReleaseCompilations.Datastores,
						ReleaseCompilationDatastoreConfig{
							AbstractComponentConfig: AbstractComponentConfig{
								Name:    c.Name,
								Type:    c.Type,
								Options: c.Options,
							},
						},
					)

					c.Type = ""
					c.Options = nil
				}
			} else {
				release.CompilationDatastore = &ReleaseCompilationDatastoreConfig{
					AbstractComponentConfig: AbstractComponentConfig{
						Name: "default",
					},
				}
			}

			{
				a := release.CompilationDatastore.AnalysisDatastore

				if a != nil {
					if a.Type != "" {
						a.Name = fmt.Sprintf("internal/release.compilation/%s", release.Name)

						out.Analyses.Datastores = append(
							out.Analyses.Datastores,
							AnalysisDatastoreConfig{
								AbstractComponentConfig: AbstractComponentConfig{
									Name:    a.Name,
									Type:    a.Type,
									Options: a.Options,
								},
							},
						)

						a.Type = ""
						a.Options = nil
					}
				} else {
					release.CompilationDatastore.AnalysisDatastore = &AnalysisDatastoreConfig{
						AbstractComponentConfig: AbstractComponentConfig{
							Name: "default",
						},
					}
				}
			}
		}

		out.Releases.Datastores = append(out.Releases.Datastores, release)
	}

	return &out, nil
}
