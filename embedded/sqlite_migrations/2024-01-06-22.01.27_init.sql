CREATE TABLE "env" (
    "id" INTEGER PRIMARY KEY NOT NULL,
    "name" TEXT NOT NULL,
    "comment" TEXT,
    "create_time" TEXT NOT NULL,
    "update_time" TEXT NOT NULL,
    UNIQUE(name)
) STRICT;
