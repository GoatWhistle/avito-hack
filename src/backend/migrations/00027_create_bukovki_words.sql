-- +goose Up
CREATE TABLE bukovki_stopwords (
    word text PRIMARY KEY
);

INSERT INTO bukovki_stopwords (word)
VALUES
    ('рядом'), ('тёмно'), ('темно'), ('очень'), ('много'), ('почти'),
    ('после'), ('перед'), ('между'), ('около'), ('через'), ('более'),
    ('менее'), ('самый'), ('такой'), ('всего'), ('можно'), ('нужно'),
    ('будет'), ('когда'), ('чтобы'), ('также'), ('сдаю'), ('штуки'),
    ('штука'), ('литра'), ('литры'), ('литров'), ('тонны'), ('тонна'),
    ('соток'), ('сотка'), ('грамм'), ('метра'), ('метры'), ('пачка'),
    ('цвета'), ('цвете'), ('волос'), ('ящика'), ('ящики'), ('кошки'),
    ('кошек'), ('новый'), ('новая'), ('новое'), ('новые'), ('размер'),
    ('размера'), ('комплект'), ('состояние'), ('отличном'), ('идеальном'),
    ('книг'), ('серо'), ('дерева'), ('гаража'), ('студию'), ('комнату'),
    ('квартире'), ('квартиру'), ('плёнке'), ('коробке'), ('укладки'),
    ('мягким'), ('ремнём'), ('зимняя'), ('вечернее'), ('колёса'), ('массив'),
    ('привод'), ('ремонт'), ('клавиша'), ('упаковке'), ('игра'), ('пара')
ON CONFLICT (word) DO NOTHING;

CREATE TABLE bukovki_words (
    word       text PRIMARY KEY,
    source     text NOT NULL DEFAULT 'items' CHECK (source IN ('items', 'seed')),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_bukovki_words_source ON bukovki_words (source);

INSERT INTO bukovki_words (word, source)
SELECT token, 'items'
FROM (
    SELECT DISTINCT lower(t.token) AS token
    FROM items i
    CROSS JOIN LATERAL regexp_split_to_table(lower(i.title), '[^а-яё]+') AS t (token)
    WHERE i.status = 'published'
      AND i.deleted_at IS NULL
) AS candidates
WHERE char_length(token) BETWEEN 4 AND 8
  AND token !~ '(ый|ой|ая|ое|ые|ие|ий|их|ым|ую|ей|ом|ах|ов|ям|ем|ами|ями|ыми|ью|его|ого|ому)$'
  AND token NOT IN (SELECT word FROM bukovki_stopwords)
ON CONFLICT (word) DO NOTHING;

INSERT INTO bukovki_words (word, source)
VALUES
    ('стол', 'seed'),
    ('полка', 'seed'),
    ('диван', 'seed'),
    ('книга', 'seed'),
    ('лампа', 'seed'),
    ('комод', 'seed'),
    ('чехол', 'seed'),
    ('робот', 'seed'),
    ('санки', 'seed'),
    ('набор', 'seed'),
    ('плита', 'seed'),
    ('халат', 'seed'),
    ('кофта', 'seed'),
    ('шарик', 'seed'),
    ('петля', 'seed'),
    ('туфли', 'seed'),
    ('ковёр', 'seed'),
    ('гитара', 'seed'),
    ('кресло', 'seed'),
    ('рюкзак', 'seed'),
    ('коляска', 'seed'),
    ('самокат', 'seed'),
    ('монитор', 'seed'),
    ('шкатулка', 'seed')
ON CONFLICT (word) DO NOTHING;

-- +goose Down
DROP INDEX IF EXISTS idx_bukovki_words_source;
DROP TABLE IF EXISTS bukovki_words;
DROP TABLE IF EXISTS bukovki_stopwords;
