package tech.triumphsystems.kobo.model;

/** The resolution details applied to a closed exception. */
public final class ExceptionResolution {

    private final ResolutionAction action;
    private final String notes;

    public ExceptionResolution(ResolutionAction action, String notes) {
        this.action = action;
        this.notes = notes;
    }

    public ResolutionAction getAction() { return action; }
    public String getNotes() { return notes; }
}
