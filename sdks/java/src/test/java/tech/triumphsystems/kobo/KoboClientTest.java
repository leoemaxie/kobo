package tech.triumphsystems.kobo;

import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import tech.triumphsystems.kobo.model.*;
import tech.triumphsystems.kobo.request.*;

import com.sun.net.httpserver.HttpServer;
import java.net.InetSocketAddress;
import java.io.OutputStream;
import java.io.IOException;

import static org.junit.jupiter.api.Assertions.*;

public class KoboClientTest {

    private static HttpServer server;
    private static String baseUrl;
    private static KoboClient client;

    @BeforeAll
    public static void setUp() throws Exception {
        server = HttpServer.create(new InetSocketAddress(0), 0);
        server.createContext("/healthz", exchange -> {
            String response = "{\"status\":\"ok\",\"db\":\"ok\"}";
            exchange.getResponseHeaders().add("Content-Type", "application/json");
            exchange.sendResponseHeaders(200, response.length());
            try (OutputStream os = exchange.getResponseBody()) {
                os.write(response.getBytes());
            }
        });
        server.createContext("/identities", exchange -> {
            if ("POST".equals(exchange.getRequestMethod())) {
                String response = "{\"id\":\"id_123\",\"state\":\"pending\",\"display_name\":\"Test User\"}";
                exchange.getResponseHeaders().add("Content-Type", "application/json");
                exchange.sendResponseHeaders(200, response.length());
                try (OutputStream os = exchange.getResponseBody()) {
                    os.write(response.getBytes());
                }
            } else {
                exchange.sendResponseHeaders(405, -1);
            }
        });
        server.start();
        baseUrl = "http://localhost:" + server.getAddress().getPort();
        client = KoboClient.builder()
                .apiKey("pk_test")
                .apiSecret("sk_test")
                .baseUrl(baseUrl)
                .build();
    }

    @AfterAll
    public static void tearDown() {
        if (server != null) {
            server.stop(0);
        }
    }

    @Test
    public void testHealth() {
        HealthResponse response = client.health();
        assertNotNull(response);
        assertEquals("ok", response.getStatus());
        assertEquals("ok", response.getDb());
    }

    @Test
    public void testIdentitiesCreate() {
        CreateIdentityRequest req = CreateIdentityRequest.builder()
                .displayName("Test User")
                .externalReference("ext_123")
                .build();
        Identity identity = client.identities().create(req);
        assertNotNull(identity);
        assertEquals("id_123", identity.getId());
        assertEquals("Test User", identity.getDisplayName());
        assertEquals(IdentityState.PENDING, identity.getState());
    }
}
