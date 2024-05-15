<a name="unreleased"></a>
## [Unreleased]


<a name="3.0.1"></a>
## [3.0.1] - 2023-09-05

<a name="3.0.0"></a>
## [3.0.0] - 2023-02-13
### Added
- FuncMap To Template
- Validation For User Input
- Validation To Placeholders
- Template Schema Validation

### Changed
- Validation Rule Names To Camel Case
- Throw An Error When Tmpl Path And Out Path Match
- More Debug Messages With Config
- Load Config Setting From File
- Make TmplToApp Cache Hidden
- Refactor To Internal
- Upgraded Kohirens Stdlib To Version 0.4.0
- Updated the templpate.json Schema
- Begin Refactoring Template Internals

### Fixed
- Rule regExp Expression Dependency
- Bug Generating A Template When There Is No Config
- Mapping Functions To Templates
- An Error Message
- Directory Mode Permissions Mask


<a name="2.0.2"></a>
## [2.0.2] - 2023-02-07
### Changed
- Upgrade Docker iamge to Alpine 3.17

### Fixed
- Typos In tmplt-type Messages


<a name="2.0.1"></a>
## [2.0.1] - 2022-10-10
### Fixed
- Building Docker Image


<a name="2.0.0"></a>
## [2.0.0] - 2022-10-10
### Added
- Command To Generate A Manifest
- Append git ref name (tag) to cache folder.
- Version to verbosity output.
- Message to inicate cahcue use in verbose logging.
- Error message.
- skip to template schema to skip files or directories.
- Copy files excluded from parsing as-is.
- config sub-command.
- Flag default-val to aid automation.
- Doc for building a template.json.
- JSON schemas for the template.json and answers JSON files.
- logf function for rendering message to the user.

### Changed
- Migrate Messaging to Various Arrays
- Append Version to Cached Download
- Exclude template.json from out-path.
- Update config usage.
- Skip .git directory from copying or parsing.
- Config stderr messages.
- Pass custom Usage method the configration.
- Migrate various log messages to errors and messages structures.
- Show errir for incorrect CLI order.
- Flag names.
- Make answers.json file optional.
- Template schema JSON.
- From -a to -i for answer short flag.
- Add more feedback when user needs to add values for missing placeholders.
- Move code to generate semver info to main.go.
- Move load anwwer errors to the error file.

### Fixed
- Merging user settings into the configuration.
- Printing usage info.
- config sub-command returned 1 when only help flag passed in.

### Removed
- Informational logging preventing template output.
- Config from verbosity output.
- Use of allowed URLs.


<a name="1.0.0"></a>
## [1.0.0] - 2022-06-26
### Added
- Centralized store for program messages.

### Changed
- Separate logic to determine if the template is local or remote.
- Separate the template type logic from template location.
- Set branch default value to main.
- Make answersPath option input.
- Rename variable.
- Verbosity code.
- Refactored error messags.

### Fixed
- Only download a zip when remote.
- But with setting positional clig flags.

### Removed
- Check URL is allowed.
- Unused variable.


<a name="0.3.1"></a>
## [0.3.1] - 2022-05-24
### Fixed
- Publishing a Docker image of the release.


<a name="0.3.0"></a>
## [0.3.0] - 2022-05-16
### Added
- Version info.

### Changed
- Upgraded CI kohirens/circleci-go image version 0.3.0.


<a name="0.2.1"></a>
## [0.2.1] - 2022-05-16
### Fixed
- Builing tagged Docker image after a release publish.


<a name="0.2.0"></a>
## [0.2.0] - 2022-05-16
### Added
- Auto Build Executables to Pipeline

### Changed
- Updated README with how to download with cURL.
- Return path from gitClone.

### Fixed
- Docker image build.


<a name="0.1.1"></a>
## [0.1.1] - 2022-05-14

<a name="0.1.0"></a>
## 0.1.0 - 2022-02-20
### Added
- Excludes field to template.json Schema.
- File for collective errors.
- Method to read in answers from template.json questions.
- Function to read the template.json file required in a template.
- Detection of 7zip in the environment.
- TODO.md
- Parsing local directories as templates.
- Error checking for appPath when extracting flag value.
- Parse local directory.
- Documentation on loadAnswers.
- loadAnswers method.
- Methdo to detect template path type.
- Methdo to detect template path type.
- Flow Diagrams.
- Local dev instructions.
- Method to detect text file via file extensions.
- Vendor dependency config.
- Template processing.
- PathExist function to stdlib.
- Method to extract a template package.
- Code coverage visuals to vscode config.
- Default config.
- Initialize a config.
- Validation of URLs domains allowed to download.
- stdlib package.
- Installation instructions to the readme.
- Algorithm.
- Debugging with Delve and VS Code Remote.
- Lead config from file.
- WIP, Download template feature.
- More tools to the dev environment.
- Errors when missing required arguments.
- Parameter parsing.
- Docker configuration.

