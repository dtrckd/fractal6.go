package model

// JsonAtom is a general interface
// for decoding unknonw structure
type JsonAtom = map[string]interface{}

type NodeId struct {
    Nameid string `json:"nameid"`
}

type MemberNode struct {
    CreatedAt string    `json:"createdAt"`
    Name string         `json:"name"`
    Nameid string       `json:"nameid"`
    Rootnameid string   `json:"rootnameid"`
	RoleType *RoleType  `json:"role_type,omitempty"`
	FirstLink *User_    `json:"first_link,omitempty"`
	Parent *NodeId      `json:"parent,omitempty"`

}

type User_ struct {
    Username string `json:"username"`
    Name *string    `json:"name,omitempty"`
}

var TensionHookPayload string = `{
  uid
  Tension.action
  Tension.emitter {
    Node.nameid
  }
  Tension.receiver {
    Node.nameid
    Node.isPrivate
    Node.charac {
        NodeCharac.userCanJoin
        NodeCharac.mode
    }
  }
  Tension.blobs %s {
    uid
    Blob.blob_type
    Blob.md
    Blob.node {
      NodeFragment.name
      NodeFragment.nameid
      NodeFragment.type_
      NodeFragment.about
      NodeFragment.mandate {
        expand(_all_)
      }

      NodeFragment.first_link
      NodeFragment.second_link
      NodeFragment.skills
      NodeFragment.role_type

      NodeFragment.children {
        NodeFragment.first_link
        NodeFragment.role_type
      }
    }
  }
}`
