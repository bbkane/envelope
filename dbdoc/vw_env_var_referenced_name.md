# vw_env_var_referenced_name

## Description

<details>
<summary><strong>Table Definition</strong></summary>

```sql
CREATE VIEW vw_env_var_referenced_name AS
SELECT
    env_var_id,
    env_id,
	(SELECT name FROM env WHERE env_id = env_var.env_id) AS env_name,
    name,
    value,
    comment,
    create_time,
    update_time
FROM env_var
```

</details>

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| env_var_id | INTEGER |  | true |  |  |  |
| env_id | INTEGER |  | true |  |  |  |
| env_name | TEXT |  | true |  |  |  |
| name | TEXT |  | true |  |  |  |
| value | TEXT |  | true |  |  |  |
| comment | TEXT |  | true |  |  |  |
| create_time | TEXT |  | true |  |  |  |
| update_time | TEXT |  | true |  |  |  |

## Referenced Tables

| Name | Columns | Comment | Type |
| ---- | ------- | ------- | ---- |
| [env](env.md) | 5 |  | table |
| [env_var](env_var.md) | 7 |  | table |

## Relations

![er](vw_env_var_referenced_name.svg)

---

> Generated by [tbls](https://github.com/k1LoW/tbls)