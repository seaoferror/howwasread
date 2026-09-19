package repository

import (
	"backend/onlineconversation/internal/entity"
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

func (r *Repository) InsertConversation(ctx context.Context, session session, conversationId uuid.UUID, novel, shortStory, poem, play, film, writtenBy, rule string, capacity int, t time.Time, length time.Duration) error {
	_, err := session.ExecContext(ctx, `
		INSERT INTO online_conversation
			(id, novel, short_story, poem, play, film, written_by, rule, capacity, time, length_minutes, current_registrants)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1)`,
		conversationId[:], novel, shortStory, poem, play, film, writtenBy, rule, capacity, t, int(length.Minutes()),
	)
	if err != nil {
		slog.Error("fail to insert new online conversation", "err", err)
		return err
	}
	return nil
}

func (r *Repository) InsertModerator(ctx context.Context, session session, conversationId, memberId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`INSERT INTO online_conversation_moderator (conversation_id, member_id) VALUES (?, ?)`,
		conversationId[:], memberId[:],
	)
	if err != nil {
		slog.Error("fail to insert moderator",
			"conversationId", conversationId, "memberId", memberId, "err", err)
		return err
	}
	return nil
}

func (r *Repository) InsertRegistrant(ctx context.Context, session session, conversationId, memberId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`INSERT IGNORE INTO online_conversation_registrant (conversation_id, member_id) VALUES (?, ?)`,
		conversationId[:], memberId[:],
	)
	if err != nil {
		slog.Error("fail to insert registrant",
			"conversationId", conversationId, "memberId", memberId, "err", err)
		return err
	}
	return nil
}

func (r *Repository) findIds(ctx context.Context, session session, query string, conversationId uuid.UUID) ([]uuid.UUID, error) {
	rows, err := session.QueryContext(ctx, query, conversationId[:])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]uuid.UUID, 0)
	for rows.Next() {
		var rawId []byte
		if err := rows.Scan(&rawId); err != nil {
			return nil, err
		}
		ids = append(ids, uuid.UUID(rawId))
	}
	return ids, rows.Err()
}

func (r *Repository) FindConversation(ctx context.Context, session session, conversationId uuid.UUID) (*entity.Conversation, error) {
	var d entity.Conversation
	var idRaw []byte
	var lengthMinutes int

	row := session.QueryRowContext(ctx, `
		SELECT id, novel, short_story, poem, play, film, written_by, rule, capacity, time, length_minutes
		FROM online_conversation
		WHERE id = ?`,
		conversationId[:],
	)
	err := row.Scan(&idRaw, &d.Novel, &d.ShortStory, &d.Poem, &d.Play, &d.Film, &d.WrittenBy, &d.Rule, &d.Capacity, &d.Time, &lengthMinutes)
	if err != nil {
		slog.Error("fail to find online conversation", "err", err)
		return nil, err
	}
	d.Id = uuid.UUID(idRaw)
	d.Length = time.Duration(lengthMinutes) * time.Minute

	return &d, nil
}

func (r *Repository) FindConversationDetail(ctx context.Context, session session, conversationId, memberId uuid.UUID) (*entity.Conversation, bool, bool, bool, error) {
	var d entity.Conversation
	var idRaw []byte
	var lengthMinutes int
	var isRegistrant, isBanned, isNotificationScheduled bool

	row := session.QueryRowContext(ctx, `
		SELECT id, novel, short_story, poem, play, film, written_by, rule, capacity, time, length_minutes,
			EXISTS(SELECT 1 FROM online_conversation_registrant WHERE conversation_id = c.id AND member_id = ?),
			EXISTS(SELECT 1 FROM online_conversation_ban WHERE conversation_id = c.id AND member_id = ?),
			EXISTS(SELECT 1 FROM online_conversation_notification WHERE conversation_id = c.id AND member_id = ?)
		FROM online_conversation c
		WHERE c.id = ?`,
		memberId[:], memberId[:], memberId[:], conversationId[:],
	)
	err := row.Scan(&idRaw, &d.Novel, &d.ShortStory, &d.Poem, &d.Play, &d.Film, &d.WrittenBy, &d.Rule, &d.Capacity, &d.Time, &lengthMinutes,
		&isRegistrant, &isBanned, &isNotificationScheduled)
	if err != nil {
		slog.Error("fail to find online conversation detail", "err", err)
		return nil, false, false, false, err
	}
	d.Id = uuid.UUID(idRaw)
	d.Length = time.Duration(lengthMinutes) * time.Minute

	return &d, isRegistrant, isBanned, isNotificationScheduled, nil
}

