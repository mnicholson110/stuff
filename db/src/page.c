#include "page.h"

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

Page *create_page(uint32_t page_id, uint32_t tuple_size)
{
    Page *page = malloc(sizeof(Page));
    if (!page)
    {
        perror("Failed to allocate memory for Page");
        return NULL;
    }
    page->page_id = page_id;
    page->tuple_size = tuple_size;
    page->num_tuples = 0;
    page->free_offset = HEADER_SIZE;
    memset(page->data, 0, sizeof(page->data));
    return page;
}

void destroy_page(Page *page)
{
    if (page)
    {
        free(page);
    }
}

void add_tuple(Page *page, const uint8_t *tuple)
{
    uint32_t free_slot = -1;  // default value; TODO: add error checking
    uint32_t *offset_array = get_offset_array(page);

    // using num_tuples isn't correct, as we could delete "middle" tuples,
    // but it's fine for single-threaded work
    for (uint32_t i = 0; i < page->num_tuples + 1; i++)
    {
        if (!offset_array[i])
        {
            free_slot = i;
        }
    }
    // TODO: need error handling here to check valid offset
    // and potentially handle case when tuple is too large for free space
    uint32_t offset =
        PAGE_SIZE - HEADER_SIZE - ((free_slot + 1) * page->tuple_size) - 1;
    memcpy(&page->data[offset], tuple, page->tuple_size);
    offset_array[free_slot] = offset;
    page->free_offset += sizeof(uint32_t);
    page->num_tuples++;
}

uint8_t *get_tuple(Page *page, uint8_t tuple_num)
{
    uint32_t *offset_array = get_offset_array(page);
    // TODO: add bounds check against page->num_tuples
    uint32_t offset = offset_array[tuple_num];
    if (!offset)
    {
        return NULL;
    }
    return &page->data[offset];
}

// this simply returns a pointer to the first byte after the header.
// need to be careful with indexing through this; checks against
// page->num_tuples needed.
uint32_t *get_offset_array(Page *page) { return (uint32_t *)&page->data; }

void remove_tuple(Page *page, uint8_t tuple_num)
{
    get_offset_array(page)[tuple_num] = 0;
};
