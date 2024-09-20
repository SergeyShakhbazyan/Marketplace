CREATE KEYSPACE "marketplace_keyspace" WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

CREATE TABLE marketplace_keyspace.userData(
                                              id UUID,
                                              firstName TEXT,
                                              lastName TEXT,
                                              avatar TEXT,
                                              phoneNumber TEXT,
                                              email TEXT,
                                              AccountType TEXT,
                                              Subscription boolean,
                                              password TEXT,
                                              createdAt TIMESTAMP,
                                              rating DECIMAL,
                                              PRIMARY KEY (id)
);

CREATE TABLE marketplace_keyspace.product (
                                              product_id UUID,
                                              owner_id UUID,
                                              title TEXT,
                                              image LIST<TEXT>,
                                              description TEXT,
                                              price INT,
                                              brandName TEXT,
                                              category_id UUID,
                                              subcategory_id UUID,
                                              created_at TIMESTAMP,
                                              keywords SET<TEXT>,
                                              PRIMARY KEY ((category_id, subcategory_id), created_at, product_id)
)WITH CLUSTERING ORDER BY (created_at desc);

CREATE TABLE marketplace_keyspace.product_views (
                                                    product_id UUID,
                                                    views COUNTER,
                                                    PRIMARY KEY (product_id)
);


CREATE INDEX ON marketplace_keyspace.product (category_id);

CREATE INDEX ON marketplace_keyspace.product (owner_id);

CREATE INDEX ON marketplace_keyspace.product (keywords);

CREATE INDEX ON marketplace_keyspace.product (product_id);


CREATE TABLE marketplace_keyspace.category(
                                              id UUID,
                                              name TEXT,
                                              image TEXT,
                                              PRIMARY KEY (id, name)
);


CREATE TABLE marketplace_keyspace.subcategories(
                                                   category_id UUID,
                                                   subcategory_id UUID,
                                                   name TEXT,
                                                   groupID UUID,
                                                   PRIMARY KEY (subcategory_id, category_id)
);

CREATE TABLE marketplace_keyspace.subcategoryGroups(
                                                       groupID UUID,
                                                       groupName TEXT,
                                                       subcategory_ids FROZEN<LIST<UUID>>,
                                                       categoryID UUID,
                                                       PRIMARY KEY (categoryID, groupName, groupID)
);


CREATE TABLE marketplace_keyspace.brands (
                                             subcategory_id UUID,
                                             id UUID,
                                             name TEXT,
                                             PRIMARY KEY (subcategory_id, id)
);

CREATE MATERIALIZED VIEW marketplace_keyspace.brands_by_name AS
SELECT subcategory_id, id, name
FROM marketplace_keyspace.brands
WHERE name IS NOT NULL AND
    id IS NOT NULL
PRIMARY KEY (name, subcategory_id, id);

CREATE TABLE marketplace_keyspace.models (
                                             brand_id UUID,
                                             id UUID,
                                             name TEXT,
                                             PRIMARY KEY (brand_id, id)
);

CREATE TABLE marketplace_keyspace.model_parameters(
                                                      modelID UUID,
                                                      parameterName TEXT,
                                                      parameterValue LIST<TEXT>,
                                                      PRIMARY KEY (modelID, parameterName)
);

SELECT name, category_id FROM marketplace_keyspace.subcategories WHERE subcategory_id = c39ba478-5e7a-11ef-b0e2-2a0b725efeb0;


select category_id, subcategory_id, name from marketplace_keyspace.subcategories where name = 'Women''s Clothing';



CREATE MATERIALIZED VIEW marketplace_keyspace.subcategories_by_name_category_id AS
SELECT category_id, subcategory_id, name
FROM marketplace_keyspace.subcategories
WHERE name IS NOT NULL AND category_id IS NOT NULL AND subcategory_id IS NOT NULL
PRIMARY KEY ((name, category_id), subcategory_id);



CREATE MATERIALIZED VIEW marketplace_keyspace.products_by_category_and_date AS
SELECT product_id, created_at
FROM marketplace_keyspace.product
WHERE category_id IS NOT NULL
  AND product_id IS NOT NULL
  AND created_at IS NOT NULL
PRIMARY KEY (category_id, product_id, created_at);

SELECT * FROM marketplace_keyspace.product WHERE category_id = 0027d084-646f-11ef-85a7-38c9863c85bd;
SELECT product_id, title, image, price FROM marketplace_keyspace.product WHERE category_id = e4c15ef6-504a-11ef-8f69-2a0b725efeae;
SELECT product_id, title, image, price, created_at FROM marketplace_keyspace.product WHERE owner_id = 6f3a8698-5a37-11ef-8605-2a0b725efeb0;
SELECT product_id, title, image, price, created_at FROM marketplace_keyspace.product WHERE keywords CONTAINS 'iphone 13';

CREATE TABLE marketplace_keyspace.product_filters (
                                                      category_id UUID,
                                                      sub_category_id UUID,
                                                      filter_name TEXT,
                                                      filter_value TEXT,
                                                      product_id UUID,
                                                      PRIMARY KEY ((category_id, sub_category_id, filter_name, filter_value), product_id)
);

CREATE MATERIALIZED VIEW marketplace_keyspace.product_filters_by_id AS
SELECT category_id, sub_category_id, filter_name, filter_value, product_id
FROM marketplace_keyspace.product_filters
WHERE product_id IS NOT NULL
  AND category_id IS NOT NULL
  AND sub_category_id IS NOT NULL
  AND filter_name IS NOT NULL
  AND filter_value IS NOT NULL
PRIMARY KEY (product_id, category_id, sub_category_id, filter_name, filter_value);


CREATE INDEX ON marketplace_keyspace.product_filters (product_id);

CREATE MATERIALIZED VIEW marketplace_keyspace.product_by_id AS
SELECT product_id, owner_id, title, image, description, price,brandName,
       category_id, subcategory_id, created_at, keywords
FROM marketplace_keyspace.product
WHERE product_id IS NOT NULL
  AND  category_id IS NOT NULL
  AND subcategory_id IS NOT NULL
  AND created_at IS NOT NULL
PRIMARY KEY (product_id, category_id, subcategory_id, created_at);

UPDATE marketplace_keyspace.product_views SET views = views + 1 WHERE product_id = 08d41ac4-661f-11ef-b64e-2a0b725efeb0;

DELETE FROM marketplace_keyspace.product WHERE product_id = 89d0d0f6-696a-11ef-bb7a-2a0b725efeb0 AND  created_at = '2024-09-02 20:32:57.138' AND category_id = 0027d084-646f-11ef-85a7-38c9863c85bd AND subcategory_id  = a198dd16-64ef-11ef-8f7e-2a0b725efeb0