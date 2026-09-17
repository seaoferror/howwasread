package repository

import (
	"backend/onlineconversation/internal/entity"
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

const Limit = 10

func (r *Repository) InsertConversation(ctx context.Context, session session, conversationId uuid.UUID, novel, shortStory, poem, play, film, writtenBy, rule string, capacity int, t time.Time, length time.Duration) error {
	_, err := session.ExecContext(ctx, `
		INSERT INTO conversations
			(id, novel, short_story, poem, play, film, written_by, rule, capacity, time, length_minutes, current_registrants)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1)`,
		conversationId[:], novel, shortStory, poem, play, film, writtenBy, rule, capacity, t, int(length.Minutes()),
	)
	if err != nil {
		slog.Error("fail to insert new conversation", "err", err)
		return err
	}
	return nil
}

func (r *Repository) InsertModerator(ctx context.Context, session session, conversationId, memberId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`INSERT INTO conversation_moderators (conversation_id, member_id) VALUES (?, ?)`,
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
		`INSERT IGNORE INTO conversation_registrants (conversation_id, member_id) VALUES (?, ?)`,
		conversationId[:], memberId[:],
	)
	if err != nil {
		slog.Error("fail to insert registrant",
			"conversationId", conversationId, "memberId", memberId, "err", err)
		return err
	}
	return nil
}

