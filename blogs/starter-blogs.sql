BEGIN TRANSACTION;

INSERT INTO categories(name) VALUES ("Category 1");
INSERT INTO categories(name) VALUES ("Category 2");
INSERT INTO categories(name) VALUES ("Category 3");

COMMIT;

BEGIN TRANSACTION;

INSERT INTO
  blogs(id,title,authors,categories,timestamp,tz_offset)
  VALUES (
    1,
    "Blog 1",
    "|Author 1|",
    "|Category 1|Category 2|",
    1736800572, -6*60*60
  )
;

INSERT INTO
  blogs(id,title,authors,categories,timestamp,tz_offset)
  VALUES (
    2,
    "Blog 2",
    "|Author 2|",
    "|Category 2|Category 3|",
    1736890072, -6*60*60
  )
;

INSERT INTO
  blogs(id,title,authors,categories,timestamp,tz_offset)
  VALUES (
    3,
    "Blog 3",
    "|Author 1|Author 2|Author 3|",
    "|Category 1|Category 3|",
    1736890572, -6*60*60
  )
;

COMMIT;

BEGIN TRANSACTION;

INSERT INTO
  edits(blog_id,description,timestamp,tz_offset,prev_hash,hash)
  VALUES (
    1,
    "",
    1736800572, -6*60*60,
    "TODO", "TODO"
  )
;

INSERT INTO
  edits(blog_id,description,timestamp,tz_offset,prev_hash,hash)
  VALUES (
    2,
    "",
    1736890072, -6*60*60
    "TODO", "TODO"
  )
;

INSERT INTO
  edits(blog_id,description,timestamp,tz_offset,prev_hash,hash)
  VALUES (
    3,
    "",
    1736890572, -6*60*60
    "TODO", "TODO"
  )
;

COMMIT;
