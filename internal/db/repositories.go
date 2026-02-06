package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hidatara-ds/evolipia-radar/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SourceRepository struct {
	db *DB
}

func NewSourceRepository(db *DB) *SourceRepository {
	return &SourceRepository{db: db}
}

func (r *SourceRepository) List(ctx context.Context) ([]models.Source, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, name, type, category, url, mapping_json, enabled, status,
		       last_test_status, last_test_message, created_at, updated_at
		FROM sources
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []models.Source
	for rows.Next() {
		var s models.Source
		var mappingJSON []byte
		err := rows.Scan(
			&s.ID, &s.Name, &s.Type, &s.Category, &s.URL, &mappingJSON,
			&s.Enabled, &s.Status, &s.LastTestStatus, &s.LastTestMessage,
			&s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if mappingJSON != nil {
			s.MappingJSON = mappingJSON
		}
		sources = append(sources, s)
	}
	return sources, rows.Err()
}

func (r *SourceRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Source, error) {
	var s models.Source
	var mappingJSON []byte
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, name, type, category, url, mapping_json, enabled, status,
		       last_test_status, last_test_message, created_at, updated_at
		FROM sources
		WHERE id = $1
	`, id).Scan(
		&s.ID, &s.Name, &s.Type, &s.Category, &s.URL, &mappingJSON,
		&s.Enabled, &s.Status, &s.LastTestStatus, &s.LastTestMessage,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if mappingJSON != nil {
		s.MappingJSON = mappingJSON
	}
	return &s, nil
}

func (r *SourceRepository) GetEnabled(ctx context.Context) ([]models.Source, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, name, type, category, url, mapping_json, enabled, status,
		       last_test_status, last_test_message, created_at, updated_at
		FROM sources
		WHERE enabled = true
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []models.Source
	for rows.Next() {
		var s models.Source
		var mappingJSON []byte
		err := rows.Scan(
			&s.ID, &s.Name, &s.Type, &s.Category, &s.URL, &mappingJSON,
			&s.Enabled, &s.Status, &s.LastTestStatus, &s.LastTestMessage,
			&s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if mappingJSON != nil {
			s.MappingJSON = mappingJSON
		}
		sources = append(sources, s)
	}
	return sources, rows.Err()
}

func (r *SourceRepository) Create(ctx context.Context, s *models.Source) error {
	var mappingJSON []byte
	if s.MappingJSON != nil {
		mappingJSON = s.MappingJSON
	}
	err := r.db.Pool.QueryRow(ctx, `
		INSERT INTO sources (name, type, category, url, mapping_json, enabled, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`, s.Name, s.Type, s.Category, s.URL, mappingJSON, s.Enabled, s.Status).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt,
	)
	return err
}

func (r *SourceRepository) UpdateTestStatus(ctx context.Context, id uuid.UUID, status, message string) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE sources
		SET last_test_status = $1, last_test_message = $2, updated_at = now()
		WHERE id = $3
	`, status, message, id)
	return err
}

func (r *SourceRepository) SetEnabled(ctx context.Context, id uuid.UUID, enabled bool, status string) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE sources
		SET enabled = $1, status = $2, updated_at = now()
		WHERE id = $3
	`, enabled, status, id)
	return err
}

type ItemRepository struct {
	db *DB
}

func NewItemRepository(db *DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	var item models.Item
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, source_id, title, url, published_at, content_hash,
		       domain, category, raw_excerpt, created_at
		FROM items
		WHERE id = $1
	`, id).Scan(
		&item.ID, &item.SourceID, &item.Title, &item.URL, &item.PublishedAt,
		&item.ContentHash, &item.Domain, &item.Category, &item.RawExcerpt,
		&item.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ItemRepository) GetByContentHash(ctx context.Context, hash string) (*models.Item, error) {
	var item models.Item
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, source_id, title, url, published_at, content_hash,
		       domain, category, raw_excerpt, created_at
		FROM items
		WHERE content_hash = $1
	`, hash).Scan(
		&item.ID, &item.SourceID, &item.Title, &item.URL, &item.PublishedAt,
		&item.ContentHash, &item.Domain, &item.Category, &item.RawExcerpt,
		&item.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ItemRepository) Create(ctx context.Context, item *models.Item) error {
	err := r.db.Pool.QueryRow(ctx, `
		INSERT INTO items (source_id, title, url, published_at, content_hash,
		                   domain, category, raw_excerpt)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`, item.SourceID, item.Title, item.URL, item.PublishedAt, item.ContentHash,
		item.Domain, item.Category, item.RawExcerpt).Scan(
		&item.ID, &item.CreatedAt,
	)
	return err
}

func (r *ItemRepository) GetTopDaily(ctx context.Context, date time.Time, topic *string, limit int) ([]models.Item, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := `
		SELECT i.id, i.source_id, i.title, i.url, i.published_at, i.content_hash,
		       i.domain, i.category, i.raw_excerpt, i.created_at
		FROM items i
		JOIN scores s ON s.item_id = i.id
		WHERE i.published_at >= $1 AND i.published_at < $2
	`
	args := []interface{}{startOfDay, endOfDay}
	argIdx := 3

	if topic != nil {
		query += fmt.Sprintf(` AND EXISTS (
			SELECT 1 FROM summaries su
			WHERE su.item_id = i.id
			AND su.tags @> $%d::jsonb
		)`, argIdx)
		args = append(args, fmt.Sprintf(`["%s"]`, *topic))
		argIdx++
	}

	query += ` ORDER BY s.final DESC LIMIT $` + fmt.Sprintf("%d", argIdx)
	args = append(args, limit)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ID, &item.SourceID, &item.Title, &item.URL, &item.PublishedAt,
			&item.ContentHash, &item.Domain, &item.Category, &item.RawExcerpt,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *ItemRepository) GetRising(ctx context.Context, window time.Duration, limit int) ([]models.Item, error) {
	cutoff := time.Now().Add(-window)

	rows, err := r.db.Pool.Query(ctx, `
		SELECT DISTINCT ON (i.id) i.id, i.source_id, i.title, i.url, i.published_at,
		       i.content_hash, i.domain, i.category, i.raw_excerpt, i.created_at
		FROM items i
		JOIN signals sig ON sig.item_id = i.id
		WHERE sig.fetched_at >= $1
		ORDER BY i.id, sig.fetched_at DESC
		LIMIT $2
	`, cutoff, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ID, &item.SourceID, &item.Title, &item.URL, &item.PublishedAt,
			&item.ContentHash, &item.Domain, &item.Category, &item.RawExcerpt,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *ItemRepository) Search(ctx context.Context, query string, topic *string, limit, offset int) ([]models.Item, int, error) {
	searchQuery := `%` + query + `%`

	baseQuery := `
		SELECT i.id, i.source_id, i.title, i.url, i.published_at, i.content_hash,
		       i.domain, i.category, i.raw_excerpt, i.created_at
		FROM items i
		WHERE (i.title ILIKE $1 OR i.raw_excerpt ILIKE $1)
	`
	args := []interface{}{searchQuery}
	argIdx := 2

	if topic != nil {
		baseQuery += fmt.Sprintf(` AND EXISTS (
			SELECT 1 FROM summaries su
			WHERE su.item_id = i.id
			AND su.tags @> $%d::jsonb
		)`, argIdx)
		args = append(args, fmt.Sprintf(`["%s"]`, *topic))
		argIdx++
	}

	baseQuery += ` ORDER BY i.published_at DESC`

	countQuery := `SELECT COUNT(*) FROM (` + baseQuery + `) sub`
	var total int
	err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	baseQuery += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Pool.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ID, &item.SourceID, &item.Title, &item.URL, &item.PublishedAt,
			&item.ContentHash, &item.Domain, &item.Category, &item.RawExcerpt,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

// GetItemsNeedingScoring returns items that need score computation or recalculation
func (r *ItemRepository) GetItemsNeedingScoring(ctx context.Context, days int, limit int) ([]models.Item, error) {
	rows, err := r.db.Pool.Query(ctx, `
		SELECT i.id, i.source_id, i.title, i.url, i.published_at, i.content_hash,
		       i.domain, i.category, i.raw_excerpt, i.created_at
		FROM items i
		LEFT JOIN scores s ON s.item_id = i.id
		WHERE i.published_at >= now() - interval '1 day' * $1
		AND (s.item_id IS NULL OR s.computed_at < i.created_at)
		LIMIT $2
	`, days, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ID, &item.SourceID, &item.Title, &item.URL, &item.PublishedAt,
			&item.ContentHash, &item.Domain, &item.Category, &item.RawExcerpt,
			&item.CreatedAt,
		)
		if err != nil {
			continue
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

type SignalRepository struct {
	db *DB
}

func NewSignalRepository(db *DB) *SignalRepository {
	return &SignalRepository{db: db}
}

func (r *SignalRepository) Create(ctx context.Context, signal *models.Signal) error {
	err := r.db.Pool.QueryRow(ctx, `
		INSERT INTO signals (item_id, points, comments, rank_pos)
		VALUES ($1, $2, $3, $4)
		RETURNING id, fetched_at
	`, signal.ItemID, signal.Points, signal.Comments, signal.RankPos).Scan(
		&signal.ID, &signal.FetchedAt,
	)
	return err
}

func (r *SignalRepository) GetLatestByItemID(ctx context.Context, itemID uuid.UUID) (*models.Signal, error) {
	var sig models.Signal
	err := r.db.Pool.QueryRow(ctx, `
		SELECT id, item_id, points, comments, rank_pos, fetched_at
		FROM signals
		WHERE item_id = $1
		ORDER BY fetched_at DESC
		LIMIT 1
	`, itemID).Scan(
		&sig.ID, &sig.ItemID, &sig.Points, &sig.Comments, &sig.RankPos, &sig.FetchedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sig, nil
}

func (r *SignalRepository) GetRisingSignals(ctx context.Context, itemID uuid.UUID, window time.Duration) ([]models.Signal, error) {
	cutoff := time.Now().Add(-window)
	rows, err := r.db.Pool.Query(ctx, `
		SELECT id, item_id, points, comments, rank_pos, fetched_at
		FROM signals
		WHERE item_id = $1 AND fetched_at >= $2
		ORDER BY fetched_at ASC
	`, itemID, cutoff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var signals []models.Signal
	for rows.Next() {
		var sig models.Signal
		err := rows.Scan(
			&sig.ID, &sig.ItemID, &sig.Points, &sig.Comments, &sig.RankPos, &sig.FetchedAt,
		)
		if err != nil {
			return nil, err
		}
		signals = append(signals, sig)
	}
	return signals, rows.Err()
}

type ScoreRepository struct {
	db *DB
}

func NewScoreRepository(db *DB) *ScoreRepository {
	return &ScoreRepository{db: db}
}

func (r *ScoreRepository) Upsert(ctx context.Context, score *models.Score) error {
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO scores (item_id, hot, relevance, credibility, novelty, final)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (item_id) DO UPDATE SET
			hot = EXCLUDED.hot,
			relevance = EXCLUDED.relevance,
			credibility = EXCLUDED.credibility,
			novelty = EXCLUDED.novelty,
			final = EXCLUDED.final,
			computed_at = now()
	`, score.ItemID, score.Hot, score.Relevance, score.Credibility, score.Novelty, score.Final)
	return err
}

