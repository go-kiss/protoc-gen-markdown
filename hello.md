# Demo

 Service demo.

 All leading comments will be copied to markdown.

- [/api/demo.Demo/Echo1](#apidemodemoecho1)
- [/api/demo.Demo/Echo2](#apidemodemoecho2)

## /api/demo.Demo/Echo1

 Rpc demo

 All leading comments will be copied to markdown.


### Request
```javascript
{
    // boolean value demo
    a: false, // type<bool>
    // 32 bit int value demo
    b: 0, // type<int32>
    // 64 bit int value demo
    c: "0", // type<int64>, stored as string
    // float value demo
    d: 0.0, // type<double>
    // string value demo
    e: "", // type<string>
    // bytes value demo
    f: "", // type<bytes>, stored as base64 string
    // message value demo
    g: {
        // string list value demo
        a: [
            ""
        ], // list<string>
        // map value demo
        b: {
            "0": ""
        }, // map<int32,string>
        // self reference value demo
        c: {
            a: {}, // type<Bar>, self-referenced message will be displayed as {}
        }, // type<Baz>
        // enum value demo
        d: "Unknown", // enum<Unknown,Male,Female>
        // message list value demo
        e: [{
            name: "", // type<string>
            age: 0, // type<int32>
        }], // list<Person>
    }, // type<Bar>
    // imported message value demo
    h: {
        // Represents seconds of UTC time since Unix epoch
        // 1970-01-01T00:00:00Z. Must be from 0001-01-01T00:00:00Z to
        // 9999-12-31T23:59:59Z inclusive.
        seconds: "0", // type<int64>
        // Non-negative fractions of a second at nanosecond resolution. Negative
        // second values with fractions must still have non-negative nanos values
        // that count forward in time. Must be from 0 to 999,999,999
        // inclusive.
        nanos: 0, // type<int32>
    }, // type<Timestamp>
}
```

### Reply
```javascript
{}
```
## /api/demo.Demo/Echo2

 Another rpc demo


### Request
```javascript
{}
```

### Reply
```javascript
{
    // boolean value demo
    a: false, // type<bool>
    // 32 bit int value demo
    b: 0, // type<int32>
    // 64 bit int value demo
    c: "0", // type<int64>, stored as string
    // float value demo
    d: 0.0, // type<double>
    // string value demo
    e: "", // type<string>
    // bytes value demo
    f: "", // type<bytes>, stored as base64 string
    // message value demo
    g: {
        // string list value demo
        a: [
            ""
        ], // list<string>
        // map value demo
        b: {
            "0": ""
        }, // map<int32,string>
        // self reference value demo
        c: {
            a: {}, // type<Bar>, self-referenced message will be displayed as {}
        }, // type<Baz>
        // enum value demo
        d: "Unknown", // enum<Unknown,Male,Female>
        // message list value demo
        e: [{
            name: "", // type<string>
            age: 0, // type<int32>
        }], // list<Person>
    }, // type<Bar>
    // imported message value demo
    h: {
        // Represents seconds of UTC time since Unix epoch
        // 1970-01-01T00:00:00Z. Must be from 0001-01-01T00:00:00Z to
        // 9999-12-31T23:59:59Z inclusive.
        seconds: "0", // type<int64>
        // Non-negative fractions of a second at nanosecond resolution. Negative
        // second values with fractions must still have non-negative nanos values
        // that count forward in time. Must be from 0 to 999,999,999
        // inclusive.
        nanos: 0, // type<int32>
    }, // type<Timestamp>
}
```

