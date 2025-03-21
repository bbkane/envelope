# migration_v2

## Description

<details>
<summary><strong>Table Definition</strong></summary>

```sql
CREATE TABLE migration_v2 (
	migration_v2_id INTEGER PRIMARY KEY,
	file_name TEXT NOT NULL,
	migrate_time TEXT NOT NULL,
	UNIQUE(file_name)
) STRICT
```

</details>

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| migration_v2_id | INTEGER |  | true |  |  |  |
| file_name | TEXT |  | false |  |  |  |
| migrate_time | TEXT |  | false |  |  |  |

## Constraints

| Name | Type | Definition |
| ---- | ---- | ---------- |
| migration_v2_id | PRIMARY KEY | PRIMARY KEY (migration_v2_id) |
| sqlite_autoindex_migration_v2_1 | UNIQUE | UNIQUE (file_name) |

## Indexes

| Name | Definition |
| ---- | ---------- |
| sqlite_autoindex_migration_v2_1 | UNIQUE (file_name) |

## Relations

![er](migration_v2.svg)

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
