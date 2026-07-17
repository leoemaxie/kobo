package tech.triumphsystems.kobo;

import tech.triumphsystems.kobo.model.*;
import tech.triumphsystems.kobo.model.IdentityState;

import tech.triumphsystems.kobo.model.TransactionStatus;
import tech.triumphsystems.kobo.model.ExceptionType;
import tech.triumphsystems.kobo.model.ExceptionStatus;
import tech.triumphsystems.kobo.model.ResolutionAction;

import java.time.Instant;
import java.util.*;

/**
 * Maps raw JSON maps (produced by {@link JsonParser}) to typed model objects.
 *
 * @internal
 */
@SuppressWarnings("unchecked")
final class JsonMapper {

    static <T> T map(Map<String, Object> m, Class<T> type) {
        if (type == Identity.class)             return type.cast(toIdentity(m));
        if (type == Transaction.class)          return type.cast(toTransaction(m));
        if (type == Statement.class)            return type.cast(toStatement(m));
        if (type == Exception_.class)           return type.cast(toException(m));
        if (type == TransactionPage.class)      return type.cast(toTransactionPage(m));
        if (type == ExceptionPage.class)        return type.cast(toExceptionPage(m));
        if (type == HealthResponse.class)       return type.cast(toHealthResponse(m));
        throw new KoboException(0, "parse_error", "No mapper for " + type.getSimpleName(), null);
    }

    // ─── Identity ─────────────────────────────────────────────────────────────

    static Identity toIdentity(Map<String, Object> m) {
        VirtualAccountSummary va = null;
        if (m.get("virtual_account") instanceof Map<?, ?> vaMap) {
            va = toVirtualAccount((Map<String, Object>) vaMap);
        }
        return new Identity(
            str(m, "id"),
            str(m, "external_reference"),
            str(m, "display_name"),
            IdentityState.of(str(m, "state")),

            va,
            m.containsKey("metadata") ? (Map<String, Object>) m.get("metadata") : null,
            str(m, "failure_reason"),
            instant(m, "created_at"),
            instant(m, "updated_at")
        );
    }

    static VirtualAccountSummary toVirtualAccount(Map<String, Object> m) {
        Long expectedAmountKobo = m.containsKey("expected_amount_kobo") && m.get("expected_amount_kobo") != null 
                ? longVal(m, "expected_amount_kobo") : null;
        boolean isExpired = m.containsKey("is_expired") && Boolean.TRUE.equals(m.get("is_expired"));
        return new VirtualAccountSummary(str(m, "account_number"), str(m, "bank_name"), str(m, "account_name"), expectedAmountKobo, isExpired);
    }

    // ─── Transaction ──────────────────────────────────────────────────────────

    static Transaction toTransaction(Map<String, Object> m) {
        return new Transaction(
            str(m, "id"),
            str(m, "account_id"),
            longVal(m, "amount_kobo"),
            str(m, "direction"),
            TransactionStatus.of(str(m, "status")),
            str(m, "monnify_reference"),
            instant(m, "occurred_at")
        );
    }

    static TransactionPage toTransactionPage(Map<String, Object> m) {
        List<Transaction> data = new ArrayList<>();
        if (m.get("data") instanceof List<?> list) {
            for (Object o : list) {
                if (o instanceof Map<?, ?> row) data.add(toTransaction((Map<String, Object>) row));
            }
        }
        return new TransactionPage(data, str(m, "next_cursor"));
    }

    // ─── Statement ────────────────────────────────────────────────────────────

    static Statement toStatement(Map<String, Object> m) {
        List<Transaction> txns = new ArrayList<>();
        if (m.get("transactions") instanceof List<?> list) {
            for (Object o : list) {
                if (o instanceof Map<?, ?> row) txns.add(toTransaction((Map<String, Object>) row));
            }
        }
        return new Statement(
            str(m, "account_id"),
            str(m, "period"),
            longVal(m, "opening_balance_kobo"),
            longVal(m, "closing_balance_kobo"),
            longVal(m, "total_inflow_kobo"),
            Collections.unmodifiableList(txns)
        );
    }

    // ─── Exception ────────────────────────────────────────────────────────────

    static Exception_ toException(Map<String, Object> m) {
        ExceptionResolution resolution = null;
        if (m.get("resolution") instanceof Map<?, ?> resMap) {
            resolution = new ExceptionResolution(
                ResolutionAction.of(str((Map<String, Object>) resMap, "action")),
                str((Map<String, Object>) resMap, "notes")
            );
        }
        Instant resolvedAt = m.containsKey("resolved_at") && m.get("resolved_at") != null
                ? Instant.parse(str(m, "resolved_at")) : null;
        return new Exception_(
            str(m, "id"),
            ExceptionType.of(str(m, "type")),
            longVal(m, "amount_kobo"),
            str(m, "monnify_reference"),
            str(m, "related_account_id"),
            ExceptionStatus.of(str(m, "status")),
            resolution,
            instant(m, "detected_at"),
            resolvedAt
        );
    }

    static ExceptionPage toExceptionPage(Map<String, Object> m) {
        List<Exception_> data = new ArrayList<>();
        if (m.get("data") instanceof List<?> list) {
            for (Object o : list) {
                if (o instanceof Map<?, ?> row) data.add(toException((Map<String, Object>) row));
            }
        }
        return new ExceptionPage(data, str(m, "next_cursor"));
    }

    // ─── Health ───────────────────────────────────────────────────────────────

    static HealthResponse toHealthResponse(Map<String, Object> m) {
        return new HealthResponse(str(m, "status"), str(m, "db"));
    }

    // ─── Helpers ──────────────────────────────────────────────────────────────

    private static String str(Map<String, Object> m, String key) {
        Object v = m.get(key);
        return v == null ? null : String.valueOf(v);
    }

    private static long longVal(Map<String, Object> m, String key) {
        Object v = m.get(key);
        if (v instanceof Long l) return l;
        if (v instanceof Number n) return n.longValue();
        return 0L;
    }

    private static Instant instant(Map<String, Object> m, String key) {
        String s = str(m, key);
        return s == null ? null : Instant.parse(s);
    }
}
