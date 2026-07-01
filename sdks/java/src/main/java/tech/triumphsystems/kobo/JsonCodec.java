package tech.triumphsystems.kobo;

import java.util.Map;

/**
 * Minimal, zero-dependency JSON codec using the standard library.
 *
 * <p>Uses hand-rolled serialisation to avoid any external dependency
 * (Jackson, Gson, etc.). Adequate for the well-defined, shallow shapes
 * returned by the Kobo API.
 *
 * @internal
 */
final class JsonCodec {

    // ─── Serialise ────────────────────────────────────────────────────────────

    String toJson(Object obj) {
        if (obj == null) return "null";
        if (obj instanceof String) return quote((String) obj);
        if (obj instanceof Boolean || obj instanceof Number) return obj.toString();
        if (obj instanceof Map<?, ?> m) return mapToJson(m);
        if (obj instanceof Iterable<?> it) return iterableToJson(it);
        // Fall back to class-level reflection-free field enumeration
        return reflectToJson(obj);
    }

    private String quote(String s) {
        return '"' + s
                .replace("\\", "\\\\")
                .replace("\"", "\\\"")
                .replace("\n", "\\n")
                .replace("\r", "\\r")
                .replace("\t", "\\t")
                + '"';
    }

    private String mapToJson(Map<?, ?> m) {
        StringBuilder sb = new StringBuilder("{");
        boolean first = true;
        for (Map.Entry<?, ?> e : m.entrySet()) {
            if (e.getValue() == null) continue; // omit nulls
            if (!first) sb.append(',');
            sb.append(quote(String.valueOf(e.getKey()))).append(':').append(toJson(e.getValue()));
            first = false;
        }
        return sb.append('}').toString();
    }

    private String iterableToJson(Iterable<?> it) {
        StringBuilder sb = new StringBuilder("[");
        boolean first = true;
        for (Object o : it) {
            if (!first) sb.append(',');
            sb.append(toJson(o));
            first = false;
        }
        return sb.append(']').toString();
    }

    /**
     * Uses the object's own {@code toJsonMap()} method (implemented on each model class)
     * so no reflection is needed.
     */
    private String reflectToJson(Object obj) {
        if (obj instanceof JsonSerializable s) return toJson(s.toJsonMap());
        throw new IllegalArgumentException("Cannot serialise " + obj.getClass().getName()
                + " to JSON. Implement JsonSerializable or use a Map.");
    }

    // ─── Deserialise ──────────────────────────────────────────────────────────

    @SuppressWarnings("unchecked")
    <T> T fromJson(String json, Class<T> type) {
        Object raw = new JsonParser(json).parseValue();
        if (type == String.class) return type.cast(String.valueOf(raw));
        if (type == Map.class) return (T) raw;
        if (raw instanceof Map<?, ?> m) return JsonMapper.map((Map<String, Object>) m, type);
        throw new KoboException(0, "parse_error", "Cannot parse response into " + type.getSimpleName(), null);
    }
}
