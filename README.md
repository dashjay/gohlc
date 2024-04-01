# Go HLC(Hybrid Logic Clock)

## Intro

To implement a mvcc manager for object storage, we usually design a version allocator usually a number that increases
monotonically, but sometimes machine crashed and time drift back to the past, and version(clock) allocated re-allocated
by mistake, so HLC is for solving this problem.

I didn't find any repo fit for my request, so I decided to designed/implemented my own hlc module.

## Detail

I use a time(unix_nano) as part of the physical clock and allocate version from this time(clock) like this:

```text
// pseudo-code for allocate one version
 
v := allocatedVersion
desire := v+1
cas(&allocatedVersion, v, desire)
```

In addition to the logic above , we have another thread to save version to persistent storage, it updated
`lastSavedTimestamp` after save and signal all clients who wait for versions.

if var `desire` catch up the `lastSavedTimestamp`, the allocate thread wait quietly, after saving `hybrid logic clock`
to persistent storage, the process will be signaled.

Allocator process need not wait for the persistence of timestamp under normal conditions, there is a timed task to do
this periodically, they even not feel the existence of persistence process.

## Architecture

- **persistence**: developers can define a persistence.Interface by themselves, save the version to where the like: eg
  boltdb, disk, etcd, mysql...(I have implemented a disk one for you to use)
- **service**: it implemented a hlc service for user to call by grpc.


## Use

It is easy to run service as a gprc server, I'll provide the full code sample later. 


