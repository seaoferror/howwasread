package repository

import (
	"backend/onlineconversation/internal/dto"
	"backend/onlineconversation/internal/projection"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

func (r *Repository) InsertConversation(ctx context.Context, session session, conversationId uuid.UUID, req dto.CreateConversationRequest) error {
	_, err := session.ExecContext(ctx, `
		INSERT INTO online_conversation
		(id, novel, short_story, poem, play, film, written_by, rule, capacity,
		time, length_minutes, current_registrants)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1)`,
		conversationId[:], req.Novel, req.ShortStory, req.Poem, req.Play, req.Film,
		req.WrittenBy, req.Rule,
		req.Capacity, req.Time, req.LengthMinutes)
	if err != nil {
		slog.Error("fail to insert new online conversation", "err", err)
		return err
	}
	return nil
}

func (r *Repository) UpdateConversationIfModerator(ctx context.Context, session session, memberId uuid.UUID, req dto.UpdateConversationRequest) (bool, error) {
	res, err := session.ExecContext(ctx, `
		UPDATE online_conversation
		SET novel=?, short_story=?, poem=?, play=?, film=?,
		written_by=?, rule=?, capacity=?, time=?, length_minutes=?
		WHERE id=?
		AND EXISTS (SELECT 1 FROM online_conversation_moderator
		WHERE conversation_id=? AND member_id=? AND deleted_at IS NULL)`,
		req.Novel, req.ShortStory, req.Poem, req.Play, req.Film,
		req.WrittenBy, req.Rule, req.Capacity, req.Time, req.LengthMinutes,
		req.Id[:], req.Id[:], memberId[:])
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n > 0, err
}

func (r *Repository) DeleteOnlineConversationIfModerator(ctx context.Context, session session, conversationId, memberId uuid.UUID) (bool, error) {
	res, err := session.ExecContext(ctx, `
		UPDATE online_conversation
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id=? AND deleted_at IS NULL
		AND EXISTS (SELECT 1 FROM online_conversation_moderator
		WHERE conversation_id=? AND member_id=? AND deleted_at IS NULL)`,
		conversationId[:], conversationId[:], memberId[:])
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n > 0, err
}

func (r *Repository) AddBanIdIfModerator(ctx context.Context, session session, id, conversationId, modId, banId uuid.UUID) (bool, error) {
	res, err := session.ExecContext(ctx, `
		INSERT INTO online_conversation_ban (id, conversation_id, member_id)
		SELECT ?, ?, ?
		WHERE EXISTS (SELECT 1 FROM online_conversation_moderator
		WHERE conversation_id=? AND member_id=? AND deleted_at IS NULL)
		AND NOT EXISTS (SELECT 1 FROM online_conversation_ban
		WHERE conversation_id=? AND member_id=? AND deleted_at IS NULL)`,
		id[:], conversationId[:], banId[:],
		conversationId[:], modId[:],
		conversationId[:], banId[:])
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n > 0, err
}

func (r *Repository) InsertModerator(ctx context.Context, session session, id, conversationId, memberId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`INSERT INTO online_conversation_moderator (id, conversation_id, member_id)
		VALUES (?, ?, ?)`,
		id[:], conversationId[:], memberId[:],
	)
	if err != nil {
		slog.Error("fail to insert moderator",
			"conversationId", conversationId, "memberId", memberId, "err", err)
		return err
	}
	return nil
}

func (r *Repository) InsertRegistrant(ctx context.Context, session session, id, conversationId, memberId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`INSERT INTO online_conversation_registrant (id, conversation_id, member_id)
		VALUES (?, ?, ?)`,
		id[:], conversationId[:], memberId[:],
	)
	if err != nil {
		slog.Error("fail to insert registrant",
			"conversationId", conversationId, "memberId", memberId, "err", err)
		return err
	}
	return nil
}

func (r *Repository) InsertRegistrantIfNotExist(ctx context.Context, session session, id, conversationId, memberId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`INSERT INTO online_conversation_registrant (id, conversation_id, member_id)
		SELECT ?, ?, ?
		WHERE NOT EXISTS (
		SELECT 1 FROM online_conversation_registrant
		WHERE conversation_id = ? AND member_id = ?)`,
		id[:], conversationId[:], memberId[:],
		conversationId[:], memberId[:],
	)
	if err != nil {
		slog.Error("fail to insert registrant",
			"conversationId", conversationId, "memberId", memberId, "err", err)
		return err
	}
	return nil
}

