package tech.triumphsystems.kobo.model;

/** Action taken to resolve an exception. */
public enum ResolutionAction {
    RETURN_TO_SENDER("return_to_sender"),
    REDIRECT_TO_SUCCESSOR("redirect_to_successor"),
    MANUAL_OVERRIDE("manual_override");

    private final String value;
    ResolutionAction(String value) { this.value = value; }

    public String getValue() { return value; }

    public static ResolutionAction of(String s) {
        for (ResolutionAction v : values()) if (v.value.equalsIgnoreCase(s)) return v;
        throw new IllegalArgumentException("Unknown ResolutionAction: " + s);
    }
}
