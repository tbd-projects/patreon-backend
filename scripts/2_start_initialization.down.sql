DELETE FROM creator_category WHERE
    name = 'Подкасты' or
    name = 'Художники' or
    name = 'Музыканты' or
    name = 'Писатели и журналисты' or
    name = 'Видеоблогер' or
    name = 'Образование' or
    name = 'Программирование' or
    name = 'Другое';

DELETE FROM posts_type WHERE
        type = 'music' or
        type = 'video' or
        type = 'files' or
        type = 'text' or
        type = 'image';