func (r *Repository) FindModeratorIds(ctx context.Context, session session, conversationId uuid.UUID) ([]uuid.UUID, error) {
	ids, err := r.findIds(ctx, session, `SELECT member_id FROM online_conversation_moderator WHERE conversation_id = ?`, conversationId)
	if err != nil {
		slog.Error("fail to find mod ids", "err", err)
		return nil, err
	}
	return ids, nil
}

func (r *Repository) IsModerator(ctx context.Context, session session, conversationId, memberId uuid.UUID) (bool, error) {
	var exists bool
	err := session.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM online_conversation_moderator WHERE conversation_id = ? AND member_id = ?)`,
		conversationId[:], memberId[:],
	).Scan(&exists)
	if err != nil {
		slog.Error("fail to check moderator", "conversationId", conversationId, "memberId", memberId, "err", err)
		return false, err
	}
	return exists, nil
}

func (r *Repository) IsRegistrant(ctx context.Context, session session, conversationId, memberId uuid.UUID) (bool, error) {
	var exists bool
	err := session.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM online_conversation_registrant WHERE conversation_id = ? AND member_id = ?)`,
		conversationId[:], memberId[:],
	).Scan(&exists)
	if err != nil {
		slog.Error("fail to check registrant", "conversationId", conversationId, "memberId", memberId, "err", err)
		return false, err
	}
	return exists, nil
}

func (r *Repository) IsBanned(ctx context.Context, session session, conversationId, memberId uuid.UUID) (bool, error) {
	var exists bool
	err := session.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM online_conversation_ban WHERE conversation_id = ? AND member_id = ?)`,
		conversationId[:], memberId[:],
	).Scan(&exists)
	if err != nil {
		slog.Error("fail to check ban", "conversationId", conversationId, "memberId", memberId, "err", err)
		return false, err
	}
	return exists, nil
}

func (r *Repository) IsNotificationScheduled(ctx context.Context, session session, conversationId, memberId uuid.UUID) (bool, error) {
	var exists bool
	err := session.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM online_conversation_notification WHERE conversation_id = ? AND member_id = ?)`,
		conversationId[:], memberId[:],
	).Scan(&exists)
	if err != nil {
		slog.Error("fail to check notification", "conversationId", conversationId, "memberId", memberId, "err", err)
		return false, err
	}
	return exists, nil
}

