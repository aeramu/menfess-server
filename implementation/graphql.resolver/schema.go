package resolver

//Schema grahql
var Schema = `
  	schema{
		query: Query
		mutation: Mutation
  	}
  	type Query{
		menfessPost(id: ID!): MenfessPost!
		menfessPostList(first: Int, after: ID, sort: Int): MenfessPostConnection!
		menfessPostRooms(ids: [ID!]!, first: Int, after: ID): MenfessPostConnection!
		menfessRoomList: MenfessRoomConnection!
		menfessAvatarList: [String!]!
	}
	type Mutation{
		postMenfessPost(name: String!, avatar: String!, body: String!, parentID: ID, repostID: ID, roomID: ID): MenfessPost!
		upvoteMenfessPost(postID: ID!): MenfessPost!
		downvoteMenfessPost(postID: ID!): MenfessPost!
	}
	type MenfessPost{
		id: ID!
		timestamp: Int!
		name: String!
		avatar: String!
		body: String!
		replyCount: Int!
		upvoteCount: Int!
		downvoteCount: Int!
		upvoted: Boolean!
		downvoted: Boolean!
		parent: MenfessPost
		repost: MenfessPost
		child(first: Int, after: ID, before: ID, sort: Int): MenfessPostConnection!
		room: String!
	}
	type MenfessRoom{
		id: ID!
		name: String!
		avatar: String!
	}
	type MenfessPostConnection{
		edges: [MenfessPost]!
		pageInfo: PageInfo!
	}
	type MenfessRoomConnection{
		edges: [MenfessRoom]!
		pageInfo: PageInfo!
	}
	type PageInfo{
		startCursor: ID
		endCursor: ID
	}
`
