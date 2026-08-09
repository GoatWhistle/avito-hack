-- +goose Up
DELETE FROM user_badges
WHERE badge_id IN (
    'quick_finger', 'discount_hunter', 'avito_director',
    'honest_critic', 'reliable_bro', 'lightning'
);

DELETE FROM badges
WHERE id IN (
    'quick_finger', 'discount_hunter', 'avito_director',
    'honest_critic', 'reliable_bro', 'lightning'
);

UPDATE badges
SET name = 'Исследователь', description = 'Добавить 5 объявлений в избранное'
WHERE id = 'explorer';

UPDATE badges
SET name = 'Ничего не скроешь', description = 'Улучшить 3 своих объявления'
WHERE id = 'nothing_hidden';

UPDATE badges
SET name = 'Зелёная планета', description = 'Продать 5 объявлений'
WHERE id = 'green_planet';

UPDATE badges
SET name = 'Заряженный стрик', description = 'Заходить в приложение 14 дней подряд'
WHERE id = 'charged_streak';

UPDATE badges
SET name = 'Друг Енота', description = 'Достичь 10 уровня'
WHERE id = 'raccoon_friend';

INSERT INTO badges (id, name, description, icon_url) VALUES
    ('curious_eye',   'Любопытный глаз', 'Посмотреть 10 объявлений',   '/assets/badges/explorer.svg'),
    ('first_listing', 'Первое объявление', 'Опубликовать 1 объявление', '/assets/badges/complete.svg'),
    ('quest_runner',  'Исполнительный',  'Выполнить 3 задания за день', '/assets/badges/streak.svg')
ON CONFLICT (id) DO UPDATE
SET name = EXCLUDED.name,
    description = EXCLUDED.description,
    icon_url = EXCLUDED.icon_url;

UPDATE rewards
SET title = 'Бейдж «Мастер объявлений»', description = 'Бейдж в профиль'
WHERE id = 'dialogue_master_badge';

-- +goose Down
UPDATE rewards
SET title = 'Бейдж «Мастер диалога»', description = 'Бейдж в профиль'
WHERE id = 'dialogue_master_badge';

DELETE FROM user_badges WHERE badge_id IN ('curious_eye', 'first_listing', 'quest_runner');
DELETE FROM badges WHERE id IN ('curious_eye', 'first_listing', 'quest_runner');

UPDATE badges
SET description = 'Добавить в избранное товары из 5 разных категорий'
WHERE id = 'explorer';

UPDATE badges
SET description = 'Заполнить объявление на 100%'
WHERE id = 'nothing_hidden';

UPDATE badges
SET description = 'Продать 5 б/у вещей'
WHERE id = 'green_planet';

INSERT INTO badges (id, name, description, icon_url) VALUES
    ('quick_finger',    'Быстрый палец',       'Написать продавцу в течение 60 минут после публикации',  '/assets/badges/quick-finger.svg'),
    ('discount_hunter', 'Охотник за скидками', 'Купить товар из раздела спецпредложений',                '/assets/badges/discount-hunter.svg'),
    ('avito_director',  'Режиссёр Авито',      'Опубликовать 5 объявлений с видео длиннее 30 секунд',    '/assets/badges/director.svg'),
    ('honest_critic',   'Честный критик',      'Оставить 10 подробных отзывов с фотографиями',           '/assets/badges/critic.svg'),
    ('reliable_bro',    'Надёжный бро',        'Оставить 10 развёрнутых отзывов с видео',                '/assets/badges/bro.svg'),
    ('lightning',       'Молния',              'Ответить на 10 сообщений быстрее чем за 15 минут',       '/assets/badges/lightning.svg')
ON CONFLICT (id) DO NOTHING;
