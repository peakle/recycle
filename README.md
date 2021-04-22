[POST] `/v1/order/create`

Required Headers:
Content-Type: application/x-www-form-urlencoded

Request example:

```
{
    UserId  int     "user_id"
    Address string  "address"
    MaxSize int     "max_size"
    EventAt string  "event_at"
}
```

[POST]  `/v1/order/subscribe` - subscribe user to event Required Headers:

Required Headers:
Content-Type: application/x-www-form-urlencoded

Request example:

```
{
    UserId      int     "user_id"
    OrderId     int     "order_id"
}
```

[GET] `/v1/order/list` - return list of active events

[GET] `/v1/order/info` - return detailed info about event
```
{
    OrderId     int     "order_id"
}
```
