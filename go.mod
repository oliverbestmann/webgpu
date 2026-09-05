module github.com/dhannyell/webgpu

go 1.25

require (
	github.com/dhannyell/webgpu/libs-android v0.0.0-20260905020343-c4069cac051c
	github.com/dhannyell/webgpu/libs-darwin v0.0.0-20260905020352-092624e3a385
	github.com/dhannyell/webgpu/libs-ios v0.0.0-20260905020358-fb13dc722108
	github.com/dhannyell/webgpu/libs-linux v0.0.0-20260905020403-4ea5c475cbbd
	github.com/dhannyell/webgpu/libs-windows v0.0.0-20260905020407-128027db48f4
	github.com/go-gl/glfw/v3.4/glfw v0.1.0-pre.1.0.20260628091122-0bd588dc30cf
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract v1.27.0 // published before deciding on a version scheme. we start at v1.0.0
