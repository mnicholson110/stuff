#pragma once
#include <stddef.h>
#include <stdint.h>

#include "page.h"

typedef struct DiskManager DiskManager;
struct DiskManager
{
    const char *filepath;
    int db_fd;
};

DiskManager *create_disk_manager(const char *filepath);
void destroy_disk_manager(DiskManager *dm);

void open_dm(DiskManager *dm);
void write_page(DiskManager *dm, Page *page);
void read_page(DiskManager *dm, uint32_t page_id, Page *page);
void close_dm(DiskManager *dm);
