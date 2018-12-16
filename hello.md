- [/api/routeguide.RouteGuide/GetFeature](#apirouteguiderouteguidegetfeature)
- [/api/routeguide.RouteGuide/ListFeatures](#apirouteguiderouteguidelistfeatures)
- [/api/routeguide.RouteGuide/RecordRoute](#apirouteguiderouteguiderecordroute)
- [/api/routeguide.RouteGuide/RouteChat](#apirouteguiderouteguideroutechat)

## /api/routeguide.RouteGuide/GetFeature

A simple RPC.

Obtains the feature at a given position.

A feature with an empty name is returned if there's no feature at the given
position.

### Method

POST

### Request
```javascript
{
    latitude: 0, // type:<int>
    longitude: 0, // type:<int>
}
```

### Reply
```javascript
{
    // The name of the feature.
    name: "", // type:<string>
    // The point where the feature is detected.
    location: {
        latitude: 0, // type:<int>
        longitude: 0, // type:<int>
    },
}
```
## /api/routeguide.RouteGuide/ListFeatures

A server-to-client streaming RPC.

Obtains the Features available within the given Rectangle.  Results are
streamed rather than returned at once (e.g. in a response message with a
repeated field), as the rectangle may cover a large area and contain a
huge number of features.

### Method

POST

### Request
```javascript
{
    // One corner of the rectangle.
    lo: {
        latitude: 0, // type:<int>
        longitude: 0, // type:<int>
    },
    // The other corner of the rectangle.
    hi: {
        latitude: 0, // type:<int>
        longitude: 0, // type:<int>
    },
}
```

### Reply
```javascript
{
    // The name of the feature.
    name: "", // type:<string>
    // The point where the feature is detected.
    location: {
        latitude: 0, // type:<int>
        longitude: 0, // type:<int>
    },
}
```
## /api/routeguide.RouteGuide/RecordRoute

A client-to-server streaming RPC.

Accepts a stream of Points on a route being traversed, returning a
RouteSummary when traversal is completed.

### Method

POST

### Request
```javascript
{
    latitude: 0, // type:<int>
    longitude: 0, // type:<int>
}
```

### Reply
```javascript
{
    // The number of points received.
    point_count: 0, // type:<int>
    // The number of known features passed while traversing the route.
    feature_count: 0, // type:<int>
    // The distance covered in metres.
    distance: 0, // type:<int>
    // The duration of the traversal in seconds.
    elapsed_time: 0, // type:<int>
}
```
## /api/routeguide.RouteGuide/RouteChat

A Bidirectional streaming RPC.

Accepts a stream of RouteNotes sent while a route is being traversed,
while receiving other RouteNotes (e.g. from other users).

### Method

POST

### Request
```javascript
{
    // The location from which the message is sent.
    location: {
        latitude: 0, // type:<int>
        longitude: 0, // type:<int>
    },
    // The message to be sent.
    message: "", // type:<string>
}
```

### Reply
```javascript
{
    // The location from which the message is sent.
    location: {
        latitude: 0, // type:<int>
        longitude: 0, // type:<int>
    },
    // The message to be sent.
    message: "", // type:<string>
}
```
