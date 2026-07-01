package tech.triumphsystems.kobo;

import java.io.IOException;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.Base64;
import java.util.Map;
import java.util.Objects;

/**
 * HTTP transport layer used by all Kobo service clients.
 * Wraps {@link java.net.http.HttpClient} (added in Java 11) — zero external dependencies.
 *
 * @internal – Not part of the public API surface.
 */
final class KoboHttpClient {

    private static final String SDK_VERSION = "0.1.0";

    private final HttpClient httpClient;
    private final String baseUrl;
    private final String authHeaderValue;
    private final JsonCodec json;

    KoboHttpClient(HttpClient httpClient, String baseUrl, String apiKey, String apiSecret) {
        this.httpClient = Objects.requireNonNull(httpClient);
        this.baseUrl = Objects.requireNonNull(baseUrl);
        this.json = new JsonCodec();

        String raw = apiKey + ":" + apiSecret;
        this.authHeaderValue = "Basic " + Base64.getEncoder().encodeToString(raw.getBytes(StandardCharsets.UTF_8));
    }

    // ─── HTTP Verbs ──────────────────────────────────────────────────────────

    <T> T get(String path, Map<String, String> queryParams, Class<T> responseType) {
        return execute(newRequest("GET", buildUri(path, queryParams), null), responseType);
    }

    <T> T post(String path, Object body, Class<T> responseType) {
        String serialized = body != null ? json.toJson(body) : "";
        return execute(newRequest("POST", buildUri(path, null), serialized), responseType);
    }

    <T> T patch(String path, Object body, Class<T> responseType) {
        String serialized = json.toJson(body);
        return execute(newRequest("PATCH", buildUri(path, null), serialized), responseType);
    }

    // ─── Internals ───────────────────────────────────────────────────────────

    private HttpRequest newRequest(String method, URI uri, String body) {
        HttpRequest.Builder builder = HttpRequest.newBuilder(uri)
                .header("Authorization", authHeaderValue)
                .header("Accept", "application/json")
                .header("User-Agent", "kobo-java/" + SDK_VERSION);

        if (body != null) {
            builder.header("Content-Type", "application/json")
                   .method(method, HttpRequest.BodyPublishers.ofString(body, StandardCharsets.UTF_8));
        } else if ("POST".equals(method) || "PATCH".equals(method)) {
            builder.method(method, HttpRequest.BodyPublishers.noBody());
        } else {
            builder.method(method, HttpRequest.BodyPublishers.noBody());
        }

        return builder.build();
    }

    private <T> T execute(HttpRequest request, Class<T> responseType) {
        HttpResponse<String> response;
        try {
            response = httpClient.send(request, HttpResponse.BodyHandlers.ofString(StandardCharsets.UTF_8));
        } catch (IOException e) {
            throw new KoboException(0, "network_error", "Request failed: " + e.getMessage(), null);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new KoboException(0, "interrupted", "Request was interrupted", null);
        }

        int status = response.statusCode();
        String responseBody = response.body();

        if (status < 200 || status >= 300) {
            // Attempt to decode a structured error body
            try {
                @SuppressWarnings("unchecked")
                Map<String, Object> err = json.fromJson(responseBody, Map.class);
                String code = String.valueOf(err.getOrDefault("code", "http_error"));
                String message = String.valueOf(err.getOrDefault("message", "HTTP " + status));
                @SuppressWarnings("unchecked")
                Map<String, Object> details = (Map<String, Object>) err.get("details");
                throw new KoboException(status, code, message, details);
            } catch (KoboException ke) {
                throw ke;
            } catch (Exception ex) {
                throw new KoboException(status, "http_error", "HTTP " + status + ": " + responseBody, null);
            }
        }

        if (responseType == Void.class || responseBody == null || responseBody.isBlank()) {
            return null;
        }
        return json.fromJson(responseBody, responseType);
    }

    private URI buildUri(String path, Map<String, String> queryParams) {
        StringBuilder sb = new StringBuilder(baseUrl).append(path);
        if (queryParams != null && !queryParams.isEmpty()) {
            sb.append('?');
            queryParams.forEach((k, v) -> {
                if (v != null) {
                    sb.append(encodeComponent(k)).append('=').append(encodeComponent(v)).append('&');
                }
            });
            // Remove trailing '&'
            sb.setLength(sb.length() - 1);
        }
        try {
            return new URI(sb.toString());
        } catch (URISyntaxException e) {
            throw new IllegalArgumentException("Invalid URI: " + sb, e);
        }
    }

    private static String encodeComponent(String s) {
        try {
            return java.net.URLEncoder.encode(s, StandardCharsets.UTF_8)
                    .replace("+", "%20");
        } catch (Exception e) {
            return s;
        }
    }
}
