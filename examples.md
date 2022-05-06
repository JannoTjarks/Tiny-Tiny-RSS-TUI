# Examples for the usage of the TTRSS API

## Login to the TTRSS server
``` sh
curl -d '{"op":"login","user":"xxx","password":"xxx"}' https://xxx/api/
```

## Get the current API version
``` sh
curl -d '{"sid":"xxx","op":"getApiLevel"}' https://xxx/api/
```

## Get the categories
``` sh
curl -d '{"sid":"xxx","op":"getCategories"}' https://xxx/api/
```

## Get the feeds
``` sh
curl -d '{"sid":"xxx","op":"getFeeds","cat_id":"4"}' https://xxx/api/
```

## Get the headlines
``` sh
curl -d '{"sid":"","op":"getHeadlines","feed_id":"24"}' https://xxx/api/
```
