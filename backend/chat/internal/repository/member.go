package repository

import (
	"context"
	"errors"
	"log/slog"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func (r *Repository) SetServerIP(ctx context.Context, memberId, ip string) error {
	result := r.client.Do(ctx, r.client.B().Sadd().Key(memberId).Member(ip).Build())
	if result.Error() != nil {
		slog.Error("fail to save member ip", "err", result.Error())
		return result.Error()
	}
	return nil
}

func (r *Repository) RemoveServerIP(ctx context.Context, memberId, ip string) error {
	result := r.client.Do(ctx, r.client.B().Srem().Key(memberId).Member(ip).Build())
	if result.Error() != nil {
		slog.Error("fail to remove member ip", "err", result.Error())
		return result.Error()
	}
	return nil
}

func (r *Repository) DeleteAccount(ctx context.Context, id gocql.UUID, email, phoneNumber string) error {
	err := r.session.Batch(gocql.CounterBatch).
		Query("DELETE FROM member_by_id WHERE id = ?", id).
		Query("DELETE FROM profile_by_id WHERE id = ?", id).
		Query("DELETE FROM member_by_email WHERE email = ?", email).
		Query("DELETE FROM member_by_phone_number WHERE phone_number = ?", phoneNumber).
		ExecContext(ctx)
	if err != nil {
		slog.Error("fail to delete account",
			"err", err,
			"id", id.String(),
		)
		return err
	}
	return nil
}

func (r *Repository) DidBlock(ctx context.Context, blockerId gocql.UUID, blockedId gocql.UUID) (bool, error) {
	var a gocql.UUID
	err := r.session.Query(
		`SELECT blocked_id FROM block WHERE blocker_id = ? AND blocked_id = ?`, blockerId, blockedId,
	).ScanContext(ctx, &a)
	if errors.Is(gocql.ErrNotFound, err) {
		return false, nil
	}
	if err != nil {
		slog.Error("fail to check block", "err", err,
			"toId", blockerId,
			"fromId", blockedId)
		return false, err
	}
	return true, nil
}

func (r *Repository) FindChatParticipantIds(ctx context.Context, roomId gocql.UUID) (ids []gocql.UUID, err error) {
	err = r.session.Query(
		`SELECT participant_ids FROM chat_room_by_id WHERE id = ?`, roomId).
		ScanContext(ctx, &ids)
	if err != nil {
		slog.Error("fail to find chat participant ids", "err", err,
			"roomId", roomId)
		return nil, err
	}
	return ids, nil
}

func (r *Repository) AddReporterIdByReportedId(ctx context.Context, reporterId gocql.UUID, reportedId gocql.UUID) error {
	err := r.session.Query(`UPDATE reporter_ids USING TTL ? SET reporter_ids = reporter_ids + ? WHERE reported_id = ?`,
		30*24*60*60, []gocql.UUID{reporterId}, reportedId).ExecContext(ctx)
	if err != nil {
		slog.Error("fail to add reporter id by reported id",
			"err", err,
			"reporterId", reporterId,
			"reportedId", reportedId)
		return err
	}
	return nil
}

func (r *Repository) FindReportCountById(ctx context.Context, reportedId gocql.UUID) (count int, err error) {
	err = r.session.Query(
		`SELECT collection_count(reporter_ids) FROM reporter_ids WHERE reported_id = ?`,
		reportedId,
	).ScanContext(ctx, &count)
	if err != nil {
		slog.Error("fail to find report count by id",
			"err", err,
			"reportedId", reportedId)
		return 0, err
	}
	return count, nil
}

func (r *Repository) FindEmailAndPhoneNumberById(ctx context.Context, id gocql.UUID) (email, phoneNumber string, err error) {
	err = r.session.Query("SELECT email, phone_number FROM member_by_id WHERE id = ?", id).
		ScanContext(ctx, &email, &phoneNumber)
	if err != nil {
		slog.Error("fail to find email and phone number by id",
			"err", err,
			"id", id)
		return "", "", err
	}
	return email, phoneNumber, nil
}

func (r *Repository) BanPhoneNumber(ctx context.Context, phoneNumber string) error {
	err := r.session.Query(`INSERT INTO banned_phone_number (phone_number) VALUES (?)`, phoneNumber).ExecContext(ctx)
	if err != nil {
		slog.Error("fail to ban phone number",
			"err", err)
		return err
	}
	return nil
}
