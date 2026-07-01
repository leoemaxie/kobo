package tech.triumphsystems.kobo;

import java.util.Map;

/**
 * Marker interface for model classes that know how to serialise themselves to
 * a {@link Map} for JSON encoding. Keeps the codec reflection-free.
 *
 * @internal
 */
interface JsonSerializable {
    /** Returns a Map representation suitable for JSON serialisation. */
    Map<String, Object> toJsonMap();
}
