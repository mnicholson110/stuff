const std = @import("std");

const PAGE_SIZE: u16 = 8192;

const diskpage = extern struct {
    // LAYOUT
    // Header should contain:
    //    PageID: u16?
    //    Tuple Size in bytes: u32?
    //    Offset to first tuple: u16?
    //    MAX DATA SIZE = PAGE_SIZE - 14
    page_id: u16 align(1),
    tuple_size: u32 align(1),
    num_tuples: u32 align(1),
    free_tuple_offset: u32 align(1),
    data: [PAGE_SIZE - 14]u8 align(1),

    // need checks on all kinds of shit here
    pub fn add_tuple(self: *diskpage, tuple: []const u8) void {
        var offset = self.free_tuple_offset;
        if (offset == 0) { // means no data, insert at end of page
            self.free_tuple_offset = PAGE_SIZE - self.tuple_size - 15;
            offset = self.free_tuple_offset;
        }
        @memcpy(self.data[offset .. offset + self.tuple_size], tuple);
        self.free_tuple_offset -= self.tuple_size;
        self.num_tuples += 1;
    }

    pub fn get_tuple(self: *diskpage, tuple_num: u8) []u8 {
        //std.debug.print("{any}: PAGE_SIZE\n", .{PAGE_SIZE});
        //std.debug.print("{any}: tuple_num\n", .{tuple_num});
        //std.debug.print("{any}: tuple_size\n", .{self.tuple_size});
        const offset: u32 = PAGE_SIZE - ((tuple_num + 1) * self.tuple_size) - 15;
        return self.data[offset .. offset + self.tuple_size];
    }
};

const diskmanager = struct {
    // TODO: add error handling throughout
    filepath: []const u8,
    db_file: std.fs.File,

    fn open(self: *diskmanager) !void {
        self.db_file = try std.fs.cwd().createFile(self.filepath, .{ .read = true });
    }

    fn write_page(self: *diskmanager, page: *diskpage) !void {
        const offset: u64 = page.page_id * PAGE_SIZE;
        _ = try self.db_file.pwrite(std.mem.asBytes(page), offset);
        try self.db_file.sync();
    }

    fn read_page(self: *diskmanager, page_id: u16, buffer: []u8) !diskpage {
        const offset: u64 = page_id * PAGE_SIZE;
        _ = try self.db_file.pread(buffer, offset);
        //doesnt work: const page = std.mem.bytesToValue(diskpage, &buffer);
        var page: diskpage = undefined;
        page.page_id = std.mem.readInt(u16, buffer[0..2], std.builtin.Endian.little);
        page.tuple_size = std.mem.readInt(u32, buffer[2..6], std.builtin.Endian.little);
        page.num_tuples = std.mem.readInt(u32, buffer[6..10], std.builtin.Endian.little);
        page.free_tuple_offset = std.mem.readInt(u32, buffer[10..14], std.builtin.Endian.little);
        @memcpy(&page.data, buffer[14..PAGE_SIZE]);
        return page;
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

    // TODO: add buffer manager to handle allocations + pages
    const buf = try std.heap.page_allocator.alloc(u8, PAGE_SIZE);
    defer std.heap.page_allocator.free(buf);

    var page1: diskpage = .{ .page_id = 0, .tuple_size = 64, .num_tuples = 0, .free_tuple_offset = 0, .data = undefined };
    page1.add_tuple(std.mem.asBytes(&one));
    page1.add_tuple(std.mem.asBytes(&two));
    page1.add_tuple(std.mem.asBytes(&three));
    try dm.write_page(&page1);

    //try dm.read_page(0, buf);
    var page1_from_disk = try dm.read_page(0, buf);
    //std.debug.print("page_id: {any}\n", .{std.mem.readInt(u16, buf[0..2], std.builtin.Endian.little)});
    //std.debug.print("tuple_size: {any}\n", .{std.mem.readInt(u32, buf[2..6], std.builtin.Endian.little)});

    const one_from_disk = std.mem.bytesToValue(TestRecord, page1_from_disk.get_tuple(0));
    std.debug.print("{any}\n", .{one_from_disk});

    const two_from_disk = std.mem.bytesToValue(TestRecord, page1_from_disk.get_tuple(1));
    std.debug.print("{any}\n", .{two_from_disk});

    const three_from_disk = std.mem.bytesToValue(TestRecord, page1_from_disk.get_tuple(2));
    std.debug.print("{any}\n", .{three_from_disk});
}

test {
    std.debug.print("Size of diskpage: {any}\n", .{@sizeOf(diskpage)});
    try std.testing.expectEqual(PAGE_SIZE, @sizeOf(diskpage));
}
