ezorm-gen-databend
---

ezorm-gen-databend is a plugin for [ezorm](github.com/ezbuy/ezorm)
to integrate with [databend](github.com/datafuselabs/databend).

## Features

* Generate databend table schema from ezorm yaml.

## Usage

> See [ezorm plugin doc](https://github.com/ezbuy/ezorm/blob/main/doc/plugin.mdd)

```shell
ezorm gen -i ./tests/ -o ./tests/ --goPackage databend --plugin databend --plugin-only
```

and then the result SQL will locate in `./tests/test_blogs_create_table.sql` and `./tests/test_users_create_table.sql`.
