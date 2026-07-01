package tech.triumphsystems.kobo.model;

import java.util.Collections;
import java.util.List;

/** Cursor-paginated response wrapper. */
public final class TransactionPage {

    private final List<Transaction> data;
    private final String nextCursor;

    public TransactionPage(List<Transaction> data, String nextCursor) {
        this.data = data != null ? data : Collections.emptyList();
        this.nextCursor = nextCursor;
    }

    public List<Transaction> getData() { return data; }
    /** Null when there are no more pages. */
    public String getNextCursor() { return nextCursor; }
    public boolean hasMore() { return nextCursor != null; }
}
