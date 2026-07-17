package exceptions

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type Repository interface {
	ListOpenExceptions(ctx context.Context, integratorID uuid.UUID, limit, offset int32) ([]Exception, error)
	ResolveException(ctx context.Context, id uuid.UUID, integratorID uuid.UUID, action string, notes string, successorIdentityID *uuid.UUID) (Exception, error)
}

type repository struct {
	q *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) Repository {
	return &repository{q: q}
}

func (r *repository) ListOpenExceptions(ctx context.Context, integratorID uuid.UUID, limit, offset int32) ([]Exception, error) {
	rows, err := r.q.ListOpenExceptions(ctx, sqlc.ListOpenExceptionsParams{
		IntegratorID: integratorID,
		Limit:        limit,
		Offset:       offset,
	})
	if err != nil {
		return nil, err
	}

	entries := make([]Exception, len(rows))
	for i, row := range rows {
		var relatedAccountID *uuid.UUID
		if row.RelatedAccountID.Valid {
			id := row.RelatedAccountID.Bytes
			uid, _ := uuid.FromBytes(id[:])
			relatedAccountID = &uid
		}

		var successorIdentityID *uuid.UUID
		if row.SuccessorIdentityID.Valid {
			id := row.SuccessorIdentityID.Bytes
			uid, _ := uuid.FromBytes(id[:])
			successorIdentityID = &uid
		}

		var resolutionAction *string
		if row.ResolutionAction.Valid {
			action := row.ResolutionAction.String
			resolutionAction = &action
		}

		var resolutionNotes *string
		if row.ResolutionNotes.Valid {
			notes := row.ResolutionNotes.String
			resolutionNotes = &notes
		}

		var resolvedAt *time.Time
		if row.ResolvedAt.Valid {
			t := row.ResolvedAt.Time
			resolvedAt = &t
		}

		entries[i] = Exception{
			ID:                  row.ID,
			IntegratorID:        row.IntegratorID,
			Type:                row.Type,
			AmountKobo:          row.AmountKobo,
			MonnifyReference:      row.MonnifyReference,
			RelatedAccountID:    relatedAccountID,
			Status:              row.Status,
			ResolutionAction:    resolutionAction,
			ResolutionNotes:     resolutionNotes,
			SuccessorIdentityID: successorIdentityID,
			DetectedAt:          row.DetectedAt,
			ResolvedAt:          resolvedAt,
			CreatedAt:           row.CreatedAt,
		}
	}
	return entries, nil
}

func (r *repository) ResolveException(ctx context.Context, id uuid.UUID, integratorID uuid.UUID, action string, notes string, successorIdentityID *uuid.UUID) (Exception, error) {
	params := sqlc.ResolveExceptionParams{
		ID:           id,
		IntegratorID: integratorID,
		ResolutionAction: pgtype.Text{
			String: action,
			Valid:  action != "",
		},
		ResolutionNotes: pgtype.Text{
			String: notes,
			Valid:  notes != "",
		},
	}
	if successorIdentityID != nil {
		params.SuccessorIdentityID = pgtype.UUID{
			Bytes: *successorIdentityID,
			Valid: true,
		}
	}

	row, err := r.q.ResolveException(ctx, params)
	if err != nil {
		return Exception{}, err
	}

	var parsedRelatedAccountID *uuid.UUID
	if row.RelatedAccountID.Valid {
		uid, _ := uuid.FromBytes(row.RelatedAccountID.Bytes[:])
		parsedRelatedAccountID = &uid
	}

	var parsedSuccessorIdentityID *uuid.UUID
	if row.SuccessorIdentityID.Valid {
		uid, _ := uuid.FromBytes(row.SuccessorIdentityID.Bytes[:])
		parsedSuccessorIdentityID = &uid
	}

	var resAction *string
	if row.ResolutionAction.Valid {
		a := row.ResolutionAction.String
		resAction = &a
	}

	var resNotes *string
	if row.ResolutionNotes.Valid {
		n := row.ResolutionNotes.String
		resNotes = &n
	}

	var resolvedAt *time.Time
	if row.ResolvedAt.Valid {
		t := row.ResolvedAt.Time
		resolvedAt = &t
	}

	return Exception{
		ID:                  row.ID,
		IntegratorID:        row.IntegratorID,
		Type:                row.Type,
		AmountKobo:          row.AmountKobo,
		MonnifyReference:      row.MonnifyReference,
		RelatedAccountID:    parsedRelatedAccountID,
		Status:              row.Status,
		ResolutionAction:    resAction,
		ResolutionNotes:     resNotes,
		SuccessorIdentityID: parsedSuccessorIdentityID,
		DetectedAt:          row.DetectedAt,
		ResolvedAt:          resolvedAt,
		CreatedAt:           row.CreatedAt,
	}, nil
}