func (r *Repository) HasNotification(ctx context.Context, session session, conversationId uuid.UUID) (bool, error) {
	var exists bool
	err := session.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM online_conversation_notification WHERE conversation_id = ?)`,
		conversationId[:],
	).Scan(&exists)
	if err != nil {
		slog.Error("fail to check notification existence", "conversationId", conversationId, "err", err)
		return false, err
	}
	return exists, nil
}

func (r *Repository) FindReporterIds(ctx context.Context, session session, conversationId uuid.UUID) ([]uuid.UUID, error) {
	ids, err := r.findIds(ctx, session, `SELECT member_id FROM online_conversation_reporter WHERE conversation_id = ?`, conversationId)
	if err != nil {
		slog.Error("fail to find reporter ids", "err", err)
		return nil, err
	}
	return ids, nil
}

func (r *Repository) FindRegistrantIds(ctx context.Context, session session, conversationId uuid.UUID) ([]uuid.UUID, error) {
	ids, err := r.findIds(ctx, session, `SELECT member_id FROM online_conversation_registrant WHERE conversation_id = ?`, conversationId)
	if err != nil {
		slog.Error("fail to find registrant ids", "err", err)
		return nil, err
	}
	return ids, nil
}

func (r *Repository) FindCapacity(ctx context.Context, session session, id uuid.UUID) (int, error) {
	var capacity int
	err := session.QueryRowContext(ctx, `SELECT capacity FROM online_conversation WHERE id = ?`, id[:]).Scan(&capacity)
	if err != nil {
		slog.Error("fail to find online conversation capacity", "err", err)
		return 0, err
	}
	return capacity, nil
}

func (r *Repository) AddBanId(ctx context.Context, session session, conversationId uuid.UUID, banId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`INSERT IGNORE INTO online_conversation_ban (conversation_id, member_id) VALUES (?, ?)`,
		conversationId[:], banId[:],
	)
	if err != nil {
		slog.Error("fail to add ban id to online conversation",
			"err", err, "conversationId", conversationId, "memberId", banId.String())
		return err
	}
	return nil
}

func (r *Repository) DeleteOnlineConversation(ctx context.Context, session session, id uuid.UUID) error {
	_, err := session.ExecContext(ctx, `UPDATE online_conversation SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?`, id[:])
	if err != nil {
		slog.Error("fail to soft delete online conversation", "err", err)
		return err
	}
	return nil
}

func (r *Repository) AddReporterId(ctx context.Context, session session, conversationId, memberId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`INSERT IGNORE INTO online_conversation_reporter (conversation_id, member_id) VALUES (?, ?)`,
		conversationId[:], memberId[:],
	)
	if err != nil {
		slog.Error("fail to add reporter id to online conversation",
			"conversationId", conversationId, "memberId", memberId, "err", err)
		return err
	}
	return nil
}

func (r *Repository) TryIncrementRegistrants(ctx context.Context, session session, conversationId uuid.UUID) (bool, error) {
	result, err := session.ExecContext(ctx,
		`UPDATE online_conversation SET current_registrants = current_registrants + 1 WHERE id = ? AND current_registrants < capacity`,
		conversationId[:],
	)
	if err != nil {
		slog.Error("fail to increment online conversation registrants",
			"conversationId", conversationId, "err", err)
		return false, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

func (r *Repository) DecrementRegistrants(ctx context.Context, session session, conversationId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`UPDATE online_conversation SET current_registrants = current_registrants - 1 WHERE id = ?`,
		conversationId[:],
	)
	if err != nil {
		slog.Error("fail to decrement online conversation registrants",
			"conversationId", conversationId, "err", err)
		return err
	}
	return nil
}

func (r *Repository) RemoveRegistrantId(ctx context.Context, session session, conversationId, memberId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`DELETE FROM online_conversation_registrant WHERE conversation_id = ? AND member_id = ?`,
		conversationId[:], memberId[:],
	)
	if err != nil {
		slog.Error("fail to remove online conversation registrant id",
			"conversationId", conversationId, "memberId", memberId, "err", err)
		return err
	}
	return nil
}

func (r *Repository) AddNotificationId(ctx context.Context, session session, conversationId, memberId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`INSERT IGNORE INTO online_conversation_notification (conversation_id, member_id) VALUES (?, ?)`,
		conversationId[:], memberId[:],
	)
	if err != nil {
		slog.Error("fail to add notification id to online conversation",
			"conversationId", conversationId, "memberId", memberId, "err", err)
		return err
	}
	return nil
}

func (r *Repository) RemoveNotificationId(ctx context.Context, session session, conversationId, memberId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`DELETE FROM online_conversation_notification WHERE conversation_id = ? AND member_id = ?`,
		conversationId[:], memberId[:],
	)
	if err != nil {
		slog.Error("fail to remove online conversation notification id",
			"conversationId", conversationId, "memberId", memberId, "err", err)
		return err
	}
	return nil
}
