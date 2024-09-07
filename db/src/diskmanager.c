#include "diskmanager.h"

#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

DiskManager *create_disk_manager(const char *filepath)
{
    DiskManager *dm = malloc(sizeof(DiskManager));
    if (!dm)
    {
        perror("Failed to allocate memory for DiskManager");
        return NULL;
    }
    dm->filepath = filepath;
    dm->db_fd = -1;
    return dm;
}

void destroy_disk_manager(DiskManager *dm)
{
    if (dm)
    {
        if (dm->db_fd != -1)
        {
            close_dm(dm);
        }
        free(dm);
    }
}

void open_dm(DiskManager *dm)
{
    if ((dm->db_fd = open(dm->filepath, O_RDWR | O_CREAT | O_APPEND, 0644)) ==
        -1)
    {
        perror("Error opening file");
    }
}

void write_page(DiskManager *dm, Page *page)
{
    uint32_t offset = page->page_id * PAGE_SIZE;
    if (pwrite(dm->db_fd, page, PAGE_SIZE, offset) == -1)
    {
        perror("Error writing page");
    }
    if (fsync(dm->db_fd) == -1)
    {
        perror("Error syncing file");
    }
}

void read_page(DiskManager *dm, uint32_t page_id, Page *page)
{
    uint32_t offset = page_id * PAGE_SIZE;
    if (pread(dm->db_fd, (void *)page, PAGE_SIZE, offset) == -1)
    {
        perror("Error reading page");
    }
}

void close_dm(DiskManager *dm)
{
    if (close(dm->db_fd) == -1)
    {
        perror("Error closing file");
    }
    dm->db_fd = -1;
}
