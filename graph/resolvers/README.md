# GraphQL Resolvers Organization

This directory contains all GraphQL resolver implementations organized by schema using a naming convention.

## File Organization

Since Go doesn't allow splitting a single package across multiple directories, all resolver files are in this directory but use prefixes to indicate which database schema they belong to:

### djong-jukung Schema
- `djong-jukung.mstpricingtiers.resolvers.go` - Pricing tiers resolver
- `djong-jukung.mstpricingtiers.resolvers_test.go` - Tests

### djong-phinisi Schema  
- `djong-phinisi.mstbooking.resolvers.go` - Booking resolver
- `djong-phinisi.mstlead.resolvers.go` - Lead resolver
- `djong-phinisi.mstsalespipeline.resolvers.go` - Sales pipeline resolver
- `djong-phinisi.*.resolvers_test.go` - Tests

### Core Files
- `resolver.go` - Root resolver struct with DB connections
- `generated.go` - Auto-generated GraphQL execution engine (DO NOT EDIT)
- `mock_db_test.go` - Test helpers

## Schema Files

Schema files ARE organized in subdirectories:
- `graph/schema/djong-jukung/` - Schemas for djong_jukung database
- `graph/schema/djong-phinisi/` - Schemas for djong_phinisi database

## Why This Structure?

**Go Limitation**: A single Go package cannot be split across multiple directories. Since all resolvers must be in `package graph`, they must all be in the same directory.

**Solution**: Use naming prefixes (`djong-jukung.`, `djong-phinisi.`) to clearly identify which schema each resolver belongs to while keeping all files in one directory.

## Regenerating Resolvers

When you run `make generate`:
1. gqlgen generates resolver files at the root level
2. The organize script renames them with appropriate prefixes
3. Existing implementations are preserved (thanks to `preserve_resolver: true`)

This gives you clear organization by schema while maintaining Go's package requirements.
