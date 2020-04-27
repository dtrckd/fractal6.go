// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type AddCommentInput struct {
	CreatedAt string   `json:"createdAt,omitempty"`
	CreatedBy *UserRef `json:"createdBy,omitempty"`
	Message   *string  `json:"message,omitempty"`
	Void      *string  `json:"_VOID,omitempty"`
}

type AddCommentPayload struct {
	Comment []*Comment `json:"comment,omitempty"`
	NumUids *int       `json:"numUids,omitempty"`
}

type AddLabelInput struct {
	Name  string  `json:"name,omitempty"`
	Color *string `json:"color,omitempty"`
}

type AddLabelPayload struct {
	Label   []*Label `json:"label,omitempty"`
	NumUids *int     `json:"numUids,omitempty"`
}

type AddMandateInput struct {
	CreatedAt        string   `json:"createdAt,omitempty"`
	CreatedBy        *UserRef `json:"createdBy,omitempty"`
	Message          *string  `json:"message,omitempty"`
	Purpose          string   `json:"purpose,omitempty"`
	Responsabilities *string  `json:"responsabilities,omitempty"`
	Domains          []string `json:"domains,omitempty"`
}

type AddMandatePayload struct {
	Mandate []*Mandate `json:"mandate,omitempty"`
	NumUids *int       `json:"numUids,omitempty"`
}

type AddNodeInput struct {
	CreatedAt    string        `json:"createdAt,omitempty"`
	CreatedBy    *UserRef      `json:"createdBy,omitempty"`
	Parent       *NodeRef      `json:"parent,omitempty"`
	Children     []*NodeRef    `json:"children,omitempty"`
	Type         NodeType      `json:"type_,omitempty"`
	Name         string        `json:"name,omitempty"`
	Nameid       string        `json:"nameid,omitempty"`
	Rootnameid   string        `json:"rootnameid,omitempty"`
	Mandate      *MandateRef   `json:"mandate,omitempty"`
	TensionsOut  []*TensionRef `json:"tensions_out,omitempty"`
	TensionsIn   []*TensionRef `json:"tensions_in,omitempty"`
	NTensionsOut *int          `json:"n_tensions_out,omitempty"`
	NTensionsIn  *int          `json:"n_tensions_in,omitempty"`
	NChildren    *int          `json:"n_children,omitempty"`
	IsRoot       bool          `json:"isRoot"`
	FirstLink    *UserRef      `json:"first_link,omitempty"`
	SecondLink   *UserRef      `json:"second_link,omitempty"`
	Skills       []string      `json:"skills,omitempty"`
}

type AddNodePayload struct {
	Node    []*Node `json:"node,omitempty"`
	NumUids *int    `json:"numUids,omitempty"`
}

type AddTensionInput struct {
	CreatedAt string        `json:"createdAt,omitempty"`
	CreatedBy *UserRef      `json:"createdBy,omitempty"`
	Message   *string       `json:"message,omitempty"`
	Title     string        `json:"title,omitempty"`
	Type      TensionType   `json:"type_,omitempty"`
	Emitter   *NodeRef      `json:"emitter,omitempty"`
	Receiver  *NodeRef      `json:"receiver,omitempty"`
	Comments  []*CommentRef `json:"comments,omitempty"`
	Labels    []*LabelRef   `json:"labels,omitempty"`
	NComments *int          `json:"n_comments,omitempty"`
}

type AddTensionPayload struct {
	Tension []*Tension `json:"tension,omitempty"`
	NumUids *int       `json:"numUids,omitempty"`
}

type AddUserInput struct {
	CreatedAt   string     `json:"createdAt,omitempty"`
	Username    string     `json:"username,omitempty"`
	Fullname    *string    `json:"fullname,omitempty"`
	Password    string     `json:"password,omitempty"`
	Roles       []*NodeRef `json:"roles,omitempty"`
	BackedRoles []*NodeRef `json:"backed_roles,omitempty"`
	Bio         *string    `json:"bio,omitempty"`
}

type AddUserPayload struct {
	User    []*User `json:"user,omitempty"`
	NumUids *int    `json:"numUids,omitempty"`
}

