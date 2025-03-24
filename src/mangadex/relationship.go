package mangadex

type Relationship struct {
	ID   string
	Type string
}

var (
	RelationshipTypeAuthor = "author"
	RelationshipTypeArtist = "artist"
	RelationshipTypeCoverArt = "cover_art"
	RelationshipTypeCreator = "creator"
	RelationshipTypeScanlationGroup = "scanlation_group"
)
