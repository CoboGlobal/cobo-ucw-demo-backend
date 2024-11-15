```mermaid
sequenceDiagram
    autonumber
    participant App as Client App
    participant SDK as UCW SDK

    App ->> SDK: Init
    activate App
    SDK ->> App: return success
    deactivate App
```