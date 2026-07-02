---
sidebar_position: 3
title: With Java
---

# Java SDK

The Kobo Java SDK allows enterprise integrations to leverage the reconciliation engine seamlessly. It has zero external dependencies and uses the native Java 11+ `HttpClient`.

## Installation (Maven)

```xml
<dependency>
    <groupId>tech.triumphsystems</groupId>
    <artifactId>kobo-sdk</artifactId>
    <version>1.0.0</version>
</dependency>
```

## Initialization

```java
import tech.triumphsystems.kobo.KoboClient;
import tech.triumphsystems.kobo.model.*;
import tech.triumphsystems.kobo.request.*;

public class Main {
    public static void main(String[] args) {
        // Production
        KoboClient kobo = KoboClient.of("kobo_live_pk_...", "kobo_live_sk_...");

        // Sandbox
        KoboClient sandbox = KoboClient.sandbox("kobo_test_pk_...", "kobo_test_sk_...");

        // Health Check
        HealthResponse health = kobo.health();
    }
}
```

---

## Identities

### Create Identity
```java
Identity identity = kobo.identities().create(
    CreateIdentityRequest.builder()
        .externalReference("user-12345")
        .displayName("John Doe")
        .build()
);
```

### Get Identity
```java
Identity identity = kobo.identities().get("idx_123abc");
```

### Update Identity
```java
Identity identity = kobo.identities().update("idx_123abc", 
    UpdateIdentityRequest.builder()
        .displayName("John D.")
        .build()
);
```

### Close Identity
```java
Identity identity = kobo.identities().close("idx_123abc", 
    CloseIdentityRequest.builder()
        .reason("User terminated")
        .build()
);
```

### Reopen Identity
```java
Identity identity = kobo.identities().reopen("idx_123abc");
```

---

## Accounts

### List Transactions
```java
// accountId, cursor, limit
TransactionPage page = kobo.accounts().listTransactions("acc_123abc", null, 50);
```

### Get Statement
```java
// accountId, period (YYYY-MM). Pass null for current month.
Statement statement = kobo.accounts().getStatement("acc_123abc", "2024-03");
```

---

## Exceptions

### List Exceptions
```java
// status, cursor, limit
ExceptionPage page = kobo.exceptions().list(ExceptionStatus.OPEN, null, 20);
```

### Resolve Exception
```java
Exception_ exception = kobo.exceptions().resolve("exc_123abc", 
    ResolveExceptionRequest.builder()
        .action("refund")
        .reason("Customer requested refund")
        .build()
);
```
