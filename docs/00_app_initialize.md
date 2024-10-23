
```mermaid
sequenceDiagram
    autonumber
    participant App as UCW Demo App
    participant SDK as UCW SDK

    App ->> SDK: Init
    activate App
    SDK ->> App: return success
    deactivate App
```
