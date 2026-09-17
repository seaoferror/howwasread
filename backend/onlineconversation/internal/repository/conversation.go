package repository

import (
	"backend/onlineconversation/internal/entity"
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

const Limit = 10

func uuidSliceToJSON(ids []uuid.UUID) ([]byte, error) {
	hexIds := make([]string, len(ids))
	for i, id := range ids {
		hexIds[i] = id.String()
	}
	return json.Marshal(hexIds)
}

func parseJSONUUIDs(raw string) ([]uuid.UUID, error) {
	var strs []string
	if err := json.Unmarshal([]byte(raw), &strs); err != nil {
		return nil, err
	}
	ids := make([]uuid.UUID, 0, len(strs))
	for _, s := range strs {
		id, err := uuid.Parse(s)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *Repository) InsertConversation(ctx context.Context, session session, conversationId uuid.UUID, novel, shortStory, poem, play, film, writtenBy, rule string, capacity int, t time.Time, length time.Duration, moderatorId uuid.UUID) error {
	moderatorJSON, err := uuidSliceToJSON([]uuid.UUID{moderatorId})
	if err != nil {
		slog.Error("fail to marshal moderator ids", "err", err)
		return err
	}
	emptyJSON := []byte("[]")

	_, err = session.ExecContext(ctx, `
		INSERT INTO conversations
			(id, novel, short_story, poem, play, film, written_by, rule, capacity, time, length_minutes, current_registrants, moderator_ids, ban_ids, reporter_ids, notification_ids)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1, ?, ?, ?, ?)`,
		conversationId[:], novel, shortStory, poem, play, film, writtenBy, rule, capacity, t, int(length.Minutes()),
		moderatorJSON, emptyJSON, emptyJSON, emptyJSON,
	)
	if err != nil {
		slog.Error("fail to insert new conversation", "err", err)
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

func (r *Repository) FindConversation(ctx context.Context, session session, conversationId uuid.UUID) (*entity.Conversation, error) {
	var d entity.Conversation
	var idRaw []byte
	var lengthMinutes int
	var moderatorRaw, banRaw, reporterRaw, notificationRaw string

	row := session.QueryRowContext(ctx, `
		SELECT id, novel, short_story, poem, play, film, written_by, rule, capacity, time, length_minutes,
		       moderator_ids, ban_ids, reporter_ids, notification_ids
		FROM conversations
		WHERE id = ?`,
		conversationId[:],
	)
	err := row.Scan(&idRaw, &d.Novel, &d.ShortStory, &d.Poem, &d.Play, &d.Film, &d.WrittenBy, &d.Rule, &d.Capacity, &d.Time, &lengthMinutes,
		&moderatorRaw, &banRaw, &reporterRaw, &notificationRaw)
	if err != nil {
		slog.Error("fail to find conversation", "err", err)
		return nil, err
	}
	d.Id = uuid.UUID(idRaw)
	d.Length = time.Duration(lengthMinutes) * time.Minute

	d.ModeratorIds, err = parseJSONUUIDs(moderatorRaw)
	if err != nil {
		slog.Error("fail to parse moderator ids", "err", err)
		return nil, err
	}
	d.BanIds, err = parseJSONUUIDs(banRaw)
	if err != nil {
		slog.Error("fail to parse ban ids", "err", err)
		return nil, err
	}
	d.ReporterIds, err = parseJSONUUIDs(reporterRaw)
	if err != nil {
		slog.Error("fail to parse reporter ids", "err", err)
		return nil, err
	}
	d.NotificationIds, err = parseJSONUUIDs(notificationRaw)
	if err != nil {
		slog.Error("fail to parse notification ids", "err", err)
		return nil, err
	}

	// Registrants still come from the separate table
	d.RegistrantIds, err = r.findRegistrantIds(ctx, session, conversationId)
	if err != nil {
		slog.Error("fail to find registrant ids", "err", err)
		return nil, err
	}

	return &d, nil
}

func (r *Repository) findRegistrantIds(ctx context.Context, session session, conversationId uuid.UUID) ([]uuid.UUID, error) {
	rows, err := session.QueryContext(ctx,
		`SELECT member_id FROM conversation_registrants WHERE conversation_id = ?`,
		conversationId[:],
	)
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

func (r *Repository) FindModeratorIds(ctx context.Context, session session, conversationId uuid.UUID) ([]uuid.UUID, error) {
	var raw string
	err := session.QueryRowContext(ctx,
		`SELECT moderator_ids FROM conversations WHERE id = ?`,
		conversationId[:],
	).Scan(&raw)
	if err != nil {
		slog.Error("fail to find moderator ids", "err", err)
		return nil, err
	}
	ids, err := parseJSONUUIDs(raw)
	if err != nil {
		slog.Error("fail to parse moderator ids", "err", err)
		return nil, err
	}
	return ids, nil
}

func (r *Repository) FindReporterIds(ctx context.Context, session session, conversationId uuid.UUID) ([]uuid.UUID, error) {
	var raw string
	err := session.QueryRowContext(ctx,
		`SELECT reporter_ids FROM conversations WHERE id = ?`,
		conversationId[:],
	).Scan(&raw)
	if err != nil {
		slog.Error("fail to find reporter ids", "err", err)
		return nil, err
	}
	ids, err := parseJSONUUIDs(raw)
	if err != nil {
		slog.Error("fail to parse reporter ids", "err", err)
		return nil, err
	}
	return ids, nil
}

func (r *Repository) FindRegistrantIds(ctx context.Context, session session, conversationId uuid.UUID) ([]uuid.UUID, error) {
	ids, err := r.findRegistrantIds(ctx, session, conversationId)
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
		`UPDATE conversations
		 SET ban_ids = JSON_ARRAY_APPEND(ban_ids, '$', CAST(? AS JSON))
		 WHERE id = ?
		   AND NOT JSON_CONTAINS(ban_ids, CAST(? AS JSON))`,
		`"`+banId.String()+`"`, conversationId[:], `"`+banId.String()+`"`,
	)
	if err != nil {
		slog.Error("fail to add ban id to conversation",
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
		`UPDATE conversations
		 SET reporter_ids = JSON_ARRAY_APPEND(reporter_ids, '$', CAST(? AS JSON))
		 WHERE id = ?
		   AND NOT JSON_CONTAINS(reporter_ids, CAST(? AS JSON))`,
		`"`+memberId.String()+`"`, conversationId[:], `"`+memberId.String()+`"`,
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
		`UPDATE conversations
		 SET notification_ids = JSON_ARRAY_APPEND(notification_ids, '$', CAST(? AS JSON))
		 WHERE id = ?
		   AND NOT JSON_CONTAINS(notification_ids, CAST(? AS JSON))`,
		`"`+memberId.String()+`"`, conversationId[:], `"`+memberId.String()+`"`,
	)
	if err != nil {
		slog.Error("fail to add notification id for conversation",
			"conversationId", conversationId, "memberId", memberId, "err", err)
		return err
	}
	return nil
}

func (r *Repository) RemoveNotificationId(ctx context.Context, session session, conversationId, memberId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`UPDATE conversations
		 SET notification_ids = JSON_REMOVE(
		     notification_ids,
		     REPLACE(JSON_SEARCH(notification_ids, 'one', ?), '"', ''))
		 WHERE id = ?
		   AND JSON_SEARCH(notification_ids, 'one', ?) IS NOT NULL`,
		memberId.String(), conversationId[:], memberId.String(),
	)
	if err != nil {
		slog.Error("fail to remove notification id",
			"conversationId", conversationId, "memberId", memberId, "err", err)
		return err
	}
	return nil
}
