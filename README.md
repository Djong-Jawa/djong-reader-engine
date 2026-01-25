# Setting Environment
1. If you are using vfox, you can follow this command : 
'''
vfox use golang
'''
2. Then you input the golang version

# Install all packages
go mod download

# Running the apps 
go run server.go

# Generate the gqlgen
gqlgen generate

# Update graphql after added new table
1. make a graphqls in directory graph/schema
2. make a file with name [modulename].graphqls, for example mstsalespipeline.graphqls
3. Then use this commend to generate resolver
'''
gqlgen generate
'''

4. follow the instruction from that gqlgen generate command response

# More detail
If you would like to become as contributor, then visit our website https://djongjawa.com





