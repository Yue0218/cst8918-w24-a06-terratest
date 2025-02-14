package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "a9fc334c-4082-43d8-95c1-4bb0f83c8a71"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "gao00094",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variable
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))

	// Confirm NIC exists and is connected to the VM
	nics := azure.GetVirtualMachineNics(t, vmName, resourceGroupName, subscriptionID)
	nicName := terraform.Output(t, terraformOptions, "nic_name")
	assert.Contains(t, nics, nicName, "NIC should be attached to the VM")

	// Confirm VM is running the correct Ubuntu version
	expectedUbuntuVersion := "24.04"
	vmUbuntuVersion := azure.GetVirtualMachineImage(t, vmName, resourceGroupName, subscriptionID)
	assert.Contains(t, vmUbuntuVersion, expectedUbuntuVersion, "VM should be running the correct Ubuntu version")
}
