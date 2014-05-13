# Overview

xuanwu is a GO ORM using code generation approach.

Supporting MongoDB & Mysql.

Definition file is using thrift.

Parser is using https://github.com/wickman/ptsd

`pip install ply`

Template is using cheetah

	sudo su
	export ARCHFLAGS=-Wno-error=unused-command-line-argument-hard-error-in-future
	pip install cheetah

# Usage

	python xuanwu.py thrift_file_path output_folder_path

# Types

Currently only supports:

* bool: A boolean value (true or false), one byte
* i32: A 32-bit signed integer
* string: Encoding agnostic text or binary string

Will support when needed:

* byte: A signed byte
* i16: A 16-bit signed integer
* i64: A 64-bit signed integer
* double: A 64-bit floating point number
* Containers: list/set/map
* Enums

## Annotations

Field could be annotated, for example:

```java
struct User {
	1: string ID
	2: string UserName (label = "User Name", index = "username")
	3: string Password
	4: string Name     (search = "simple", search = "full-10")
	5: string Email
	6: string Intro    (label = "Self Introduction", search = "full-2")
	7: string Picture  (index = "adminWithPictureIndex")
	8: string Remark
	9: bool IsAdmin    (index = "adminWithPictureIndex")
	10: string UserGroupID
}
```

In above example, xuanwu will use `label` defined when generating HTML widgets.

Two indexes:

* UserName
* Union Index Picture + IsAdmin

Two types of search:

* simple: only seach for name
* full: seach for name (weight 10) & intro (weight 2)

# Convention

* All struct must has the field `1: string ID`
  * xuanwu will map this field to mongo's `_id` specially
* First character of field name must be **upper case**
* Field name ends with ID is the foreign key to another struct
  * `HospitalID` => foreign key to `Hospital` struct


# TBD

Generate crud template?
  * How to update?

Use reflect to support crud?

How to embed groupcache?
