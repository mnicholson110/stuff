#include "bufferpool.h"

#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

BufferPool *create_buffer_pool()
{
    BufferPool *bp = malloc(sizeof(BufferPool));
    if (!bp)
    {
        perror("Failed to allocate memory for Buffer Pool");
        return NULL;
    }

    memset(bp->in_use, 0, sizeof(bp->in_use));
    return bp;
}

Page *get_page_from_buffer_pool(BufferPool *bp, DiskManager *dm,
                                uint32_t page_id)
{
    // Check if the page is already in the buffer pool
    for (int i = 0; i < NUM_BUFFERS; ++i)
    {
        if (bp->in_use[i] && bp->page_dir[i] == page_id)
        {
            return &bp->page_buffers[i];
        }
    }

    // Page is not in the buffer pool, so load it directly into the buffer pool
    // slot
    for (int i = 0; i < NUM_BUFFERS; ++i)
    {
        if (!bp->in_use[i])
        {
            read_page(dm, page_id, &bp->page_buffers[i]);
            bp->page_dir[i] = page_id;
            bp->in_use[i] = 1;
            return &bp->page_buffers[i];
        }
    }

    fprintf(stderr, "Error: Buffer pool is full. Cannot load page %u.\n",
            page_id);
    return NULL;
}

// Set (or replace) a page in the buffer pool

void remove_page_from_buffer_pool(BufferPool *bp, uint32_t page_id)
{
    for (int i = 0; i < NUM_BUFFERS; ++i)
    {
        if (bp->in_use[i] && bp->page_dir[i] == page_id)
        {
            // Clear the buffer
            bp->in_use[i] = 0;
            bp->page_dir[i] = (uint32_t)-1; // Mark as invalid
            return;
        }
    }

    fprintf(stderr, "Error: Page %u is not in the buffer pool.\n", page_id);
}

void destroy_buffer_pool(BufferPool *bp)
{
    if (bp)
    {
        free(bp);
    }
}