func (r *ScoreRepository) GetByItemID(ctx context.Context, itemID uuid.UUID) (*models.Score, error) {
	var score models.Score
	err := r.db.Pool.QueryRow(ctx, `
		SELECT item_id, hot, relevance, credibility, novelty, final, computed_at
		FROM scores
		WHERE item_id = $1
	`, itemID).Scan(
		&score.ItemID, &score.Hot, &score.Relevance, &score.Credibility,
		&score.Novelty, &score.Final, &score.ComputedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &score, nil
}

type SummaryRepository struct {
	db *DB
}

func NewSummaryRepository(db *DB) *SummaryRepository {
	return &SummaryRepository{db: db}
}

func (r *SummaryRepository) Upsert(ctx context.Context, summary *models.Summary) error {
	tagsJSON, err := json.Marshal(summary.Tags)
	if err != nil {
		return err
	}
	_, err = r.db.Pool.Exec(ctx, `
		INSERT INTO summaries (item_id, tldr, why_it_matters, tags, method)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (item_id) DO UPDATE SET
			tldr = EXCLUDED.tldr,
			why_it_matters = EXCLUDED.why_it_matters,
			tags = EXCLUDED.tags,
			method = EXCLUDED.method
	`, summary.ItemID, summary.TLDR, summary.WhyItMatters, tagsJSON, summary.Method)
	return err
}

func (r *SummaryRepository) GetByItemID(ctx context.Context, itemID uuid.UUID) (*models.Summary, error) {
	var summary models.Summary
	var tagsJSON []byte
	err := r.db.Pool.QueryRow(ctx, `
		SELECT item_id, tldr, why_it_matters, tags, method, created_at
		FROM summaries
		WHERE item_id = $1
	`, itemID).Scan(
		&summary.ItemID, &summary.TLDR, &summary.WhyItMatters,
		&tagsJSON, &summary.Method, &summary.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(tagsJSON, &summary.Tags); err != nil {
		return nil, err
	}
	return &summary, nil
}

type FetchRunRepository struct {
	db *DB
}

func NewFetchRunRepository(db *DB) *FetchRunRepository {
	return &FetchRunRepository{db: db}
}

func (r *FetchRunRepository) Create(ctx context.Context, run *models.FetchRun) error {
	err := r.db.Pool.QueryRow(ctx, `
		INSERT INTO fetch_runs (source_id, status, error, items_fetched, items_inserted)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, fetched_at
	`, run.SourceID, run.Status, run.Error, run.ItemsFetched, run.ItemsInserted).Scan(
		&run.ID, &run.FetchedAt,
	)
	return err
}
