#pragma once

#include "common.h"
#include <stdint.h>

// Constants for B+ Tree node properties
#define MAX_KEYS 3                  // Maximum number of keys in a node (depends on degree)
#define MAX_CHILDREN (MAX_KEYS + 1) // Maximum number of children for an internal node

typedef struct Value Value;
struct Value
{
    uint32_t page_id;
    uint32_t offset;
    uint32_t size;
};

// Structure for a leaf node in the B+ Tree
typedef struct BTreeLeafNode BTreeLeafNode;
struct BTreeLeafNode
{
    uint32_t keys[MAX_KEYS]; // Array of keys in the leaf node
    Value values[MAX_KEYS];  // Array of values
    BTreeLeafNode *next;     // Pointer to right sibling
    BTreeLeafNode *prev;     // Pointer to left sibling
    uint32_t numKeys;        // Number of keys in the leaf node
};

// Structure for an internal node in the B+ Tree
typedef struct BTreeInternalNode BTreeInternalNode;
struct BTreeInternalNode
{
    uint32_t keys[MAX_KEYS];                 // Array of keys in the internal node
    union BTreeNode *children[MAX_CHILDREN]; // Array of child pointers (to either internal or leaf nodes)
    uint32_t numKeys;                        // Number of keys in the internal node
};

typedef union BTreeNode BTreeNode;
union BTreeNode {
    BTreeLeafNode leaf;
    BTreeInternalNode internal;
};

// Structure for the B+ Tree
typedef struct BTree BTree;
struct BTree
{
    BTreeNode *root; // Pointer to the root node (either leaf or internal)
    uint32_t height; // Height of the tree (optional)
};
