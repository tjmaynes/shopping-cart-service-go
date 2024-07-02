-- migrate:up
CREATE TABLE item (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR (255) NOT NULL,
  price BIGINT NOT NULL,
  manufacturer VARCHAR (255) NOT NULL
);

-- migrate:down
DROP TABLE IF EXISTS item;
