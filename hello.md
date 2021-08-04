# Demo

 服务示例

 服务级的注释都会同步到 markdown 文件

- [/api/demo.Demo/Echo1](#apidemodemoecho1)
- [/api/demo.Demo/Echo2](#apidemodemoecho2)

## /api/demo.Demo/Echo1

 接口示例

 接口注释会同步到 markdown 文件的对应位置


### Request
```javascript
{
    // 布尔值示例
    a: false, // type<bool>
    // 三十二位整数示例
    b: 0, // type<int32>
    // 六十四位整数示例
    c: "0", // type<int64>, 使用 string 保存
    // 浮点数示例
    d: 0.0, // type<double>
    // 字符示例
    e: "", // type<string>
    // 字节序列示例
    f: "", // type<bytes>, 使用 base64 string 保存
    // 对象示例
    g: {
        // 字符列表示例
        a: [
            ""
        ], // list<string>
        // 字典示例
        b: {
            "0": ""
        }, // map<int32,string>
        // 递归定义示例
        c: {
            a: {}, // type<Bar>, 被引用的对象不再展开
        }, // type<Baz>
        // 枚举示例
        d: "Unknown", // enum<Unknown,Male,Female>
        // 对象列表示例
        e: [{
            name: "", // type<string>
            age: 0, // type<int32>
        }], // list<Person>
    }, // type<Bar>
    // 导入对象示例
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

 接口示例二


### Request
```javascript
{}
```

### Reply
```javascript
{
    // 布尔值示例
    a: false, // type<bool>
    // 三十二位整数示例
    b: 0, // type<int32>
    // 六十四位整数示例
    c: "0", // type<int64>, 使用 string 保存
    // 浮点数示例
    d: 0.0, // type<double>
    // 字符示例
    e: "", // type<string>
    // 字节序列示例
    f: "", // type<bytes>, 使用 base64 string 保存
    // 对象示例
    g: {
        // 字符列表示例
        a: [
            ""
        ], // list<string>
        // 字典示例
        b: {
            "0": ""
        }, // map<int32,string>
        // 递归定义示例
        c: {
            a: {}, // type<Bar>, 被引用的对象不再展开
        }, // type<Baz>
        // 枚举示例
        d: "Unknown", // enum<Unknown,Male,Female>
        // 对象列表示例
        e: [{
            name: "", // type<string>
            age: 0, // type<int32>
        }], // list<Person>
    }, // type<Bar>
    // 导入对象示例
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

