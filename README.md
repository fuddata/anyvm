# AnyVM
## Usage
### List VMs
```powershell
# From all providers
$VMs = (Invoke-RestMethod http://192.168.8.40:8080/api/v1/vms).data

# From one provider
(Invoke-RestMethod http://192.168.8.40:8080/api/v1/vms?provider=azure).data
```

### Create VM
```powershell
$payload = @{
    provider          = "" # azure , aws or gcp
    vmName            = "mynewtestvm"
    resourceGroupName = "script-test"
    location          = "westeurope"
    vmSize            = "small"
    adminUsername     = "azureuser"
    adminPassword     = "P@ssw0rd1234"
    nicId             = "/subscriptions/54e30869-75a2-47ed-8b32-1057e61707f0/resourceGroups/script-test/providers/Microsoft.Network/networkInterfaces/myNIC"
}
$jsonPayload = $payload | ConvertTo-Json -Depth 5
$apiUrl = "http://192.168.8.40:8080/api/v1/vms/create"
Invoke-RestMethod -Method Post -Uri $apiUrl -Body $jsonPayload -ContentType "application/json"
```