func (r *Repository) AddNotificationId(ctx context.Context, session session, id, conversationId, memberId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`INSERT INTO online_conversation_notification (id, conversation_id, member_id)
		SELECT ?, ?, ?
		WHERE NOT EXISTS (
		SELECT 1 FROM online_conversation_notification
		WHERE conversation_id = ? AND member_id = ?)`,
		id[:], conversationId[:], memberId[:],
		conversationId[:], memberId[:])
	if err != nil {
		slog.Error("fail to add notification id to online conversation",
			"conversationId", conversationId, "memberId", memberId, "err", err)
		return err
	}
	return nil
}

func (r *Repository) FindConversationDetail(ctx context.Context, session session, conversationId, memberId uuid.UUID) (d projection.Detail, err error) {
	row := session.QueryRowContext(ctx, `
		SELECT novel, short_story, poem, play, film, written_by, rule, capacity,
		time, length_minutes,
		EXISTS(SELECT 1 FROM online_conversation_moderator
		WHERE conversation_id = c.id AND member_id = ? AND deleted_at IS NULL),
		EXISTS(SELECT 1 FROM online_conversation_registrant
		WHERE conversation_id = c.id AND member_id = ? AND deleted_at IS NULL),
		EXISTS(SELECT 1 FROM online_conversation_ban
		WHERE conversation_id = c.id AND member_id = ? AND deleted_at IS NULL),
		EXISTS(SELECT 1 FROM online_conversation_notification
		WHERE conversation_id = c.id AND member_id = ? AND deleted_at IS NULL)
		FROM online_conversation c
		WHERE c.id = ?`,
		memberId[:], memberId[:], memberId[:], memberId[:], conversationId[:],
	)
	err = row.Scan(
		&d.Novel, &d.ShortStory, &d.Poem, &d.Play, &d.Film, &d.WrittenBy, &d.Rule, &d.Capacity, &d.Time, &d.LengthMinutes,
		&d.IsModerator, &d.IsRegistrant, &d.IsBanned, &d.IsNotificationScheduled)
	if err != nil {
		slog.Error("fail to find online conversation detail", "err", err)
		return projection.Detail{}, err
	}
	return d, nil
}

func (r *Repository) DeleteOnlineConversation(ctx context.Context, session session, id uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`UPDATE online_conversation SET
		deleted_at = CURRENT_TIMESTAMP WHERE id = ?`, id[:])
	if err != nil {
		slog.Error("fail to soft delete online conversation", "err", err)
		return err
	}
	return nil
}

func (r *Repository) TryIncrementRegistrants(ctx context.Context, session session, conversationId uuid.UUID) (bool, error) {
	result, err := session.ExecContext(ctx,
		`UPDATE online_conversation SET
		current_registrants = current_registrants + 1
		WHERE id = ? AND current_registrants < capacity`,
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
		`UPDATE online_conversation SET
		current_registrants = current_registrants - 1 WHERE id = ?`,
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
		`UPDATE online_conversation_registrant SET
		deleted_at = CURRENT_TIMESTAMP WHERE conversation_id = ?
		AND member_id = ? AND deleted_at IS NULL`,
		conversationId[:], memberId[:],
	)
	if err != nil {
		slog.Error("fail to remove online conversation registrant id",
			"conversationId", conversationId, "memberId", memberId, "err", err)
		return err
	}
	return nil
}

func (r *Repository) RemoveNotificationId(ctx context.Context, session session, conversationId, memberId uuid.UUID) error {
	_, err := session.ExecContext(ctx,
		`UPDATE online_conversation_notification SET deleted_at = CURRENT_TIMESTAMP
		WHERE conversation_id = ? AND member_id = ? AND deleted_at IS NULL`,
		conversationId[:], memberId[:],
	)
	if err != nil {
		slog.Error("fail to remove online conversation notification id",
			"conversationId", conversationId, "memberId", memberId, "err", err)
		return err
	}
	return nil
}
