# Kobo Java SDK

Official Java SDK for the Kobo API. This SDK is built with zero external dependencies and uses the native `java.net.http.HttpClient` introduced in Java 11.

## Features
- Zero external dependencies
- Hand-rolled JSON codec using the standard library
- Thread-safe HTTP Client
- Requires Java 11+

## Installation

Add the dependency to your `pom.xml`:

```xml
<dependency>
    <groupId>tech.triumphsystems</groupId>
    <artifactId>kobo-sdk</artifactId>
    <version>0.1.0</version>
</dependency>
```

### Using Locally (Without Package Manager)

To use the SDK locally without fetching from Maven Central or using a build tool for dependency management:

1. Build the SDK JAR from source:
   ```bash
   cd sdks/java
   mvn clean package -DskipTests
   ```
2. Include the generated JAR file (`target/kobo-sdk-0.1.0.jar`) in your project's classpath. Since the SDK has zero external dependencies, no other JARs are required.

*(Alternatively, you can install it to your local Maven repository by running `mvn clean install`)*

## Quick Start

```java
import tech.triumphsystems.kobo.KoboClient;
import tech.triumphsystems.kobo.KoboException;
import tech.triumphsystems.kobo.model.Identity;

import tech.triumphsystems.kobo.model.TransactionPage;
import tech.triumphsystems.kobo.request.CreateIdentityRequest;

import java.util.Map;

public class App {
    public static void main(String[] args) {
        // Initialize the client
        KoboClient kobo = KoboClient.of("kobo_live_pk_...", "kobo_live_sk_...");

        // Or use the sandbox environment
        // KoboClient kobo = KoboClient.sandbox("kobo_test_pk_...", "kobo_test_sk_...");

        try {
            // Check API health
            System.out.println("Health: " + kobo.health().getStatus());

            // Create an identity
            CreateIdentityRequest req = CreateIdentityRequest.builder()
                .externalReference("customer_12345")
                .displayName("Jane Doe")

                .metadata(Map.of("source", "web_signup"))
                .build();

            Identity identity = kobo.identities().create(req);
            System.out.println("Created Identity: " + identity.getId());

            // Fetch an identity
            Identity fetchedIdentity = kobo.identities().get(identity.getId());
            System.out.println("Identity State: " + fetchedIdentity.getState());

            // List identities
            Identity[] identities = kobo.identities().list(null, 10, 0);
            System.out.println("Found " + identities.length + " identities");

            // List transactions for an account
            if (fetchedIdentity.getVirtualAccount() != null) {
                 // Wait for account to be provisioned, then list txns
                 // Assuming you have the accountId
                 // TransactionPage page = kobo.accounts().listTransactions(accountId, null, 50);
            }

        } catch (KoboException e) {
            System.err.println("API Error: " + e.getCode() + " - " + e.getMessage());
        }
    }
}
```

## Error Handling

The SDK throws a `KoboException` when the API returns a non-2xx response. You should branch your error handling logic on the `code` property, not the `message` string.

```java
try {
    kobo.identities().get("invalid-id");
} catch (KoboException e) {
    if ("not_found".equals(e.getCode())) {
        System.out.println("Identity not found.");
    } else {
        System.out.println("API Error: " + e.getCode());
    }
}
```

## License

MIT
