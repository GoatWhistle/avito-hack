-- +goose Up
CREATE TABLE badges (
    id          text PRIMARY KEY,
    name        text NOT NULL,
    description text NOT NULL DEFAULT '',
    icon_url    text NOT NULL DEFAULT ''
);

CREATE TABLE user_badges (
    id        uuid PRIMARY KEY,
    user_id   uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    badge_id  text NOT NULL REFERENCES badges (id),
    earned_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (user_id, badge_id)
);

CREATE INDEX idx_user_badges_user ON user_badges (user_id, earned_at DESC);

CREATE TABLE rewards (
    id              text PRIMARY KEY,
    title           text NOT NULL,
    description     text NOT NULL DEFAULT '',
    kind            text NOT NULL CHECK (kind IN ('promo', 'utility', 'cosmetic')),
    condition_type  text NOT NULL CHECK (condition_type IN ('level', 'streak', 'achievement')),
    condition_value integer NOT NULL DEFAULT 0 CHECK (condition_value >= 0)
);

CREATE TABLE user_rewards (
    id           uuid PRIMARY KEY,
    user_id      uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    reward_id    text NOT NULL REFERENCES rewards (id),
    status       text NOT NULL DEFAULT 'granted'
                 CHECK (status IN ('granted', 'activated', 'expired')),
    code         text,
    granted_at   timestamptz NOT NULL DEFAULT now(),
    activated_at timestamptz,
    expires_at   timestamptz,
    UNIQUE (user_id, reward_id)
);

CREATE INDEX idx_user_rewards_user ON user_rewards (user_id, granted_at DESC);
CREATE UNIQUE INDEX idx_user_rewards_code ON user_rewards (code) WHERE code IS NOT NULL;

CREATE TABLE xp_events (
    id         uuid PRIMARY KEY,
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    action     text NOT NULL,
    subject_id uuid,
    amount     integer NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_xp_events_user_created ON xp_events (user_id, created_at DESC);
-- Частичный уникальный индекс: одно начисление на пару (действие, объект).
-- Чек-ины с subject_id IS NULL под ограничение не попадают и могут повторяться.
CREATE UNIQUE INDEX idx_xp_events_subject ON xp_events (user_id, action, subject_id)
    WHERE subject_id IS NOT NULL;

CREATE TABLE favorites (
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    item_id    uuid NOT NULL REFERENCES items (id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, item_id)
);

CREATE INDEX idx_favorites_user_created ON favorites (user_id, created_at DESC, item_id DESC);
CREATE INDEX idx_favorites_item ON favorites (item_id);

INSERT INTO badges (id, name, description, icon_url) VALUES
    ('explorer',        'Исследователь',        'Добавить в избранное товары из 5 разных категорий',      '/assets/badges/explorer.svg'),
    ('quick_finger',    'Быстрый палец',        'Написать продавцу в течение 60 минут после публикации',  '/assets/badges/quick-finger.svg'),
    ('discount_hunter', 'Охотник за скидками',  'Купить товар из раздела спецпредложений',                '/assets/badges/discount-hunter.svg'),
    ('avito_director',  'Режиссёр Авито',       'Опубликовать 5 объявлений с видео длиннее 30 секунд',    '/assets/badges/director.svg'),
    ('honest_critic',   'Честный критик',       'Оставить 10 подробных отзывов с фотографиями',           '/assets/badges/critic.svg'),
    ('reliable_bro',    'Надёжный бро',         'Оставить 10 развёрнутых отзывов с видео',                '/assets/badges/bro.svg'),
    ('nothing_hidden',  'Ничего не скроешь',    'Заполнить объявление на 100%',                           '/assets/badges/complete.svg'),
    ('lightning',       'Молния',               'Ответить на 10 сообщений быстрее чем за 15 минут',       '/assets/badges/lightning.svg'),
    ('green_planet',    'Зелёная планета',      'Продать 5 б/у вещей',                                    '/assets/badges/green.svg'),
    ('charged_streak',  'Заряженный стрик',     'Заходить в приложение 14 дней подряд',                   '/assets/badges/streak.svg'),
    ('raccoon_friend',  'Друг Енота',           'Достичь 10 уровня',                                      '/assets/badges/friend.svg');

INSERT INTO rewards (id, title, description, kind, condition_type, condition_value) VALUES
    ('starter_status',                     'Статус «Новичок»',               'Доступ к Еноту',                              'cosmetic', 'level', 1),
    ('raccoon_accessory',                  'Кепка и очки для Енота',         'Кастомизация питомца',                        'cosmetic', 'level', 2),
    ('services_discount_10',               'Скидка 10% на «Выделение»',      'Скидка на платные услуги продвижения',        'promo',    'level', 3),
    ('attentive_badge',                    'Бейдж «Внимательный»',           'Бейдж в профиль',                             'cosmetic', 'level', 4),
    ('listing_badge_or_delivery_discount', 'Значок на объявление',           'Значок на 3 дня или −20% на доставку',        'promo',    'level', 5),
    ('raccoon_avatar_frame',               'Цветная рамка аватарки',         'Рамка для аватарки Енота',                    'cosmetic', 'level', 6),
    ('favorites_personal_offer',           'Персональный промокод',          'Спецпредложение на товары из избранного',     'promo',    'level', 7),
    ('dialogue_master_badge',              'Бейдж «Мастер диалога»',         'Бейдж в профиль',                             'cosmetic', 'level', 8),
    ('autoteka_discount_30',               'Скидка 30% на Автотеку',         'Скидка на отчёт о проверке авто',             'promo',    'level', 9),
    ('free_delivery_500',                  'Бесплатная доставка',            'Промокод на доставку до 500 ₽',               'promo',    'level', 10),
    ('xl_listing_discount_50',             'Скидка 50% на XL-объявление',    'Скидка на услугу «XL-объявление»',            'promo',    'level', 11),
    ('treasure_hunter_badge',              'Бейдж «Охотник за сокровищами»', 'Бейдж в профиль',                             'cosmetic', 'level', 12),
    ('views_boost_24h',                    'Увеличение просмотров',          'Бесплатная услуга на 24 часа',                'promo',    'level', 13),
    ('raccoon_sticker_pack',               'Стикер-пачка с Енотом',          'Эксклюзивные стикеры для чатов Авито',        'cosmetic', 'level', 14),
    ('free_delivery_and_avito_guru_badge', 'Купон и бейдж «Гуру Авито»',     'Бесплатная доставка и топ-статус лояльности', 'promo',    'level', 15);

-- +goose Down
DROP INDEX IF EXISTS idx_favorites_item;
DROP INDEX IF EXISTS idx_favorites_user_created;
DROP TABLE IF EXISTS favorites;
DROP INDEX IF EXISTS idx_xp_events_subject;
DROP INDEX IF EXISTS idx_xp_events_user_created;
DROP TABLE IF EXISTS xp_events;
DROP INDEX IF EXISTS idx_user_rewards_code;
DROP INDEX IF EXISTS idx_user_rewards_user;
DROP TABLE IF EXISTS user_rewards;
DROP TABLE IF EXISTS rewards;
DROP INDEX IF EXISTS idx_user_badges_user;
DROP TABLE IF EXISTS user_badges;
DROP TABLE IF EXISTS badges;
