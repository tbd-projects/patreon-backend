ALTER TABLE awards
    ADD CONSTRAINT creator_id_price_key UNIQUE (creator_id, price);