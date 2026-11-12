module github.com/oliverbestmann/webgpu

go 1.25

require (
	github.com/go-gl/glfw/v3.4/glfw v0.1.0-pre.1.0.20260628091122-0bd588dc30cf
	github.com/oliverbestmann/webgpu/libs-android v0.0.0-20260906104815-e05521a4bbbb
	github.com/oliverbestmann/webgpu/libs-darwin v0.0.0-20260906104759-17c42a852178
	github.com/oliverbestmann/webgpu/libs-ios v0.0.0-20261112094144-2425638651e5
	github.com/oliverbestmann/webgpu/libs-linux v0.0.0-20260906104812-41e8774c38c4
	github.com/oliverbestmann/webgpu/libs-windows v0.0.0-20260628152801-f47d1b682eb8
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract v1.27.0 // published before deciding on a version scheme. we start at v1.0.0
