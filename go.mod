module github.com/oliverbestmann/webgpu

go 1.25

require (
	github.com/go-gl/glfw/v3.4/glfw v0.1.0-pre.1.0.20260628091122-0bd588dc30cf
	github.com/oliverbestmann/webgpu/libs-android v0.0.0-20261112101134-9fc0777bdf01
	github.com/oliverbestmann/webgpu/libs-darwin v0.0.0-20261112101121-f12e930cb707
	github.com/oliverbestmann/webgpu/libs-ios v0.0.0-20261112101124-0a18a0d4d894
	github.com/oliverbestmann/webgpu/libs-linux v0.0.0-20261112101129-823eab7d2391
	github.com/oliverbestmann/webgpu/libs-windows v0.0.0-20261112110423-e4fd7225cf83
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract v1.27.0 // published before deciding on a version scheme. we start at v1.0.0
