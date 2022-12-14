# go-sstable

[![codecov](https://codecov.io/gh/hjcian/go-sstable/branch/main/graph/badge.svg?token=W80KS02GV8)](https://codecov.io/gh/hjcian/go-sstable)

# Goal

This is my learning journey of trying to implement the philosophy of [SSTable](https://www.igvita.com/2012/02/06/sstable-and-log-structured-storage-leveldb/) using [Go](https://go.dev/).

The beginning of this journey is refer to the chapter of "Data Structures That Power Your Database (資料結構：資料庫動力之源)" in Designing Data-Intensive Applications. So, let's get started.


# References
- [資料密集型應用系統設計 (translate to Mandarin Chinese)](https://www.tenlong.com.tw/products/9789865028350)
- [Designing Data-Intensive Applications](https://www.oreilly.com/library/view/designing-data-intensive-applications/9781491903063/)
- [SSTable and Log Structured Storage: LevelDB](https://www.igvita.com/2012/02/06/sstable-and-log-structured-storage-leveldb/)
- [Starting from Zero: Build an LSM Database with 500 Lines of Code](https://www.alibabacloud.com/blog/starting-from-zero-build-an-lsm-database-with-500-lines-of-code_598114)


# Development Notes
- [Q: How to maintain the sparse index in a LSM-tree?](https://stackoverflow.com/questions/69103575/how-to-maintain-the-sparse-index-in-a-lsm-tree)
  - A: A typical approach is to have a separate index per segment file, and this index is re-generated during compaction/merging of segment files. *answered by [Martin Kleppmann](https://stackoverflow.com/a/69103900), the author of book*

# TODO
- implement a self-balanced tree structure for memtable?
  - [AVL tree](https://josephjsf2.github.io/data/structure/and/algorithm/2019/06/22/avl-tree.html)
  - [Red Black Tree](https://josephjsf2.github.io/data/structure/and/algorithm/2020/04/28/red-black-tree-part-1.html)

