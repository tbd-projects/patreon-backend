\connect restapi_dev

INSERT INTO creator_category (name)
    VALUES('Подкасты'),
    ('Музыканты'),
    ('Художники'),
    ('Писатели и журналисты'),
    ('Видеоблогер'),
    ('Образование'),
    ('Программирование'),
    ('Другое');

INSERT INTO posts_type (type)
VALUES('music'),
      ('video'),
      ('image'),
      ('text');

\disconnect

