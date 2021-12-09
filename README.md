```mermaid
flowchart LR
  subgraph accounts
    direction LR
    accounts_server[server] --> accounts_db[(postgres)]
  end
  request[\request\] --> gateway
  gateway -->|/accounts| accounts
```