


1. 登陆

```javascript
    LoginRes = client.Login(email)
    client.CreateVault()
    if (LoginRes.MainKeyGroupID != "") {
        if (databaseFileExist) { // db initialized
            password = local.readpassword()
            success = sdk.open(password)
            if (!success) {
                navigate(resetView)
            } else {
                nodeID = sdk.GetNodeID()
                group = client.GetGroup(nodeID)
                
                if (group.Type == "MainKey") {
                    navigate(walletView)
                } else if (group.Type == "RecoveryKey") {
                    navigate(walletView)
                } else {
                    navigate(resetView)
                }
            }
        } else {
            if (pincodeSetted) {
                navigate(pinCodeView)
            } else {
                navigate(resetView)
            }
        }
    } else {
        if (!pincodeSetted) {
            navigate(pinCodeView)
        } else {
            password = local.readpassword()
            sdk.open(password)
            navigate(keyGenView)
        }
    }
```

2. PinCode

```javascript
// pincode setup success
    local.savePinCode()
    if (keygensuccess) {
        navigate(keyGenView)
    }
    password = local.randomPasswordAndSave()
    nodeID = sdk.intialize(password)
    client.CreateVault()
    client.CreateGroup()
    sdk.Open(password)
    
```