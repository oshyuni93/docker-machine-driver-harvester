package harvester

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckNetworkData(t *testing.T) {
	type args struct {
		networkDataStr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				networkDataStr: "",
			},
			wantErr: false,
		},
		{
			name: "without network section key",
			args: args{
				networkDataStr: `
version: 1
config:
- type: physical
  name: enp1s0
  subnets:
  - type: dhcp
`,
			},
			wantErr: false,
		},
		{
			name: "1 dhcp",
			args: args{
				networkDataStr: `
network:
  version: 1
  config:
   - type: physical
     name: enp1s0
     subnets:
     - type: dhcp
`,
			},
			wantErr: false,
		},
		{
			name: "1 static",
			args: args{
				networkDataStr: `
network:
  version: 1
  config:
   - type: physical
     name: enp1s0
     subnets:
     - type: static
       address: 192.168.5.91/24
       gateway: 192.168.5.1
   - type: nameserver
     interface: enp1s0
     address:
        - 192.168.5.1
`,
			},
			wantErr: false,
		},
		{
			name: "1 static without gateway",
			args: args{
				networkDataStr: `
network:
  version: 1
  config:
   - type: physical
     name: enp1s0
     subnets:
     - type: static
       address: 192.168.5.91/24
   - type: nameserver
     interface: enp1s0
     address:
        - 192.168.5.1
`,
			},
			wantErr: true,
		},
		{
			name: "1 static without nameserver",
			args: args{
				networkDataStr: `
network:
  version: 1
  config:
   - type: physical
     name: enp1s0
     subnets:
     - type: static
       address: 192.168.5.91/24
       gateway: 192.168.5.1
`,
			},
			wantErr: true,
		},
		{
			name: "2 dhcp",
			args: args{
				networkDataStr: `
network:
  version: 1
  config:
   - type: physical
     name: enp1s0
     subnets:
     - type: dhcp
   - type: physical
     name: enp2s0
     subnets:
     - type: dhcp
`,
			},
			wantErr: false,
		},
		{
			name: "1 dhcp and 1 static with gateway",
			args: args{
				networkDataStr: `
network:
  version: 1
  config:
   - type: physical
     name: enp1s0
     subnets:
     - type: dhcp
   - type: physical
     name: enp2s0
     subnets:
     - type: static
       address: 192.168.5.91/24
       gateway: 192.168.5.1
`,
			},
			wantErr: false,
		},
		{
			name: "1 dhcp and 1 static without gateway",
			args: args{
				networkDataStr: `
network:
  version: 1
  config:
   - type: physical
     name: enp1s0
     subnets:
     - type: dhcp
   - type: physical
     name: enp2s0
     subnets:
     - type: static
       address: 192.168.5.91/24
`,
			},
			wantErr: false,
		},
		{
			name: "2 static with gateway",
			args: args{
				networkDataStr: `
network:
  version: 1
  config:
   - type: physical
     name: enp1s0
     subnets:
     - type: static
       address: 192.168.5.91/24
       gateway: 192.168.5.1
   - type: physical
     name: enp2s0
     subnets:
     - type: static
       address: 192.168.5.92/24
       gateway: 192.168.5.1
`,
			},
			wantErr: true,
		},
		{
			name: "2 static without gateway",
			args: args{
				networkDataStr: `
network:
  version: 1
  config:
   - type: physical
     name: enp1s0
     subnets:
     - type: static
       address: 192.168.5.91/24
   - type: physical
     name: enp2s0
     subnets:
     - type: static
       address: 192.168.5.91/24
   - type: nameserver
     interface: enp1s0
     address:
        - 192.168.5.1
`,
			},
			wantErr: true,
		},
		{
			name: "1 static with gateway and 1 static without gateway",
			args: args{
				networkDataStr: `
network:
  version: 1
  config:
   - type: physical
     name: enp1s0
     subnets:
     - type: static
       address: 192.168.5.91/24
       gateway: 192.168.5.1
   - type: physical
     name: enp2s0
     subnets:
     - type: static
       address: 192.168.5.92/24
   - type: nameserver
     interface: enp1s0
     address:
        - 192.168.5.1
`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkNetworkData(tt.args.networkDataStr); (err != nil) != tt.wantErr {
				t.Errorf("CheckNetworkData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_parseVGPUInfo(t *testing.T) {
	vObj := &VGPUInfo{
		VGPURequests: []VGPURequest{
			{
				DeviceName: "nvidia.com/NVIDIA_A2-2Q",
			},
			{
				DeviceName: "nvidia.com/NVIDIA_A2-1Q",
			},
		},
	}
	vgpuInfoString := `{"vGPURequests":[{"name":"","deviceName":"nvidia.com/NVIDIA_A2-2Q"},{"name":"","deviceName":"nvidia.com/NVIDIA_A2-1Q"}]}`
	assert := require.New(t)
	v, err := parseVGPUInfo(vgpuInfoString)
	assert.NoError(err)
	assert.Equal(v, vObj, "expected request to match predefined object")
}

func Test_parseGPUInfo(t *testing.T) {
	assert := require.New(t)

	// 1. Valid GPU info
	gpuInfoString := `{"enabled":true,"vendor":"nvidia","model":"A100","pciPassthrough":true,"pciDevice":{"name":"gpu01-0000e2000","address":"0000:e2:00.0","nodeName":"gpu01","resourceName":"nvidia.com/GA100_A100_SXM4_40GB"}}`
	g, err := parseGPUInfo(gpuInfoString)
	assert.NoError(err)
	assert.NotNil(g)
	assert.True(g.Enabled)
	assert.Equal("A100", g.Model)
	assert.NotNil(g.PCIDevice)
	assert.Equal("gpu01-0000e2000", g.PCIDevice.Name)
	assert.Equal("nvidia.com/GA100_A100_SXM4_40GB", g.PCIDevice.ResourceName)
	assert.Equal("gpu01", g.PCIDevice.NodeName)

	// 2. Defensive checks: empty, undefined, null, invalid, disabled
	for _, invalidStr := range []string{"", "   ", "undefined", "null", "{}", "invalid json", `{"enabled":false}`} {
		res, err := parseGPUInfo(invalidStr)
		assert.NoError(err, "should not return error for: "+invalidStr)
		assert.Nil(res, "expected nil GPUInfo for: "+invalidStr)
	}
}
