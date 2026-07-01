package tech.triumphsystems.kobo;

import java.util.*;

/**
 * Hand-rolled recursive descent JSON parser.
 * Parses into: String, Long, Double, Boolean, null, List, Map.
 * Zero external dependencies.
 *
 * @internal
 */
final class JsonParser {

    private final String src;
    private int pos = 0;

    JsonParser(String json) {
        this.src = json == null ? "null" : json.trim();
    }

    Object parseValue() {
        skipWhitespace();
        if (pos >= src.length()) return null;
        char c = src.charAt(pos);
        return switch (c) {
            case '"' -> parseString();
            case '{' -> parseObject();
            case '[' -> parseArray();
            case 't', 'f' -> parseBoolean();
            case 'n' -> parseNull();
            default -> parseNumber();
        };
    }

    // ─── Primitives ───────────────────────────────────────────────────────────

    private String parseString() {
        expect('"');
        StringBuilder sb = new StringBuilder();
        while (pos < src.length()) {
            char c = src.charAt(pos++);
            if (c == '"') return sb.toString();
            if (c == '\\') {
                char esc = src.charAt(pos++);
                sb.append(switch (esc) {
                    case '"' -> '"';
                    case '\\' -> '\\';
                    case '/' -> '/';
                    case 'n' -> '\n';
                    case 'r' -> '\r';
                    case 't' -> '\t';
                    case 'b' -> '\b';
                    case 'f' -> '\f';
                    case 'u' -> (char) Integer.parseInt(src.substring(pos, pos += 4), 16);
                    default -> esc;
                });
            } else {
                sb.append(c);
            }
        }
        throw new KoboException(0, "parse_error", "Unterminated string", null);
    }

    private Boolean parseBoolean() {
        if (src.startsWith("true", pos)) { pos += 4; return Boolean.TRUE; }
        if (src.startsWith("false", pos)) { pos += 5; return Boolean.FALSE; }
        throw new KoboException(0, "parse_error", "Unexpected token at " + pos, null);
    }

    private Object parseNull() {
        if (src.startsWith("null", pos)) { pos += 4; return null; }
        throw new KoboException(0, "parse_error", "Unexpected token at " + pos, null);
    }

    private Number parseNumber() {
        int start = pos;
        boolean isDouble = false;
        if (pos < src.length() && src.charAt(pos) == '-') pos++;
        while (pos < src.length() && Character.isDigit(src.charAt(pos))) pos++;
        if (pos < src.length() && src.charAt(pos) == '.') { isDouble = true; pos++; }
        while (pos < src.length() && Character.isDigit(src.charAt(pos))) pos++;
        if (pos < src.length() && (src.charAt(pos) == 'e' || src.charAt(pos) == 'E')) {
            isDouble = true; pos++;
            if (pos < src.length() && (src.charAt(pos) == '+' || src.charAt(pos) == '-')) pos++;
            while (pos < src.length() && Character.isDigit(src.charAt(pos))) pos++;
        }
        String raw = src.substring(start, pos);
        if (isDouble) return Double.parseDouble(raw);
        try { return Long.parseLong(raw); }
        catch (NumberFormatException e) { return Double.parseDouble(raw); }
    }

    // ─── Composites ───────────────────────────────────────────────────────────

    private Map<String, Object> parseObject() {
        expect('{');
        Map<String, Object> map = new LinkedHashMap<>();
        skipWhitespace();
        if (peek() == '}') { pos++; return map; }
        while (true) {
            skipWhitespace();
            String key = parseString();
            skipWhitespace();
            expect(':');
            skipWhitespace();
            Object value = parseValue();
            map.put(key, value);
            skipWhitespace();
            char next = src.charAt(pos++);
            if (next == '}') return map;
            if (next != ',') throw new KoboException(0, "parse_error", "Expected , or } at " + pos, null);
        }
    }

    private List<Object> parseArray() {
        expect('[');
        List<Object> list = new ArrayList<>();
        skipWhitespace();
        if (peek() == ']') { pos++; return list; }
        while (true) {
            skipWhitespace();
            list.add(parseValue());
            skipWhitespace();
            char next = src.charAt(pos++);
            if (next == ']') return list;
            if (next != ',') throw new KoboException(0, "parse_error", "Expected , or ] at " + pos, null);
        }
    }

    // ─── Helpers ──────────────────────────────────────────────────────────────

    private void skipWhitespace() {
        while (pos < src.length() && Character.isWhitespace(src.charAt(pos))) pos++;
    }

    private char peek() {
        return pos < src.length() ? src.charAt(pos) : '\0';
    }

    private void expect(char c) {
        if (pos >= src.length() || src.charAt(pos) != c) {
            throw new KoboException(0, "parse_error", "Expected '" + c + "' at " + pos + " but got '" + peek() + "'", null);
        }
        pos++;
    }
}