### Changed
- Renamed from Bootup to TmplToApp.
- Where extract places files.
- Error message array name to errs.
- Download zips to a unique name in the cache.
- Upgrade to Kohirens STDLib to version 0.1.1, go.sum updated.
- Upgraded Kohirens STDLIB version to 0.1.2.
- Dev env container name.
- Read values from template.json file.
- No longer throw an error when skipping files to run through the template engine.
- Update deve environment.
- Rebrand
- Refactoring Docker
- Allow user to set the file extensions to filter in the app config.
- Set app config property tmpl path when passing in CLI local path.
- Pull in default app config using unmarshal to struct type.
- Use one app config per run.
- Removed verbosity level from app config and add tmpl.
- Move directory where go-gitter stores files.
- Updated documentation.
- Refactor flags code.
- Error handling on config.settings function.
- Refactor function to initialize a config.
- Replace hard coded value with a constant.
- Updated local development docs.
- Editor config to be compatible with GoLang code fmt.
- Made function private that are not part of the public API.
- Refactoring for better organization.
- Refactor converge template package into main.
- Upgraded to go 1.16.
- Dev default environment variables.
- Extract the stdlib out to its own repo.
- Updated the development environment.
- Turn on Go tools extension checking to stay up to date.
- Map SSH keys in dev environment.
- Rename download to template.
- Pre-install vscode-server.
- Constant TEST_TMP to be the same in all test files.
- Rename getArgs to parseArgs.
- To use a flagset in place of global flag space.
- To Go module.
- Verbosity flag name.
- editorcofig line endings to Linux.
- Replaced rsyslog with tini and tail.
- Settle on a container running in dev.

### Fixed
- Unit test.
- Erros with reading template.json file.
- VS Code Remote container setup.
- Downloaded zips and ignore empty files.
- Malformed JSON in the default config.
- Test output going to wrong directory.
- Setting exit code from main when no errors.
- Get template.parseDir working.
- Return correct error from template.parseDir.
- Test template repository name for test arguement.
- Unit test in stdlib.
- Test in stdlib.
- Broken test for downloading templates.
- Download to save archive based on URL name.
- Testing suite.
- logic bugs.
- Missed file required to run settings test.
- Settings in VS Code devcontainer config.
- Added missing constants.
- passing flags to main and running tests.

### Removed
- Unused CopyDir function.

### BREAKING CHANGE

* Renamed `tplPath` flag to `tmplPath`.
* Changed `-v` no longer set verbosity level, but is a short for
  version.


[Unreleased]: https://github.com/kohirens/tmpltoapp.git/compare/3.0.1...HEAD
[3.0.1]: https://github.com/kohirens/tmpltoapp.git/compare/3.0.0...3.0.1
[3.0.0]: https://github.com/kohirens/tmpltoapp.git/compare/2.0.2...3.0.0
[2.0.2]: https://github.com/kohirens/tmpltoapp.git/compare/2.0.1...2.0.2
[2.0.1]: https://github.com/kohirens/tmpltoapp.git/compare/2.0.0...2.0.1
[2.0.0]: https://github.com/kohirens/tmpltoapp.git/compare/1.0.0...2.0.0
[1.0.0]: https://github.com/kohirens/tmpltoapp.git/compare/0.3.1...1.0.0
[0.3.1]: https://github.com/kohirens/tmpltoapp.git/compare/0.3.0...0.3.1
[0.3.0]: https://github.com/kohirens/tmpltoapp.git/compare/0.2.1...0.3.0
[0.2.1]: https://github.com/kohirens/tmpltoapp.git/compare/0.2.0...0.2.1
[0.2.0]: https://github.com/kohirens/tmpltoapp.git/compare/0.1.1...0.2.0
[0.1.1]: https://github.com/kohirens/tmpltoapp.git/compare/0.1.0...0.1.1
