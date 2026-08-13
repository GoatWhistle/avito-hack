package infra

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/games/domain/bukovki"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

type PgWordPool struct {
	pool     *pgxpool.Pool
	fallback *bukovki.FallbackWordPool
}

func NewPgWordPool(pool *pgxpool.Pool) *PgWordPool {
	return &PgWordPool{pool: pool, fallback: bukovki.NewFallbackWordPool()}
}

const vocabularyCTE = `
	WITH vocabulary AS (
		SELECT translate(word, 'ё', 'е') AS word FROM bukovki_words
		UNION
		SELECT translate(token, 'ё', 'е') FROM (
			SELECT DISTINCT lower(t.token) AS token
			FROM items i
			CROSS JOIN LATERAL regexp_split_to_table(lower(i.title), '[^а-яё]+') AS t (token)
			WHERE i.status = 'published'
			  AND i.deleted_at IS NULL
		) AS candidates
		WHERE token !~ '(ый|ой|ая|ое|ые|ие|ий|их|ым|ую|ей|ом|ах|ов|ям|ем|ами|ями|ыми|ью|его|ого|ому)$'
		  AND token NOT IN (SELECT word FROM bukovki_stopwords)
	)`

func (p *PgWordPool) RandomWord(ctx context.Context, minLength, maxLength int) (string, error) {
	const query = vocabularyCTE + `
		SELECT word
		FROM vocabulary v
		WHERE char_length(v.word) BETWEEN $1 AND $2
		  AND EXISTS (
			SELECT 1
			FROM items i
			WHERE i.status = 'published'
			  AND i.deleted_at IS NULL
			  AND translate(lower(i.title), 'ё', 'е') LIKE '%' || v.word || '%'
		  )
		ORDER BY random()
		LIMIT 1`

	var word string

	err := postgres.QuerierFrom(ctx, p.pool).
		QueryRow(ctx, query, minLength, maxLength).
		Scan(&word)
	if errors.Is(err, pgx.ErrNoRows) {
		return p.fallback.RandomWord(ctx, minLength, maxLength)
	}
	if err != nil {
		return "", fmt.Errorf("draw random bukovki word: %w", err)
	}

	return word, nil
}

func (p *PgWordPool) IsValidWord(ctx context.Context, word string, length int) bool {
	return p.fallback.IsValidWord(ctx, bukovki.NormalizeWord(word), length)
}

const minStemLength = 5

const listingsByPatternQuery = `
	SELECT i.display_id, i.title, i.price_kopeks,
	       coalesce((SELECT ph.url FROM item_photos ph
	         WHERE ph.item_id = i.id
	         ORDER BY ph.position
	         LIMIT 1), '') AS photo_url
	FROM items i
	WHERE i.status = 'published'
	  AND i.deleted_at IS NULL
	  AND translate(lower(i.title), 'ё', 'е') LIKE '%' || $1 || '%'
	ORDER BY EXISTS (SELECT 1 FROM item_photos ph WHERE ph.item_id = i.id) DESC, random()
	LIMIT $2`

func (p *PgWordPool) ListingsByWord(ctx context.Context, word string, limit int) ([]bukovki.Listing, error) {
	normalized := bukovki.NormalizeWord(word)

	listings, err := p.listingsByPattern(ctx, normalized, limit)
	if err != nil {
		return nil, err
	}

	if len(listings) > 0 {
		return listings, nil
	}

	stem := WordStem(normalized, minStemLength)
	if stem == normalized {
		return listings, nil
	}

	return p.listingsByPattern(ctx, stem, limit)
}

func (p *PgWordPool) listingsByPattern(ctx context.Context, pattern string, limit int) ([]bukovki.Listing, error) {
	rows, err := postgres.QuerierFrom(ctx, p.pool).Query(ctx, listingsByPatternQuery, pattern, limit)
	if err != nil {
		return nil, fmt.Errorf("query bukovki listings: %w", err)
	}
	defer rows.Close()

	listings := make([]bukovki.Listing, 0, limit)

	for rows.Next() {
		var listing bukovki.Listing

		if err := rows.Scan(&listing.DisplayID, &listing.Title, &listing.PriceKopeks, &listing.PhotoURL); err != nil {
			return nil, fmt.Errorf("scan bukovki listing: %w", err)
		}

		listings = append(listings, listing)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate bukovki listings: %w", err)
	}

	return listings, nil
}

func WordStem(word string, minLength int) string {
	runes := []rune(word)
	if len(runes) <= minLength {
		return word
	}

	trimmed := len(runes) - 1
	if trimmed < minLength {
		trimmed = minLength
	}

	return string(runes[:trimmed])
}