func (r *Repository) FindConversations(ctx context.Context, session session, page int, t time.Time) ([]entity.Conversation, error) {
	rows, err := session.QueryContext(ctx, `
		SELECT id, novel, short_story, poem, play, film, written_by, time
		FROM conversations
		WHERE time > ?
		ORDER BY time ASC
		LIMIT ? OFFSET ?`,
		t.Add(-9*time.Hour), Limit, (page-1)*5,
	)
	if err != nil {
		slog.Error("fail to find next conversations page", "err", err)
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.Conversation, 0, Limit)
	for rows.Next() {
		var d entity.Conversation
		var idRaw []byte
		err = rows.Scan(&idRaw, &d.Novel, &d.ShortStory, &d.Poem, &d.Play, &d.Film, &d.WrittenBy, &d.Time)
		if err != nil {
			slog.Error("fail to scan conversation row", "err", err)
			return nil, err
		}
		d.Id = uuid.UUID(idRaw)
		items = append(items, d)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return items, nil
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
		FROM conversations
		WHERE id = ?`,
		conversationId[:],
	)
	err := row.Scan(&idRaw, &d.Novel, &d.ShortStory, &d.Poem, &d.Play, &d.Film, &d.WrittenBy, &d.Rule, &d.Capacity, &d.Time, &lengthMinutes)
	if err != nil {
		slog.Error("fail to find conversation", "err", err)
		return nil, err
	}
	d.Id = uuid.UUID(idRaw)
	d.Length = time.Duration(lengthMinutes) * time.Minute

	d.ModeratorIds, err = r.findIds(ctx, session, `SELECT member_id FROM conversation_moderators WHERE conversation_id = ?`, conversationId)
	if err != nil {
		slog.Error("fail to find conversation", "err", err)
		return nil, err
	}
	d.RegistrantIds, err = r.findIds(ctx, session, `SELECT member_id FROM conversation_registrants WHERE conversation_id = ?`, conversationId)
	if err != nil {
		slog.Error("fail to find conversation", "err", err)
		return nil, err
	}
	d.BanIds, err = r.findIds(ctx, session, `SELECT member_id FROM conversation_bans WHERE conversation_id = ?`, conversationId)
	if err != nil {
		slog.Error("fail to find conversation", "err", err)
		return nil, err
	}
	d.NotificationIds, err = r.findIds(ctx, session, `SELECT member_id FROM conversation_notifications WHERE conversation_id = ?`, conversationId)
	if err != nil {
		slog.Error("fail to find conversation", "err", err)
		return nil, err
	}
	return &d, nil
}

func (r *Repository) FindModeratorIds(ctx context.Context, session session, conversationId uuid.UUID) ([]uuid.UUID, error) {
	ids, err := r.findIds(ctx, session, `SELECT member_id FROM conversation_moderators WHERE conversation_id = ?`, conversationId)
	if err != nil {
		slog.Error("fail to find mod ids", "err", err)
		return nil, err
	}
	return ids, nil
}

func (r *Repository) FindReporterIds(ctx context.Context, session session, conversationId uuid.UUID) ([]uuid.UUID, error) {
	ids, err := r.findIds(ctx, session, `SELECT member_id FROM conversation_reporters WHERE conversation_id = ?`, conversationId)
	if err != nil {
		slog.Error("fail to find reporter ids", "err", err)
		return nil, err
	}
	return ids, nil
}

func (r *Repository) FindRegistrantIds(ctx context.Context, session session, conversationId uuid.UUID) ([]uuid.UUID, error) {
	ids, err := r.findIds(ctx, session, `SELECT member_id FROM conversation_registrants WHERE conversation_id = ?`, conversationId)
	if err != nil {
		slog.Error("fail to find registrant ids", "err", err)
		return nil, err
	}
	return ids, nil
}

func (r *Repository) FindCapacity(ctx context.Context, session session, id uuid.UUID) (int, error) {
	var capacity int
	err := session.QueryRowContext(ctx, `SELECT capacity FROM conversations WHERE id = ?`, id[:]).Scan(&capacity)
	if err != nil {
		slog.Error("fail to find capacity", "err", err)
		return 0, err
	}
	return capacity, nil
}

func (r *Repository) AddBanId(ctx context.Context, session session, conversationId uuid.UUID, banId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`INSERT IGNORE INTO conversation_bans (conversation_id, member_id) VALUES (?, ?)`,
		conversationId[:], banId[:],
	)
	if err != nil {
		slog.Error("fail to add participant member id to conversation",
			"err", err, "conversationId", conversationId, "memberId", banId.String())
		return err
	}
	return nil
}

func (r *Repository) DeleteOnlineConversation(ctx context.Context, session session, id uuid.UUID) error {
	_, err := session.ExecContext(ctx, `DELETE FROM conversations WHERE id = ?`, id[:])
	if err != nil {
		slog.Error("fail to delete conversation", "err", err)
		return err
	}
	return nil
}

func (r *Repository) AddReporterId(ctx context.Context, session session, conversationId, memberId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`INSERT IGNORE INTO conversation_reporters (conversation_id, member_id) VALUES (?, ?)`,
		conversationId[:], memberId[:],
	)
	if err != nil {
		slog.Error("fail to update reporter ids for conversation",
			"conversationId", conversationId, "memberId", memberId, "err", err)
		return err
	}
	return nil
}

func (r *Repository) LockConversationRegistrants(ctx context.Context, session session, conversationId uuid.UUID) (int, int, error) {
	var currentRegistrants, capacity int
	err := session.QueryRowContext(ctx,
		`SELECT current_registrants, capacity FROM conversations WHERE id = ? FOR UPDATE`,
		conversationId[:],
	).Scan(&currentRegistrants, &capacity)
	if err != nil {
		slog.Error("fail to lock conversation for registration",
			"conversationId", conversationId, "err", err)
		return 0, 0, err
	}
	return currentRegistrants, capacity, nil
}

func (r *Repository) IncrementRegistrants(ctx context.Context, session session, conversationId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`UPDATE conversations SET current_registrants = current_registrants + 1 WHERE id = ?`,
		conversationId[:],
	)
	if err != nil {
		slog.Error("fail to increment registrants",
			"conversationId", conversationId, "err", err)
		return err
	}
	return nil
}

func (r *Repository) DecrementRegistrants(ctx context.Context, session session, conversationId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`UPDATE conversations SET current_registrants = current_registrants - 1 WHERE id = ?`,
		conversationId[:],
	)
	if err != nil {
		slog.Error("fail to decrement registrants",
			"conversationId", conversationId, "err", err)
		return err
	}
	return nil
}

func (r *Repository) RemoveRegistrantId(ctx context.Context, session session, conversationId, memberId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`DELETE FROM conversation_registrants WHERE conversation_id = ? AND member_id = ?`,
		conversationId[:], memberId[:],
	)
	if err != nil {
		slog.Error("fail to remove registrant id",
			"conversationId", conversationId, "memberId", memberId, "err", err)
		return err
	}
	return nil
}

func (r *Repository) AddNotificationId(ctx context.Context, session session, conversationId, memberId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`INSERT IGNORE INTO conversation_notifications (conversation_id, member_id) VALUES (?, ?)`,
		conversationId[:], memberId[:],
	)
	if err != nil {
		slog.Error("fail to update reporter ids for conversation",
			"conversationId", conversationId, "memberId", memberId, "err", err)
		return err
	}
	return nil
}

func (r *Repository) RemoveNotificationId(ctx context.Context, session session, conversationId, memberId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`DELETE FROM conversation_notifications WHERE conversation_id = ? AND member_id = ?`,
		conversationId[:], memberId[:],
	)
	if err != nil {
		slog.Error("fail to remove registrant id",
			"conversationId", conversationId, "memberId", memberId, "err", err)
		return err
	}
	return nil
}