type Comment struct {
	Message   string `json:"message,omitempty"`
	ID        string `json:"id,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	CreatedBy *User  `json:"createdBy,omitempty"`
}

type CommentFilter struct {
	ID        []string              `json:"id,omitempty"`
	CreatedAt *DateTimeFilter       `json:"createdAt,omitempty"`
	Message   *StringFullTextFilter `json:"message,omitempty"`
	And       *CommentFilter        `json:"and,omitempty"`
	Or        *CommentFilter        `json:"or,omitempty"`
	Not       *CommentFilter        `json:"not,omitempty"`
}

type CommentOrder struct {
	Asc  *CommentOrderable `json:"asc,omitempty"`
	Desc *CommentOrderable `json:"desc,omitempty"`
	Then *CommentOrder     `json:"then,omitempty"`
}

type CommentPatch struct {
	CreatedAt *string  `json:"createdAt,omitempty"`
	CreatedBy *UserRef `json:"createdBy,omitempty"`
	Message   *string  `json:"message,omitempty"`
	Void      *string  `json:"_VOID,omitempty"`
}

type CommentRef struct {
	ID        *string  `json:"id,omitempty"`
	CreatedAt *string  `json:"createdAt,omitempty"`
	CreatedBy *UserRef `json:"createdBy,omitempty"`
	Message   *string  `json:"message,omitempty"`
	Void      *string  `json:"_VOID,omitempty"`
}

type DateTimeFilter struct {
	Eq *string `json:"eq,omitempty"`
	Le *string `json:"le,omitempty"`
	Lt *string `json:"lt,omitempty"`
	Ge *string `json:"ge,omitempty"`
	Gt *string `json:"gt,omitempty"`
}

type DeleteCommentPayload struct {
	Msg     *string `json:"msg,omitempty"`
	NumUids *int    `json:"numUids,omitempty"`
}

type DeleteLabelPayload struct {
	Msg     *string `json:"msg,omitempty"`
	NumUids *int    `json:"numUids,omitempty"`
}

type DeleteMandatePayload struct {
	Msg     *string `json:"msg,omitempty"`
	NumUids *int    `json:"numUids,omitempty"`
}

type DeleteNodePayload struct {
	Msg     *string `json:"msg,omitempty"`
	NumUids *int    `json:"numUids,omitempty"`
}

type DeletePostPayload struct {
	Msg     *string `json:"msg,omitempty"`
	NumUids *int    `json:"numUids,omitempty"`
}

type DeleteTensionPayload struct {
	Msg     *string `json:"msg,omitempty"`
	NumUids *int    `json:"numUids,omitempty"`
}

type DeleteUserPayload struct {
	Msg     *string `json:"msg,omitempty"`
	NumUids *int    `json:"numUids,omitempty"`
}

type FloatFilter struct {
	Eq *float64 `json:"eq,omitempty"`
	Le *float64 `json:"le,omitempty"`
	Lt *float64 `json:"lt,omitempty"`
	Ge *float64 `json:"ge,omitempty"`
	Gt *float64 `json:"gt,omitempty"`
}

type IntFilter struct {
	Eq *int `json:"eq,omitempty"`
	Le *int `json:"le,omitempty"`
	Lt *int `json:"lt,omitempty"`
	Ge *int `json:"ge,omitempty"`
	Gt *int `json:"gt,omitempty"`
}

type Label struct {
	ID    string  `json:"id,omitempty"`
	Name  string  `json:"name,omitempty"`
	Color *string `json:"color,omitempty"`
}

type LabelFilter struct {
	ID   []string          `json:"id,omitempty"`
	Name *StringHashFilter `json:"name,omitempty"`
	And  *LabelFilter      `json:"and,omitempty"`
	Or   *LabelFilter      `json:"or,omitempty"`
	Not  *LabelFilter      `json:"not,omitempty"`
}

type LabelOrder struct {
	Asc  *LabelOrderable `json:"asc,omitempty"`
	Desc *LabelOrderable `json:"desc,omitempty"`
	Then *LabelOrder     `json:"then,omitempty"`
}

type LabelPatch struct {
	Color *string `json:"color,omitempty"`
}

type LabelRef struct {
	ID    *string `json:"id,omitempty"`
	Name  *string `json:"name,omitempty"`
	Color *string `json:"color,omitempty"`
}

type Mandate struct {
	Purpose          string   `json:"purpose,omitempty"`
	Responsabilities *string  `json:"responsabilities,omitempty"`
	Domains          []string `json:"domains,omitempty"`
	ID               string   `json:"id,omitempty"`
	CreatedAt        string   `json:"createdAt,omitempty"`
	CreatedBy        *User    `json:"createdBy,omitempty"`
	Message          *string  `json:"message,omitempty"`
}

type MandateFilter struct {
	ID        []string              `json:"id,omitempty"`
	CreatedAt *DateTimeFilter       `json:"createdAt,omitempty"`
	Message   *StringFullTextFilter `json:"message,omitempty"`
	Purpose   *StringFullTextFilter `json:"purpose,omitempty"`
	And       *MandateFilter        `json:"and,omitempty"`
	Or        *MandateFilter        `json:"or,omitempty"`
	Not       *MandateFilter        `json:"not,omitempty"`
}

type MandateOrder struct {
	Asc  *MandateOrderable `json:"asc,omitempty"`
	Desc *MandateOrderable `json:"desc,omitempty"`
	Then *MandateOrder     `json:"then,omitempty"`
}

type MandatePatch struct {
	CreatedAt        *string  `json:"createdAt,omitempty"`
	CreatedBy        *UserRef `json:"createdBy,omitempty"`
	Message          *string  `json:"message,omitempty"`
	Purpose          *string  `json:"purpose,omitempty"`
	Responsabilities *string  `json:"responsabilities,omitempty"`
	Domains          []string `json:"domains,omitempty"`
}

type MandateRef struct {
	ID               *string  `json:"id,omitempty"`
	CreatedAt        *string  `json:"createdAt,omitempty"`
	CreatedBy        *UserRef `json:"createdBy,omitempty"`
	Message          *string  `json:"message,omitempty"`
	Purpose          *string  `json:"purpose,omitempty"`
	Responsabilities *string  `json:"responsabilities,omitempty"`
	Domains          []string `json:"domains,omitempty"`
}

type Node struct {
	ID           string     `json:"id,omitempty"`
	CreatedAt    string     `json:"createdAt,omitempty"`
	CreatedBy    *User      `json:"createdBy,omitempty"`
	Parent       *Node      `json:"parent,omitempty"`
	Children     []*Node    `json:"children,omitempty"`
	Type         NodeType   `json:"type_,omitempty"`
	Name         string     `json:"name,omitempty"`
	Nameid       string     `json:"nameid,omitempty"`
	Rootnameid   string     `json:"rootnameid,omitempty"`
	Mandate      *Mandate   `json:"mandate,omitempty"`
	TensionsOut  []*Tension `json:"tensions_out,omitempty"`
	TensionsIn   []*Tension `json:"tensions_in,omitempty"`
	NTensionsOut *int       `json:"n_tensions_out,omitempty"`
	NTensionsIn  *int       `json:"n_tensions_in,omitempty"`
	NChildren    *int       `json:"n_children,omitempty"`
	IsRoot       bool       `json:"isRoot"`
	FirstLink    *User      `json:"first_link,omitempty"`
	SecondLink   *User      `json:"second_link,omitempty"`
	Skills       []string   `json:"skills,omitempty"`
}

type NodeFilter struct {
	ID         []string          `json:"id,omitempty"`
	CreatedAt  *DateTimeFilter   `json:"createdAt,omitempty"`
	Type       *NodeTypeHash     `json:"type_,omitempty"`
	Name       *StringTermFilter `json:"name,omitempty"`
	Nameid     *StringHashFilter `json:"nameid,omitempty"`
	Rootnameid *StringHashFilter `json:"rootnameid,omitempty"`
	IsRoot     *bool             `json:"isRoot"`
	Skills     *StringTermFilter `json:"skills,omitempty"`
	And        *NodeFilter       `json:"and,omitempty"`
	Or         *NodeFilter       `json:"or,omitempty"`
	Not        *NodeFilter       `json:"not,omitempty"`
}

type NodeOrder struct {
	Asc  *NodeOrderable `json:"asc,omitempty"`
	Desc *NodeOrderable `json:"desc,omitempty"`
	Then *NodeOrder     `json:"then,omitempty"`
}

type NodePatch struct {
	CreatedAt    *string       `json:"createdAt,omitempty"`
	CreatedBy    *UserRef      `json:"createdBy,omitempty"`
	Parent       *NodeRef      `json:"parent,omitempty"`
	Children     []*NodeRef    `json:"children,omitempty"`
	Type         *NodeType     `json:"type_,omitempty"`
	Name         *string       `json:"name,omitempty"`
	Rootnameid   *string       `json:"rootnameid,omitempty"`
	Mandate      *MandateRef   `json:"mandate,omitempty"`
	TensionsOut  []*TensionRef `json:"tensions_out,omitempty"`
	TensionsIn   []*TensionRef `json:"tensions_in,omitempty"`
	NTensionsOut *int          `json:"n_tensions_out,omitempty"`
	NTensionsIn  *int          `json:"n_tensions_in,omitempty"`
	NChildren    *int          `json:"n_children,omitempty"`
	IsRoot       *bool         `json:"isRoot"`
	FirstLink    *UserRef      `json:"first_link,omitempty"`
	SecondLink   *UserRef      `json:"second_link,omitempty"`
	Skills       []string      `json:"skills,omitempty"`
}

type NodeRef struct {
	ID           *string       `json:"id,omitempty"`
	CreatedAt    *string       `json:"createdAt,omitempty"`
	CreatedBy    *UserRef      `json:"createdBy,omitempty"`
	Parent       *NodeRef      `json:"parent,omitempty"`
	Children     []*NodeRef    `json:"children,omitempty"`
	Type         *NodeType     `json:"type_,omitempty"`
	Name         *string       `json:"name,omitempty"`
	Nameid       *string       `json:"nameid,omitempty"`
	Rootnameid   *string       `json:"rootnameid,omitempty"`
	Mandate      *MandateRef   `json:"mandate,omitempty"`
	TensionsOut  []*TensionRef `json:"tensions_out,omitempty"`
	TensionsIn   []*TensionRef `json:"tensions_in,omitempty"`
	NTensionsOut *int          `json:"n_tensions_out,omitempty"`
	NTensionsIn  *int          `json:"n_tensions_in,omitempty"`
	NChildren    *int          `json:"n_children,omitempty"`
	IsRoot       *bool         `json:"isRoot"`
	FirstLink    *UserRef      `json:"first_link,omitempty"`
	SecondLink   *UserRef      `json:"second_link,omitempty"`
	Skills       []string      `json:"skills,omitempty"`
}

type NodeTypeHash struct {
	Eq NodeType `json:"eq,omitempty"`
}

type Post struct {
	ID        string  `json:"id,omitempty"`
	CreatedAt string  `json:"createdAt,omitempty"`
	CreatedBy *User   `json:"createdBy,omitempty"`
	Message   *string `json:"message,omitempty"`
}

type PostFilter struct {
	ID        []string              `json:"id,omitempty"`
	CreatedAt *DateTimeFilter       `json:"createdAt,omitempty"`
	Message   *StringFullTextFilter `json:"message,omitempty"`
	And       *PostFilter           `json:"and,omitempty"`
	Or        *PostFilter           `json:"or,omitempty"`
	Not       *PostFilter           `json:"not,omitempty"`
}

type PostOrder struct {
	Asc  *PostOrderable `json:"asc,omitempty"`
	Desc *PostOrderable `json:"desc,omitempty"`
	Then *PostOrder     `json:"then,omitempty"`
}

type PostPatch struct {
	CreatedAt *string  `json:"createdAt,omitempty"`
	CreatedBy *UserRef `json:"createdBy,omitempty"`
	Message   *string  `json:"message,omitempty"`
}

type PostRef struct {
	ID string `json:"id,omitempty"`
}

type StringExactFilter struct {
	Eq *string `json:"eq,omitempty"`
	Le *string `json:"le,omitempty"`
	Lt *string `json:"lt,omitempty"`
	Ge *string `json:"ge,omitempty"`
	Gt *string `json:"gt,omitempty"`
}

type StringFullTextFilter struct {
	Alloftext *string `json:"alloftext,omitempty"`
	Anyoftext *string `json:"anyoftext,omitempty"`
}

type StringHashFilter struct {
	Eq *string `json:"eq,omitempty"`
}

type StringRegExpFilter struct {
	Regexp *string `json:"regexp,omitempty"`
}

type StringTermFilter struct {
	Allofterms *string `json:"allofterms,omitempty"`
	Anyofterms *string `json:"anyofterms,omitempty"`
}

type Tension struct {
	Title     string      `json:"title,omitempty"`
	Type      TensionType `json:"type_,omitempty"`
	Emitter   *Node       `json:"emitter,omitempty"`
	Receiver  *Node       `json:"receiver,omitempty"`
	Comments  []*Comment  `json:"comments,omitempty"`
	Labels    []*Label    `json:"labels,omitempty"`
	NComments *int        `json:"n_comments,omitempty"`
	ID        string      `json:"id,omitempty"`
	CreatedAt string      `json:"createdAt,omitempty"`
	CreatedBy *User       `json:"createdBy,omitempty"`
	Message   *string     `json:"message,omitempty"`
}

type TensionFilter struct {
	ID        []string              `json:"id,omitempty"`
	CreatedAt *DateTimeFilter       `json:"createdAt,omitempty"`
	Message   *StringFullTextFilter `json:"message,omitempty"`
	Title     *StringTermFilter     `json:"title,omitempty"`
	Type      *TensionTypeHash      `json:"type_,omitempty"`
	And       *TensionFilter        `json:"and,omitempty"`
	Or        *TensionFilter        `json:"or,omitempty"`
	Not       *TensionFilter        `json:"not,omitempty"`
}

type TensionOrder struct {
	Asc  *TensionOrderable `json:"asc,omitempty"`
	Desc *TensionOrderable `json:"desc,omitempty"`
	Then *TensionOrder     `json:"then,omitempty"`
}

type TensionPatch struct {
	CreatedAt *string       `json:"createdAt,omitempty"`
	CreatedBy *UserRef      `json:"createdBy,omitempty"`
	Message   *string       `json:"message,omitempty"`
	Title     *string       `json:"title,omitempty"`
	Type      *TensionType  `json:"type_,omitempty"`
	Emitter   *NodeRef      `json:"emitter,omitempty"`
	Receiver  *NodeRef      `json:"receiver,omitempty"`
	Comments  []*CommentRef `json:"comments,omitempty"`
	Labels    []*LabelRef   `json:"labels,omitempty"`
	NComments *int          `json:"n_comments,omitempty"`
}

type TensionRef struct {
	ID        *string       `json:"id,omitempty"`
	CreatedAt *string       `json:"createdAt,omitempty"`
	CreatedBy *UserRef      `json:"createdBy,omitempty"`
	Message   *string       `json:"message,omitempty"`
	Title     *string       `json:"title,omitempty"`
	Type      *TensionType  `json:"type_,omitempty"`
	Emitter   *NodeRef      `json:"emitter,omitempty"`
	Receiver  *NodeRef      `json:"receiver,omitempty"`
	Comments  []*CommentRef `json:"comments,omitempty"`
	Labels    []*LabelRef   `json:"labels,omitempty"`
	NComments *int          `json:"n_comments,omitempty"`
}

type TensionTypeHash struct {
	Eq TensionType `json:"eq,omitempty"`
}

type UpdateCommentInput struct {
	Filter *CommentFilter `json:"filter,omitempty"`
	Set    *CommentPatch  `json:"set,omitempty"`
	Remove *CommentPatch  `json:"remove,omitempty"`
}

type UpdateCommentPayload struct {
	Comment []*Comment `json:"comment,omitempty"`
	NumUids *int       `json:"numUids,omitempty"`
}

type UpdateLabelInput struct {
	Filter *LabelFilter `json:"filter,omitempty"`
	Set    *LabelPatch  `json:"set,omitempty"`
	Remove *LabelPatch  `json:"remove,omitempty"`
}

type UpdateLabelPayload struct {
	Label   []*Label `json:"label,omitempty"`
	NumUids *int     `json:"numUids,omitempty"`
}

type UpdateMandateInput struct {
	Filter *MandateFilter `json:"filter,omitempty"`
	Set    *MandatePatch  `json:"set,omitempty"`
	Remove *MandatePatch  `json:"remove,omitempty"`
}

type UpdateMandatePayload struct {
	Mandate []*Mandate `json:"mandate,omitempty"`
	NumUids *int       `json:"numUids,omitempty"`
}

type UpdateNodeInput struct {
	Filter *NodeFilter `json:"filter,omitempty"`
	Set    *NodePatch  `json:"set,omitempty"`
	Remove *NodePatch  `json:"remove,omitempty"`
}

type UpdateNodePayload struct {
	Node    []*Node `json:"node,omitempty"`
	NumUids *int    `json:"numUids,omitempty"`
}

type UpdatePostInput struct {
	Filter *PostFilter `json:"filter,omitempty"`
	Set    *PostPatch  `json:"set,omitempty"`
	Remove *PostPatch  `json:"remove,omitempty"`
}

type UpdatePostPayload struct {
	Post    []*Post `json:"post,omitempty"`
	NumUids *int    `json:"numUids,omitempty"`
}

type UpdateTensionInput struct {
	Filter *TensionFilter `json:"filter,omitempty"`
	Set    *TensionPatch  `json:"set,omitempty"`
	Remove *TensionPatch  `json:"remove,omitempty"`
}

type UpdateTensionPayload struct {
	Tension []*Tension `json:"tension,omitempty"`
	NumUids *int       `json:"numUids,omitempty"`
}

type UpdateUserInput struct {
	Filter *UserFilter `json:"filter,omitempty"`
	Set    *UserPatch  `json:"set,omitempty"`
	Remove *UserPatch  `json:"remove,omitempty"`
}

type UpdateUserPayload struct {
	User    []*User `json:"user,omitempty"`
	NumUids *int    `json:"numUids,omitempty"`
}

type User struct {
	ID          string  `json:"id,omitempty"`
	CreatedAt   string  `json:"createdAt,omitempty"`
	Username    string  `json:"username,omitempty"`
	Fullname    *string `json:"fullname,omitempty"`
	Password    string  `json:"password,omitempty"`
	Roles       []*Node `json:"roles,omitempty"`
	BackedRoles []*Node `json:"backed_roles,omitempty"`
	Bio         *string `json:"bio,omitempty"`
}

type UserFilter struct {
	ID        []string          `json:"id,omitempty"`
	CreatedAt *DateTimeFilter   `json:"createdAt,omitempty"`
	Username  *StringHashFilter `json:"username,omitempty"`
	And       *UserFilter       `json:"and,omitempty"`
	Or        *UserFilter       `json:"or,omitempty"`
	Not       *UserFilter       `json:"not,omitempty"`
}

type UserOrder struct {
	Asc  *UserOrderable `json:"asc,omitempty"`
	Desc *UserOrderable `json:"desc,omitempty"`
	Then *UserOrder     `json:"then,omitempty"`
}

type UserPatch struct {
	CreatedAt   *string    `json:"createdAt,omitempty"`
	Fullname    *string    `json:"fullname,omitempty"`
	Password    *string    `json:"password,omitempty"`
	Roles       []*NodeRef `json:"roles,omitempty"`
	BackedRoles []*NodeRef `json:"backed_roles,omitempty"`
	Bio         *string    `json:"bio,omitempty"`
}

type UserRef struct {
	ID          *string    `json:"id,omitempty"`
	CreatedAt   *string    `json:"createdAt,omitempty"`
	Username    *string    `json:"username,omitempty"`
	Fullname    *string    `json:"fullname,omitempty"`
	Password    *string    `json:"password,omitempty"`
	Roles       []*NodeRef `json:"roles,omitempty"`
	BackedRoles []*NodeRef `json:"backed_roles,omitempty"`
	Bio         *string    `json:"bio,omitempty"`
}

type CommentOrderable string

const (
	CommentOrderableCreatedAt CommentOrderable = "createdAt"
	CommentOrderableMessage   CommentOrderable = "message"
	CommentOrderableVoid      CommentOrderable = "_VOID"
)

var AllCommentOrderable = []CommentOrderable{
	CommentOrderableCreatedAt,
	CommentOrderableMessage,
	CommentOrderableVoid,
}

func (e CommentOrderable) IsValid() bool {
	switch e {
	case CommentOrderableCreatedAt, CommentOrderableMessage, CommentOrderableVoid:
		return true
	}
	return false
}

func (e CommentOrderable) String() string {
	return string(e)
}

func (e *CommentOrderable) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = CommentOrderable(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid CommentOrderable", str)
	}
	return nil
}

func (e CommentOrderable) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type DgraphIndex string

const (
	DgraphIndexInt      DgraphIndex = "int"
	DgraphIndexFloat    DgraphIndex = "float"
	DgraphIndexBool     DgraphIndex = "bool"
	DgraphIndexHash     DgraphIndex = "hash"
	DgraphIndexExact    DgraphIndex = "exact"
	DgraphIndexTerm     DgraphIndex = "term"
	DgraphIndexFulltext DgraphIndex = "fulltext"
	DgraphIndexTrigram  DgraphIndex = "trigram"
	DgraphIndexRegexp   DgraphIndex = "regexp"
	DgraphIndexYear     DgraphIndex = "year"
	DgraphIndexMonth    DgraphIndex = "month"
	DgraphIndexDay      DgraphIndex = "day"
	DgraphIndexHour     DgraphIndex = "hour"
)

var AllDgraphIndex = []DgraphIndex{
	DgraphIndexInt,
	DgraphIndexFloat,
	DgraphIndexBool,
	DgraphIndexHash,
	DgraphIndexExact,
	DgraphIndexTerm,
	DgraphIndexFulltext,
	DgraphIndexTrigram,
	DgraphIndexRegexp,
	DgraphIndexYear,
	DgraphIndexMonth,
	DgraphIndexDay,
	DgraphIndexHour,
}

func (e DgraphIndex) IsValid() bool {
	switch e {
	case DgraphIndexInt, DgraphIndexFloat, DgraphIndexBool, DgraphIndexHash, DgraphIndexExact, DgraphIndexTerm, DgraphIndexFulltext, DgraphIndexTrigram, DgraphIndexRegexp, DgraphIndexYear, DgraphIndexMonth, DgraphIndexDay, DgraphIndexHour:
		return true
	}
	return false
}

func (e DgraphIndex) String() string {
	return string(e)
}

func (e *DgraphIndex) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DgraphIndex(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DgraphIndex", str)
	}
	return nil
}

func (e DgraphIndex) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type LabelOrderable string

const (
	LabelOrderableName  LabelOrderable = "name"
	LabelOrderableColor LabelOrderable = "color"
)

var AllLabelOrderable = []LabelOrderable{
	LabelOrderableName,
	LabelOrderableColor,
}

func (e LabelOrderable) IsValid() bool {
	switch e {
	case LabelOrderableName, LabelOrderableColor:
		return true
	}
	return false
}

func (e LabelOrderable) String() string {
	return string(e)
}

func (e *LabelOrderable) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = LabelOrderable(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid LabelOrderable", str)
	}
	return nil
}

func (e LabelOrderable) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type MandateOrderable string

const (
	MandateOrderableCreatedAt        MandateOrderable = "createdAt"
	MandateOrderableMessage          MandateOrderable = "message"
	MandateOrderablePurpose          MandateOrderable = "purpose"
	MandateOrderableResponsabilities MandateOrderable = "responsabilities"
	MandateOrderableDomains          MandateOrderable = "domains"
)

var AllMandateOrderable = []MandateOrderable{
	MandateOrderableCreatedAt,
	MandateOrderableMessage,
	MandateOrderablePurpose,
	MandateOrderableResponsabilities,
	MandateOrderableDomains,
}

func (e MandateOrderable) IsValid() bool {
	switch e {
	case MandateOrderableCreatedAt, MandateOrderableMessage, MandateOrderablePurpose, MandateOrderableResponsabilities, MandateOrderableDomains:
		return true
	}
	return false
}

func (e MandateOrderable) String() string {
	return string(e)
}

func (e *MandateOrderable) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MandateOrderable(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MandateOrderable", str)
	}
	return nil
}

func (e MandateOrderable) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type NodeOrderable string

const (
	NodeOrderableCreatedAt    NodeOrderable = "createdAt"
	NodeOrderableName         NodeOrderable = "name"
	NodeOrderableNameid       NodeOrderable = "nameid"
	NodeOrderableRootnameid   NodeOrderable = "rootnameid"
	NodeOrderableNTensionsOut NodeOrderable = "n_tensions_out"
	NodeOrderableNTensionsIn  NodeOrderable = "n_tensions_in"
	NodeOrderableNChildren    NodeOrderable = "n_children"
	NodeOrderableSkills       NodeOrderable = "skills"
)

var AllNodeOrderable = []NodeOrderable{
	NodeOrderableCreatedAt,
	NodeOrderableName,
	NodeOrderableNameid,
	NodeOrderableRootnameid,
	NodeOrderableNTensionsOut,
	NodeOrderableNTensionsIn,
	NodeOrderableNChildren,
	NodeOrderableSkills,
}

func (e NodeOrderable) IsValid() bool {
	switch e {
	case NodeOrderableCreatedAt, NodeOrderableName, NodeOrderableNameid, NodeOrderableRootnameid, NodeOrderableNTensionsOut, NodeOrderableNTensionsIn, NodeOrderableNChildren, NodeOrderableSkills:
		return true
	}
	return false
}

func (e NodeOrderable) String() string {
	return string(e)
}

func (e *NodeOrderable) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = NodeOrderable(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid NodeOrderable", str)
	}
	return nil
}

func (e NodeOrderable) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type NodeType string

const (
	NodeTypeCircle NodeType = "Circle"
	NodeTypeRole   NodeType = "Role"
)

var AllNodeType = []NodeType{
	NodeTypeCircle,
	NodeTypeRole,
}

func (e NodeType) IsValid() bool {
	switch e {
	case NodeTypeCircle, NodeTypeRole:
		return true
	}
	return false
}

func (e NodeType) String() string {
	return string(e)
}

func (e *NodeType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = NodeType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid NodeType", str)
	}
	return nil
}

func (e NodeType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type PostOrderable string

const (
	PostOrderableCreatedAt PostOrderable = "createdAt"
	PostOrderableMessage   PostOrderable = "message"
)

var AllPostOrderable = []PostOrderable{
	PostOrderableCreatedAt,
	PostOrderableMessage,
}

func (e PostOrderable) IsValid() bool {
	switch e {
	case PostOrderableCreatedAt, PostOrderableMessage:
		return true
	}
	return false
}

func (e PostOrderable) String() string {
	return string(e)
}

func (e *PostOrderable) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PostOrderable(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PostOrderable", str)
	}
	return nil
}

func (e PostOrderable) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type TensionOrderable string

const (
	TensionOrderableCreatedAt TensionOrderable = "createdAt"
	TensionOrderableMessage   TensionOrderable = "message"
	TensionOrderableTitle     TensionOrderable = "title"
	TensionOrderableNComments TensionOrderable = "n_comments"
)

var AllTensionOrderable = []TensionOrderable{
	TensionOrderableCreatedAt,
	TensionOrderableMessage,
	TensionOrderableTitle,
	TensionOrderableNComments,
}

func (e TensionOrderable) IsValid() bool {
	switch e {
	case TensionOrderableCreatedAt, TensionOrderableMessage, TensionOrderableTitle, TensionOrderableNComments:
		return true
	}
	return false
}

func (e TensionOrderable) String() string {
	return string(e)
}

func (e *TensionOrderable) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TensionOrderable(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TensionOrderable", str)
	}
	return nil
}

func (e TensionOrderable) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type TensionType string

const (
	TensionTypeGovernance  TensionType = "Governance"
	TensionTypeOperational TensionType = "Operational"
	TensionTypePersonal    TensionType = "Personal"
	TensionTypeHelp        TensionType = "Help"
	TensionTypeAlert       TensionType = "Alert"
)

var AllTensionType = []TensionType{
	TensionTypeGovernance,
	TensionTypeOperational,
	TensionTypePersonal,
	TensionTypeHelp,
	TensionTypeAlert,
}

func (e TensionType) IsValid() bool {
	switch e {
	case TensionTypeGovernance, TensionTypeOperational, TensionTypePersonal, TensionTypeHelp, TensionTypeAlert:
		return true
	}
	return false
}

func (e TensionType) String() string {
	return string(e)
}

func (e *TensionType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TensionType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TensionType", str)
	}
	return nil
}

func (e TensionType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type UserOrderable string

const (
	UserOrderableCreatedAt UserOrderable = "createdAt"
	UserOrderableUsername  UserOrderable = "username"
	UserOrderableFullname  UserOrderable = "fullname"
	UserOrderablePassword  UserOrderable = "password"
	UserOrderableBio       UserOrderable = "bio"
)

var AllUserOrderable = []UserOrderable{
	UserOrderableCreatedAt,
	UserOrderableUsername,
	UserOrderableFullname,
	UserOrderablePassword,
	UserOrderableBio,
}

func (e UserOrderable) IsValid() bool {
	switch e {
	case UserOrderableCreatedAt, UserOrderableUsername, UserOrderableFullname, UserOrderablePassword, UserOrderableBio:
		return true
	}
	return false
}

func (e UserOrderable) String() string {
	return string(e)
}

func (e *UserOrderable) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = UserOrderable(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UserOrderable", str)
	}
	return nil
}

func (e UserOrderable) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
