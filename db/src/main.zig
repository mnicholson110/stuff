const std = @import("std");

const diskmanager = struct {
    // TODO: add error handling throughout
    filepath: []const u8,
    db_file: std.fs.File,

    fn open(self: *diskmanager) !void {
        self.db_file = try std.fs.cwd().createFile(self.filepath, .{ .read = true });
    }

    fn write_page(self: *diskmanager, page_id: u32, page_data: []const u8) !void {
        const offset: u64 = page_id * 4096;
        _ = try self.db_file.pwrite(page_data, offset);
        try self.db_file.sync();
    }

    fn read_page(self: *diskmanager, page_id: u32, page_data: []u8) !void {
        const offset: u64 = page_id * 4096;
        _ = try self.db_file.pread(page_data, offset);
    }

    fn close(self: *diskmanager) void {
        self.db_file.close();
    }
};

const TestRecord = packed struct {
    a: u64,
    b: u64,
    c: u64,
    d: u64,
    e: u64,
    f: u64,
    g: u64,
    h: u64,
};

pub fn main() !void {
    const one = TestRecord{
        .a = 1,
        .b = 1,
        .c = 1,
        .d = 1,
        .e = 1,
        .f = 1,
        .g = 1,
        .h = 1,
    };
    const two = TestRecord{
        .a = 2,
        .b = 2,
        .c = 2,
        .d = 2,
        .e = 2,
        .f = 2,
        .g = 2,
        .h = 2,
    };
    const three = TestRecord{
        .a = 3,
        .b = 3,
        .c = 3,
        .d = 3,
        .e = 3,
        .f = 3,
        .g = 3,
        .h = 3,
    };

    var dm = diskmanager{ .filepath = "db.txt", .db_file = undefined };
    defer dm.close();

    try dm.open();

    const buf = try std.heap.page_allocator.alloc(u8, 4096);
    defer std.heap.page_allocator.free(buf);

    try dm.write_page(1, std.mem.asBytes(&one));
    try dm.write_page(2, std.mem.asBytes(&two));
    try dm.write_page(0, std.mem.asBytes(&three));

    try dm.read_page(1, buf);
    const one_from_buf = std.mem.bytesToValue(TestRecord, buf);
    std.debug.print("{any}\n", .{one_from_buf});

    try dm.read_page(2, buf);
    const two_from_buf = std.mem.bytesToValue(TestRecord, buf);
    std.debug.print("{any}\n", .{two_from_buf});

    try dm.read_page(0, buf);
    const three_from_buf = std.mem.bytesToValue(TestRecord, buf);
    std.debug.print("{any}\n", .{three_from_buf});
}
