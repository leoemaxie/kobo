package tech.triumphsystems.kobo.model;

import java.util.Collections;
import java.util.List;

/** Cursor-paginated response wrapper. */
public final class ExceptionPage {

    private final List<Exception_> data;
    private final String nextCursor;

    public ExceptionPage(List<Exception_> data, String nextCursor) {
        this.data = data != null ? data : Collections.emptyList();
        this.nextCursor = nextCursor;
    }

    public List<Exception_> getData() { return data; }
    /** Null when there are no more pages. */
    public String getNextCursor() { return nextCursor; }
    public boolean hasMore() { return nextCursor != null; }
}
