# go-client

The `userup` go client package provides a convenience wrapper for the GRPC interface. It supplies common functionality and a flexible query interface to ease integration of the userservice package into your application.

## Installation

Import the package into your project:
```go
import (
    userup "github.com/hillside-labs/userservice-go-sdk/go-client"
)
```

## Usage

The `userup` package provides a `UserService` struct which can be used to interact with the userservice GRPC interface.
Simply supply it with the address of the userservice GRPC server

```go
client, err := userup.NewClient("localhost:9000")
defer client.Close()
```

### Creating a User

```go
user := &userup.User{
    Username: "jdoe2",
    Attributes: map[string]interface{}{
        "user_type": "admin",
        "email":     "jdoe@localhost.com",
        "ranking":   5,
    },
}

ctx := context.Background()
userRet, err := client.AddUser(ctx, user)
```

### Retrieve a User

```go
user, err = client.GetUser(ctx, userId)
```

### Update a User

```go
user, err := client.UpdateUser(ctx, user)
```
Any attributes or traits specified in the user struct will be added/updated during this operation.

### Add an Attribute/Trait to a User

Attribute
```go
err = client.AddAttribute(ctx, user.Id, "user_type", "admin")
```

Trait
```go
err = client.AddTrait(ctx, user.Id, "selected_value", "gibson")
```

### Log an Event

```go
// Initialize the logger
logger := userup.NewLogger(userup.NewLoggerConfig("https://userup.io/sample-client", client))

...

logger.LogEvent(context.Background(), user.Id, "io.userup.user.created", "user", strconv.FormatUint(user.Id, 10), user)
```

## Query Usage

The `Query` struct provides a flexible way to construct and execute queries in the userservice package. This is an experimental portion of the SDK and will likely change as it develops.

### Creating a Query

To create a new query, you can simply initialize a `Query` struct:

```go
query := userservice.NewQuery()
```

### Filtering Users on username
```go
query := userup.Query{
		Filter: map[string]userup.Condition{
			"username": "jdoe2",
		}
}
users, err := client.QueryUsers(ctx, &query)
```


### Filtering with different operators

If no operator is specified, the default operator is `$eq` equality. Other supported operators are 
* `$gt` greater than 
* `$gte` greater than or equal to
* `$lt` less than
* `$lte` less than or equal to
* `$ne` not equal
* `$regex` regular expression
* `$in` in
* `$nin` not in
* `$like` like
* `$ilike` case insensitive like

```go
query := userup.Query{
		Filter: map[string]userup.Condition{
			"id":  map[string]interface{}{"$gte": float64(70)},
		}
}
users, err := client.QueryUsers(ctx, &query)
```

### Filtering users based on attributes value

```go
query := userup.Query{
    Joins: []userup.Join{
        {
            Table: "attributes",
            On:    "users.id = attributes.user_id",
            Filter: map[string]userup.Condition{
                "attributes": {
                    "alias": "dumbledore",
                },
            },
        },
    },
}
users, err := client.QueryUsers(ctx, &query)
```

### Filtering with an OR condition

```go
query := userup.Query{
    Joins: []Join{
        {
            Table: "attributes",
            On:    "users.id = attributes.user_id",
            Filter: map[string]Condition{
                "attributes": {
                    "$or": []Condition{
                        {
                            "2fa": false,
                        },
                        {
                            "vip_level": float64(2),
                        },
                    },
                },
            },
        },
    },
}
```