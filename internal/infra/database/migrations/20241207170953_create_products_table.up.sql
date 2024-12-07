CREATE TABLE products (
  id         uuid PRIMARY KEY,
  name       varchar(50) NOT NULL,
  value      integer NOT NULL,
  unityType  char(2) NOT NULL,
  isDeleted  boolean DEFAULT false NOT NULL,
  createdAt  timestamp DEFAULT now() NOT NULL,
  updatedAt  timestamp
);
